package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

type publisher struct {
	Name    string
	Subject string
}

func (p *publisher) send(m *message, econn *nats.EncodedConn, w io.Writer) error {
	enc := json.NewEncoder(w)
	econn.Publish(p.Subject, m.toSend(time.Now()))
	econn.Flush()
	if err := econn.LastError(); err != nil {
		return fmt.Errorf("Error publishing a mesage: %s, error:%s", m, err)
	}
	if err := enc.Encode(m); err != nil {
		return err
	}

	log.Printf("Message %q \nPublished at %s", m, m.SendAt)
	return nil
}

func (p *publisher) sendNextMessage(m *message, conn *nats.EncodedConn, w io.Writer) error {
	for {
		m = m.next()
		if err := p.send(m, conn, w); err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}
func (p *publisher) sendMessages(mm []*message, conn *nats.EncodedConn, w io.Writer) error {
	for _, m := range mm {
		if err := p.send(m, conn, w); err != nil {
			return err
		}
	}
	return nil
}
func (p *publisher) sendMessagesFromFile(file string, conn interface{}) error {
	var econn *nats.EncodedConn
	if v, ok := conn.(*nats.EncodedConn); ok {
		econn = v
	} else if v, ok := conn.(*nats.Conn); ok {
		ec, err := nats.NewEncodedConn(v, nats.JSON_ENCODER)
		if err != nil {
			return fmt.Errorf("Error gerimg nats json encoded connection %s ", err)
		}
		econn = ec
	} else {
		return fmt.Errorf("Unnoun connection provider %t ", econn)
	}
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatal("Unable to open file %s, err: %s", file, err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var m *message
	m = m.new()
	mm := make([]*message, 0, 1)
	for scanner.Scan() {
		if err := json.Unmarshal([]byte(scanner.Text()), m); err != nil {
			return fmt.Errorf("Fail to decode messages from %s, with error: %s", file, err)
		}
		mm = append(mm, m)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if len(mm) == 0 {
		if err := p.send(m, econn, f); err != nil {
			return err
		}
	} else {
		m = mm[len(mm)-1]
	}
	if err := p.sendNextMessage(m, econn, f); err != nil {
		return err
	}
	return nil
}
