package connection

import (
	"log"
	"net"
	"time"
)

const (
	PEER_IP   = "192.168.100.113"
	PEER_PORT = 51820
)

var conn *net.UDPConn

func InitClient() {
	var err error
	addr := &net.UDPAddr{IP: net.ParseIP(PEER_IP), Port: PEER_PORT}
	conn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("InitClient dial failed: %v", err)
	}
	log.Printf("InitClient: connected to %s (remote=%v)", addr.String(), conn.RemoteAddr())
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
	addr := &net.UDPAddr{
		IP:   net.ParseIP(PEER_IP),
		Port: PEER_PORT,
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}
	defer conn.Close()

	log.Printf("sending dummy packets to %s:%d", PEER_IP, PEER_PORT)

	// dummy packet content
	packet := []byte("hello from client ddd ddd 1221")

	for {
		_, err := conn.Write(packet)
		if err != nil {
			log.Printf("failed to send packet: %v", err)
		} else {
			log.Printf("packet sent")
		}
		time.Sleep(1 * time.Second) // send every second
	}
}
