package main

import (
	"log"

	"murali.dfss/p2p"
)

func main() {
	tr := p2p.NewTCPTransport(":3000")
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}

}

func NOPHandshakeFunc(any) error {
	return nil
}
