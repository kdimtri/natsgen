package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/nats-io/nats.go"
)

func usage() {
	fmt.Println("Usage:\n./nats [ pub | sub ] [ pubFile.json | subFile.json ] \nExample:\n./nats pub ./pub.json\n./nats sub ./sub.json\nOptional flags:\n")
	flag.PrintDefaults()
}

func main() {
	var url = flag.String("s", "demo.nats.io", "NATS server URL")
	var subj = flag.String("j", "default", "Subject to publish to/subscribe on")
	var delay = flag.Int("d", 1, "delay betwin messages (in seconds)")
	var wait = flag.Int("w", 60, "For how long to wait for reconections(in seconds)")
	var help = flag.Bool("h", false, "Show help message")
	flag.Parse()
	if *help {
		usage()
		os.Exit(0)
	}
	args := flag.Args()
	if len(args) != 2 {
		usage()
		os.Exit(1)
	}
	w, file := args[0], args[1]
	wt, dt := time.Duration(*wait)*time.Second, time.Duration(*delay)*time.Second
	if w == "pub" {
		name := "nats file message publisher"
		opts := setOpts(name, wt, dt)
		conn, err := nats.Connect(*url, opts...)
		if err != nil {
			log.Fatal("Unable to conntct to Nats server")
		}
		pub := &publisher{name, *subj}
		go log.Fatal(pub.sendMessagesFromFile(file, conn))
	} else if w == "sub" {
		name := "nats file message subscriber"
		opts := setOpts(name, wt, dt)
		conn, err := nats.Connect(*url, opts...)
		if err != nil {
			log.Fatal("Unable to conntct to Nats server")
		}
		sub := &subscriber{name, *subj}
		sub.SubscribeWithFile(file, conn)
	} else {
		fmt.Println("Undefined nats worker", args[0])
		usage()
		os.Exit(2)
	}
	runtime.Goexit()
}

func setOpts(name string, wait time.Duration, delay time.Duration) []nats.Option {
	return []nats.Option{
		nats.Name(name),
		nats.ReconnectWait(wait),
		nats.MaxReconnects(int(wait / delay)),
		nats.ReconnectHandler(reconnected),
		nats.DisconnectErrHandler(disconnected),
		nats.ClosedHandler(closed),
	}
}

func disconnected(conn *nats.Conn, err error) {
	log.Printf("%s disconected with err: %s", conn.Opts.Name, err)
}
func reconnected(conn *nats.Conn) {
	log.Printf("%s reconected to ", conn.Opts.Name, conn.ConnectedUrl())
}
func closed(conn *nats.Conn) {
	err := conn.LastError()
	if err != nil {
		log.Fatal(fmt.Sprintf("Server %s close %s connection with an err: %s", conn.ConnectedUrl(), conn.Opts.Name, err), err)
	}
	log.Fatal("Server %s closed %s connection")
}
