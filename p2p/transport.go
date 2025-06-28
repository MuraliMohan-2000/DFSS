package p2p

// Peer is an interface that represents remote node
type Peer interface{}

// Transport is anything that handles the communication between nodes and the network
// This can be of form of  tcp, udp, websockets.
type Transport interface {
	ListenAndAccept() error
}
