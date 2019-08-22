// Package transport provides implementations for simple address-based packet
// transfers.
package transport

import (
	pb "github.com/wollac/autopeering/proto"
)

const (
	// MaxPacketSize specifies the maximum allowed size of packets.
	// Packets larger than this will be cut and thus treated as invalid.
	MaxPacketSize = 1280
)

// Transport is generic network connection to transfer protobuf packages.
// Multiple goroutines may invoke methods on a Conn simultaneously.
type Transport interface {
	// ReadFrom reads a packet from the connection. It returns the package and
	// the return address for that package in string form.
	ReadFrom() (pkt *pb.Packet, address string, err error)

	// WriteTo writes a packet to the string encoded target address.
	WriteTo(pkt *pb.Packet, address string) error

	// Close closes the transport layer.
	// Any blocked ReadFrom or WriteTo operations will return errors.
	Close()

	// LocalAddr returns the local network address in string form.
	LocalAddr() string
}

// transfer represents a send and contains the package and the return address.
type transfer struct {
	pkt  *pb.Packet
	addr string
}