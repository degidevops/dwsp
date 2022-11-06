package main

import (
	"flag"
	"io"
	"log"
	"net"
	"sync"
	"time"
)


type libStr struct {
	address, dstAddress,
	handshakeCode *string
}

type ssh interface {
	Write([]byte) (int, error)
	Read([]byte) (int, error)
	Close() error
}

func main() {
	var wg sync.WaitGroup
	SSH := libStr{
    address:       flag.String("l", "localhost:8080", "Set port listen. Ex.: localhost:8080"),
    dstAddress:    flag.String("d", "localhost:22", "Set internal ssh address. Ex.: localhost:22"),
    handshakeCode: flag.String("hsc", "101 Upgrade Protocol", "Set custom respon code. Ex.: 101/200.. etc."),
	}
	flag.Parse()
	go SSH.Svr(&wg)

	wg.Add(1)
	wg.Wait()
}



func (p *libStr) Svr(wg *sync.WaitGroup) {
	defer wg.Done()
	taddr, errNet := net.ResolveTCPAddr("tcp", *p.address)
	if errNet != nil {
		log.Println(errNet.Error())
		return
	}
	conn, errConn := net.ListenTCP("tcp", taddr)

	if errConn != nil {
		log.Fatal("Failed to listen HTTP server ", errConn.Error())
	}
	log.Println("HTTP Server listening on ", *p.address, "redirect to -> ", *p.dstAddress)
	for {
		ClientConn, err := conn.AcceptTCP()
		if err != nil {
			log.Println("Failed to accepted stream on HTTP mode: ", err.Error())
			continue
		}
		err = ClientConn.SetKeepAlive(true)
		if err != nil {
			log.Println("Failed to set keepalive: ", err.Error())
			continue
		}
		go p.Handler(ClientConn)
	}
}



func (p *libStr) Handler(ClientConn ssh) {
	if len(*p.handshakeCode) > 0 {
		_, err := ClientConn.Write([]byte("HTTP/1.1 " + *p.handshakeCode+ "\r\n\r\n"))
		if err != nil {
			return
		}
	} else {
		_, err := ClientConn.Write([]byte("HTTP/1.1 101 Upgrade Protocol\r\n\r\n"))
		if err != nil {
			return
		}
	}
	sshConn, err := net.DialTimeout("tcp", *p.dstAddress, 15*time.Second)
	if err != nil {
		log.Println("Failed to call destination. ", err.Error())
		return
	}

		if p.discardPayload(ClientConn) == nil {
			go copyStream(sshConn, ClientConn)
			go copyStream(ClientConn, sshConn)
		} else {
			log.Println("Failed on receive payload", p.discardPayload(ClientConn))
      return
		}
	}

func (p *libStr) discardPayload(ClientConn ssh) error {
	bft := make([]byte, 2048)
	_, err := io.ReadAtLeast(ClientConn, bft, 5)

	if err != nil {
    log.Println(err)
		return err
	} else {
		return nil
	}
}

func copyStream(input, output ssh) {
	_, err := io.Copy(input, output)
	if err != nil {
		return
	}
}


