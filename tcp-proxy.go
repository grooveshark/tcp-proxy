package main

import (
	"io"
	"log"
	"net"

	"github.com/mediocregopher/flagconfig"
)

var Local, Remote string

func CopyClose(dst io.WriteCloser, src io.ReadCloser) {
	io.Copy(dst, src)
	dst.Close()
	src.Close()
}

func main() {
	fc := flagconfig.New("tcp-proxy")
	fc.DisallowConfig()
	fc.StrParam("local", "Address to listen on", ":4444")
	fc.RequiredStrParam("remote", "Address to proxy to")

	if err := fc.Parse(); err != nil {
		log.Fatal(err)
	}

	Local, Remote = fc.GetStr("local"), fc.GetStr("remote")

	l, err := net.Listen("tcp", Local)
	if err != nil {
		log.Fatal(err)
	}

	for {
		lconn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		rconn, err := net.Dial("tcp", Remote)
		if err != nil {
			log.Print(err)
			lconn.Close()
			continue
		}

		go CopyClose(lconn, rconn)
		go CopyClose(rconn, lconn)
	}
}
