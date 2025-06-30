package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"murali.dfss/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {

	tcpTranspotOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}

	tcpTransport := p2p.NewTCPTransport(tcpTranspotOpts)

	fileServerOPts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTRansformFunc,
		Transport:         tcpTransport,
		BootStrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOPts)

	tcpTransport.OnPeer = s.onPeer

	return s

}

func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")
	go func() {

		log.Fatal(s1.Start())

	}()

	time.Sleep(2 * time.Second)

	go s2.Start()
	time.Sleep(2 * time.Second)

	for i := 0; i < 10; i++ {
		data := bytes.NewReader([]byte("my big data file here!"))
		s2.Store(fmt.Sprintf("myprivatedata_%d", i), data)
		time.Sleep(5 * time.Millisecond)
	}

	// r, err := s1.Get("myprivatedata")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// b, err := io.ReadAll(r)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(b))

	select {}

}
