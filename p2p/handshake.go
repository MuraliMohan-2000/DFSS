package p2p

// Handshakefunc ....?
type HandshakeFunc func(any) error

func NOPHandshakeFunc(any) error {
	return nil
}
