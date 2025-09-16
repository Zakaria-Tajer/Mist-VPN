package reader

import (
	"log"
	"net"
	"os"
	"os/exec"
	"zakaria/mist-vpn/client/connection"

	"github.com/songgao/water"
	"golang.org/x/net/ipv4"
)

const (
	// I use TUN interface, so only plain IP packet, no ethernet header + mtu is set to 1300
	BUFFERSIZE = 1600
	MTU        = "1300"
	TUN_IP     = "10.0.0.2"
)

func ReadPacketsFromTun() {

	iface, err := water.New(water.Config{})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("tun interface: %s", iface.Name())

	RunBin("/bin/ip", "link", "set", "dev", iface.Name(), "mtu", MTU)
	RunBin("/bin/ip", "addr", "add", TUN_IP, "dev", iface.Name())
	RunBin("/bin/ip", "link", "set", "dev", iface.Name(), "up")

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

		if header.Dst.Equal(net.ParseIP(TUN_IP)) {

			payload := buf[header.Len:n]

			connection.Client(payload)

		} else {
			continue
		}

		// if header.Protocol != 17 {
		// 	continue
		// }

		// log.Printf("Forwarded UDP packet %d bytes: %s -> %s", n, header.Src, header.Dst)

	}

}

func RunBin(bin string, args ...string) {
	cmd := exec.Command(bin, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}
