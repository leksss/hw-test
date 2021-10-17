package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	DefaultHost    = "127.0.0.1"
	DefaultPort    = "4242"
	DefaultTimeout = 10 * time.Second
)

func main() {
	flagTimeout := flag.Duration("timeout", 10*time.Second, "connect timeout")

	flag.Parse()
	host := DefaultHost
	if h := flag.Arg(0); h != "" {
		host = h
	}
	port := DefaultPort
	if p := flag.Arg(1); p != "" {
		port = p
	}
	timeout := DefaultTimeout
	if flagTimeout != nil {
		timeout = *flagTimeout
	}

	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	defer func() {
		if err := client.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		if err := client.Send(); err != nil {
			log.Fatalln(err)
		}
		cancel()
	}()

	go func() {
		if err := client.Receive(); err != nil {
			log.Fatalln(err)
		}
		cancel()
	}()

	<-ctx.Done()
}
