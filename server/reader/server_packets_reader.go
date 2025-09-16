package server

import (
	"log"
	"net"
	"zakaria/mist-vpn/helpers"

	"github.com/songgao/water"
	"golang.org/x/net/ipv4"
)

const (
	SERVER_TUN_IP   = "10.0.0.1"
	CLIENT_TUN_IP   = "10.0.0.2"
	SERVER_UDP_PORT = 41822
	CLIENT_UDP_PORT = 54321
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

		payload := buf[:n]

		udpLen := 8 + len(payload)
		udpHeader := helpers.BuildUDPHeader(CLIENT_UDP_PORT, SERVER_UDP_PORT, udpLen)

		totalLen := 20 + udpLen
		ipHeader := helpers.BuildIPv4Header(
			net.ParseIP(CLIENT_TUN_IP),
			net.ParseIP(SERVER_TUN_IP),
			totalLen,
			0x1234,
			64,
			17,
		)

		packet := append(ipHeader, udpHeader...)
		packet = append(packet, payload...)

		_, err = tunIface.Write(packet)
		if err != nil {
			log.Printf("failed to write packet to TUN: %v", err)
		}

		// log.Printf("Wrapped %d bytes from %s into TUN packet (total %d bytes)", n, addr, len(packet))

		log.Printf("Received %d bytes from %s: %q", n, addr, packet)
	}

}

func ForwardTunToClient(tunIface *water.Interface, clientConn *net.UDPConn, clientAddr *net.UDPAddr) {
	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := tunIface.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		header, _ := ipv4.ParseHeader(buf[:n])
		payload := buf[header.Len:n]

		_, err = clientConn.WriteToUDP(payload, clientAddr)
		if err != nil {
			log.Printf("failed to send to client: %v", err)
		}
	}
}
