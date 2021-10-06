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

func main() {
	// default host
	host := "127.0.0.1"

	// default port
	port := "4242"

	// default timeout
	timeout := 10 * time.Second

	flagTimeout := flag.Duration("timeout", 10*time.Second, "connect timeout")

	flag.Parse()
	if h := flag.Arg(0); h != "" {
		host = h
	}
	if p := flag.Arg(1); p != "" {
		port = p
	}

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
