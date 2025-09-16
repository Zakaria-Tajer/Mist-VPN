package connection

import (
	"log"
	"net"
	"time"
	"zakaria/mist-vpn/helpers"
)

const (
	CLIENT_TUN_IP   = "10.0.0.2"
	SERVER_TUN_IP   = "10.0.0.1"
	SERVER_UDP_IP   = "192.168.100.113" // server’s LAN/public IP
	SERVER_UDP_PORT = 41822
)

var conn *net.UDPConn

func InitClient() {
	var err error
	addr := &net.UDPAddr{IP: net.ParseIP(SERVER_UDP_IP), Port: SERVER_UDP_PORT}
	conn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("InitClient dial failed: %v", err)
	}
	log.Printf("InitClient: connected to %s (remote=%v)", addr.String(), conn.RemoteAddr())
	// return conn
}

func Client(packet []byte) {
	if conn == nil {
		log.Fatal("Client called before InitClient")
	}
	n, err := conn.Write(packet)
	if err != nil {
		log.Printf("failed to send packet: %v", err)
	} else {
		log.Printf("Sent %d bytes to %v", n, conn.RemoteAddr())
	}
}

func SendDummyContent() {

	payload := []byte("hello from client ddd ddd 1221")
	srcPort := 12345 // any random client port
	dstPort := 8080  // pick something meaningful

	// Build UDP header
	udpLen := 8 + len(payload)
	udpHeader := helpers.BuildUDPHeader(srcPort, dstPort, udpLen)

	// Build IPv4 header
	totalLen := 20 + udpLen
	ipHeader := helpers.BuildIPv4Header(
		net.ParseIP(CLIENT_TUN_IP),
		net.ParseIP(SERVER_TUN_IP),
		totalLen,
		0,
		64,
		17,
	)

	// Final packet = IPv4 header + UDP header + payload
	packet := append(ipHeader, append(udpHeader, payload...)...)

	for {
		_, err := conn.Write(packet)
		if err != nil {
			log.Printf("failed to send packet: %v", err)
		} else {
			log.Printf("packet sent to %s:%d", SERVER_UDP_IP, SERVER_UDP_PORT)
		}
		time.Sleep(1 * time.Second)
	}
}
