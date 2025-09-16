package helpers

import (
	"encoding/binary"
	"net"
)

// buildIPv4Header constructs a minimal IPv4 header (20 bytes) with no options.
func BuildIPv4Header(srcIP, dstIP net.IP, totalLen int, identification uint16, ttl uint8, protocol uint8) []byte {
	// Ensure 4-byte IPv4 addresses
	src := srcIP.To4()
	dst := dstIP.To4()
	if src == nil || dst == nil {
		return nil
	}
	h := make([]byte, 20)
	h[0] = (4 << 4) | 5
	// DSCP/ECN
	h[1] = 0
	binary.BigEndian.PutUint16(h[2:4], uint16(totalLen))
	// Identification
	binary.BigEndian.PutUint16(h[4:6], identification)
	// Flags(0) + Fragment Offset(0)
	binary.BigEndian.PutUint16(h[6:8], 0)
	// TTL
	h[8] = byte(ttl)
	// Protocol (UDP=17)
	h[9] = byte(protocol)
	// Header checksum (set 0 for calculation)
	h[10] = 0
	h[11] = 0
	// Src/Dst
	copy(h[12:16], src)
	copy(h[16:20], dst)
	cs := Ipv4HeaderChecksum(h)
	binary.BigEndian.PutUint16(h[10:12], cs)
	return h
}

// buildUDPHeader constructs a UDP header with checksum set to 0 (optional for IPv4).
func BuildUDPHeader(srcPort, dstPort, length int) []byte {
	h := make([]byte, 8)
	binary.BigEndian.PutUint16(h[0:2], uint16(srcPort))
	binary.BigEndian.PutUint16(h[2:4], uint16(dstPort))
	binary.BigEndian.PutUint16(h[4:6], uint16(length))
	// Checksum 0 indicates no checksum for IPv4 UDP
	h[6] = 0
	h[7] = 0
	return h
}

// ipv4HeaderChecksum computes the Internet checksum for the 20-byte IPv4 header.
func Ipv4HeaderChecksum(header []byte) uint16 {
	var sum uint32
	for i := 0; i < len(header); i += 2 {
		sum += uint32(binary.BigEndian.Uint16(header[i : i+2]))
	}
	for (sum >> 16) != 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	return ^uint16(sum)
}
