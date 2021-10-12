package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout, default is 10s")
	flag.Parse()
	host := flag.Arg(0)
	port := flag.Arg(1)
	if host == "" || port == "" {
		log.Fatal("host and port is required params")
	}
	addr := net.JoinHostPort(host, port)

	c := NewTelnetClient(addr, *timeout, os.Stdin, os.Stdout)
	if err := c.Connect(); err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	io.WriteString(os.Stderr, fmt.Sprintf("...connected to: %s\n", addr))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		err := c.Receive()
		if err != nil {
			io.WriteString(os.Stderr, fmt.Sprintf("...%s\n", err))
		}
		cancel()
	}()

	go func() {
		err := c.Send()
		if err != nil {
			io.WriteString(os.Stderr, fmt.Sprintf("...%s\n", err))
		}
		cancel()
	}()

	<-ctx.Done()
}
