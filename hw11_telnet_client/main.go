package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	logger := log.New(os.Stderr, "", 0)

	// --timeout=10s host port
	timeout := flag.Duration("timeout", time.Second*10, "connection timeout")
	flag.Parse()
	if flag.NArg() != 2 {
		fmt.Println("use: go-telnet [parameters] host port")
		return
	}
	addr := flag.Arg(0) + ":" + flag.Arg(1)

	client := NewTelnetClient(addr, *timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		logger.Fatalf("connection error: %v", err)
	}
	defer client.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	stopped := false

	// receiver
	go func() {
		if err := client.Receive(); err != nil && !stopped {
			logger.Fatalf("receiver error: %v", err)
		}
		logger.Println("connection closed by peer")
		stopped = true
		stop()
	}()

	// sender
	go func() {
		if err := client.Send(); err != nil && !stopped {
			logger.Fatalf("sender error: %v", err)
		}
		logger.Println("end of input")
		stopped = true
		stop()
	}()

	// wait for ctrl+c
	<-ctx.Done()
}
