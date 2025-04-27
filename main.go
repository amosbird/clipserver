package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/atotto/clipboard"
)

var socket = flag.String("s", fmt.Sprintf("%s/clipservice.sock", os.TempDir()), "domain socket file")
var server = flag.Bool("S", false, "is server?")

func serve(c net.Conn) {
	data, err := ioutil.ReadAll(c)
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := clipboard.WriteAll(string(data)); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	if !*server {
		c, err := net.Dial("unix", *socket)
		if err != nil {
			log.Fatal("Dial error", err)
		}
		defer c.Close()
		_, err = io.Copy(c, os.Stdin)
		if err != nil {
			log.Fatal("Write error:", err)
		}
		return
	}

	clipboard.Primary = false

	syscall.Unlink(*socket)
	ln, err := net.Listen("unix", *socket)
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(ln net.Listener, c chan os.Signal) {
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		ln.Close()
		os.Exit(0)
	}(ln, sigc)

	for {
		fd, err := ln.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}
		go serve(fd)
	}
}
