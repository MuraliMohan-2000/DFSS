package main

import (
	"log"
	"time"

	"murali.dfss/p2p"
)

func main() {

	tcpTranspotOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		//ToDO: onPeer func
	}

	tcpTranspot := p2p.NewTCPTransport(tcpTranspotOpts)

	fileServerOPts := FileServerOpts{
		StorageRoot:       "3000_network",
		PathTransformFunc: CASPathTRansformFunc,
		Transport:         tcpTranspot,
	}

	s := NewFileServer(fileServerOPts)

	go func() {
		time.Sleep(time.Second * 3)
		s.Stop()
	}()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
