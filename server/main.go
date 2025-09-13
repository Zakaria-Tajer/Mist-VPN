package main

import (
	"log"
	"net"
	"zakaria/mist-vpn/helpers"
	server "zakaria/mist-vpn/server/reader"

	"github.com/songgao/water"
)

const (
	TUN_IP   = "10.0.0.1/30"  // server TUN IP
	MTU      = "1300"
	BUFFERSIZE = 1600
	UDP_PORT = 51820
)

func main() {
	iface, err := water.New(water.Config{DeviceType: water.TUN})
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Server TUN interface: %s", iface.Name())

    // Bring interface up
    helpers.RunBin("/bin/ip", "link", "set", "dev", iface.Name(), "mtu", "1300")
    helpers.RunBin("/bin/ip", "addr", "add", TUN_IP, "dev", iface.Name())
    helpers.RunBin("/bin/ip", "link", "set", "dev", iface.Name(), "up")

    // Create UDP listener
    addr := net.UDPAddr{IP: net.IPv4zero, Port: UDP_PORT}

    conn, err := net.ListenUDP("udp", &addr)

    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    log.Printf("Listening on UDP %s", conn.LocalAddr())


    
    
    
    // Call your server loop from server_packet_reader.go
    go server.LogTunPackets(iface)
    server.ServerSideReader(conn, iface)
}
