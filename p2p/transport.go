package p2p

import "net"

// Peer is an interface that represents remote node
type Peer interface {
	net.Conn
	send([]byte) error
}

// Transport is anything that handles the communication between nodes and the network
// This can be of form of  tcp, udp, websockets.
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
