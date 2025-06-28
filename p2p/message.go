package p2p

import "net"

// RPC represents a any arbitrary that is being sent over each
// transport between two nodes in the network
type RPC struct {
	From    net.Addr
	Payload []byte
}
