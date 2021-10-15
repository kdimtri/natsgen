package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"time"

	lorem "github.com/drhodes/golorem"
)

type message struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Hash   string `json:"hash"`
	SendAt string `json:"sendAt"`
}

func (m *message) String() string {
	return fmt.Sprintf("Message( %d ): %s, hash: %s", m.ID, m.Text, m.Hash)
}
func (m *message) bytes() (data []byte) {
	if m != nil && m.ID != 0 {
		_ = json.Unmarshal(data, *m)
	}
	return data
}
func (m *message) hash(b []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(b))
}
func (m *message) new() *message {
	return &message{
		ID:     0,
		Text:   text(),
		SendAt: "",
		Hash:   fmt.Sprintf("%x", sha1.Sum([]byte(""))),
	}
}
func (m *message) next() *message {
	nm := &message{
		ID:     m.ID + 1,
		Text:   text(),
		SendAt: "",
		Hash:   fmt.Sprintf("%x", sha1.Sum([]byte(m.Hash+string(m.bytes())))),
	}
	return nm
}

func (m *message) toSend(t time.Time) *message {
	m.SendAt = t.Format(time.Stamp)
	return m
}
func text() string {
	return lorem.Sentence(2, 3)
}
