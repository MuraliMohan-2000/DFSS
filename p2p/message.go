package p2p

const (
	IncommingMessage = 0x1
	IncommingStream  = 0x2
)

// RPC represents a any arbitrary that is being sent over each
// transport between two nodes in the network
type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}
