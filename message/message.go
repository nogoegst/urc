// message.go - message reciever and ticker.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of urc, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package message

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const messageTTL = 30 * time.Minute
const MaxMessageSize = 255

func UnixSocketMessageCheck(ch chan<- string) {
	defer close(ch)
	sockpath := filepath.Join(os.Getenv("HOME"), "urc.sock")
	// Remove socket left from dead urc
	os.Remove(sockpath)
	l, err := net.Listen("unix", sockpath)
	if err != nil {
		log.Printf("Unable to listen on socket: %v", err)
		return
	}
	defer os.Remove(sockpath)

	for {
		c, err := l.Accept()
		if err != nil {
			log.Printf("Unable to accept connection: %v", err)
			continue
		}
		buf := make([]byte, MaxMessageSize)
		n, err := c.Read(buf)
		if err != nil {
			log.Printf("Unable to read from connection: %v", err)
			continue
		}
		c.Close()
		ch <- string(buf[:n])
	}
}

const maxMsgLength = 64

type Message struct {
	Text      string
	Timestamp time.Time
}

func (m *Message) Format() string {
	fm := strings.TrimRight(m.Text, "\n\r")
	if len(fm) > maxMsgLength {
		fm = fm[:maxMsgLength] + "[...]"
	}
	if fm != "" {
		fm += fmt.Sprintf(" %dm", int(time.Since(m.Timestamp).Minutes()))
	}
	return fm
}

func MessageBufferedCheck(out chan<- Message, mchk func(chan<- string)) {
	messageCh := make(chan string)
	go mchk(messageCh)
	messageTimer := time.NewTimer(time.Duration(0))
	messageTicker := time.NewTicker(time.Minute)
	messageTicker.Stop()
	for {
		select {
		case msg := <-messageCh:
			messageTimer.Reset(messageTTL)
			messageTicker.Stop()
			messageTicker = time.NewTicker(time.Minute)
			out <- Message{Text: msg, Timestamp: time.Now()}
		case <-messageTicker.C:
		case <-messageTimer.C:
			messageTicker.Stop()
			out <- Message{}
		}
	}
}

func WatchMessages() <-chan Message {
	ch := make(chan Message)
	go MessageBufferedCheck(ch, UnixSocketMessageCheck)
	return ch
}
