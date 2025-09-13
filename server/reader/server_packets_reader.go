package server

import (
	"encoding/binary"
	"log"
	"net"

	"github.com/songgao/water"
	"golang.org/x/net/ipv4"
)

const (
	SERVER_TUN_IP   = "10.0.0.1/30"
	SERVER_UDP_PORT = 51820
	BUFFERSIZE      = 1600
)

func ServerSideReader(conn *net.UDPConn, tunIface *water.Interface) {
	buf := make([]byte, 1600)

	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("failed to receive packet: %v", err)
			continue
		}

		_, err = tunIface.Write(buf[:n])
		if err != nil {
			log.Printf("failed to write packet to TUN: %v", err)
		}

		log.Printf("Received %d bytes from %s", n, addr)
	}

}

func LogTunPackets(iface *water.Interface) {
	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := iface.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		if n < 20 {
			continue
		}

		header, err := ipv4.ParseHeader(buf[:n])
		if err != nil || header == nil {
			log.Printf("failed parse header: %v", err)
			continue
		}

		log.Printf("TUN packet: %s -> %s proto=%d len=%d", header.Src, header.Dst, header.Protocol, n)

		if !header.Dst.Equal(net.IP(SERVER_TUN_IP)) {

			continue
		}
		payload := buf[header.Len:n]

		if header.Protocol == 6 { // TCP
			srcPort := binary.BigEndian.Uint16(payload[0:2])
			dstPort := binary.BigEndian.Uint16(payload[2:4])
			// seqNum := binary.BigEndian.Uint32(payload[4:8])
			// ackNum := binary.BigEndian.Uint32(payload[8:12])
			dataOffset := (payload[12] >> 4) * 4
			tcpData := payload[dataOffset:]
			log.Printf("TCP %d -> %d, payload: %x / %q", srcPort, dstPort, tcpData, tcpData)
		}
	}
}
