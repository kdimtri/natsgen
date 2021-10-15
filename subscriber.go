package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

type subscriber struct {
	Name    string
	Subject string
}

func (s *subscriber) SubscribeWithFile(file string, conn *nats.Conn) error {
	econn, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		return fmt.Errorf("Error gerimg nats json encoded connection %s ", err)
	}
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatal("Unable to open file %s, err: %s", file, err)
	}
	econn.Subscribe(s.Subject, func(msg *nats.Msg) {
		mmap := make(map[string]interface{})
		if err := json.Unmarshal(msg.Data, &mmap); err != nil {
			log.Fatal(err)

		}
		rt := time.Now().Format(time.Stamp)
		log.Printf("Mesage: %v\nReceived at %s", mmap, rt)
		mmap["rcvTime"] = rt
		if _, err := f.Write([]byte(fmt.Sprintln(mmap))); err != nil {
			log.Fatal(err)
		}
	})
	econn.Flush()
	if err := econn.LastError(); err != nil {
		return (err)
	}
	log.Printf("Waiting for messages on %s", s.Subject)
	return nil
}
