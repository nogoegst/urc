// urc.go - stuff displayed in upper right corner.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of urc, using the creative
// commons "cc0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"log"
	"fmt"
	"net"
	"os"
	"time"
	"strings"

	"github.com/nogoegst/bulb"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

const timeLayout = "Mon 02.01 15:04:05"
const torReconnectDelay = 5 * time.Second
const messageTTL = 30 * time.Minute
const maxMsgLength = 64

type Status struct {
	Time		time.Time
	TorLiveness	string
	Message		string
	MessageTimestamp	time.Time
}

func (s *Status) Format() (string) {
	fMsg := strings.TrimRight(s.Message, "\n\r")
	if len(fMsg) > maxMsgLength {
		fMsg = fMsg[:maxMsgLength] + "[...]"
	}
	if fMsg != "" {
		fMsg += fmt.Sprintf(" %dm", int(time.Since(s.MessageTimestamp).Minutes()))
	}
	fTorLiveness := strings.ToLower(s.TorLiveness)
	fTime := s.Time.Format(timeLayout)
	return fmt.Sprintf("%s | Λ > 0 | tor is %s | %s ", fMsg, fTorLiveness, fTime)
}

func messageCheck(messageCh chan<- string) {
	sockpath := os.Getenv("HOME")+"/urc.sock"
	os.Remove(sockpath)
	l, err := net.Listen("unix", sockpath)
	if err != nil {
		log.Printf("Unable to listen on socket: %v", err)
		close(messageCh)
		return
	}
	defer os.Remove(sockpath)

	for {
		c, err := l.Accept()
		if err != nil {
			log.Printf("Unable to accept connection: %v", err)
			continue
		}
		buf := make([]byte, 255)
		n, err := c.Read(buf)
		if err != nil {
			log.Printf("Unable to read from connection: %v", err)
			continue
		}
		c.Close()
		messageCh <- string(buf[:n])
	}
}

func livenessCheck(livenessCh chan<- string) {
	for {
		c, err := bulb.DialURL("default://")
		if err != nil {
			log.Printf("Failed to connect to control socket: %v", err)
			livenessCh <- "not running"
			time.Sleep(torReconnectDelay)
			continue
		}
		defer c.Close()
		if err := c.Authenticate("ExamplePassword"); err != nil {
			log.Printf("Authentication failed: %v", err)
			close(livenessCh)
			return
		}
		livenessCh <- "running"
		c.StartAsyncReader()
		resp, err := c.Request("GETINFO network-liveness")
		if err != nil {
			log.Fatalf("GETINFO failed: %v", err)
		}
		livenessCh <- strings.TrimPrefix(resp.Data[0], "network-liveness=")
		if _, err := c.Request("SETEVENTS NETWORK_LIVENESS"); err != nil {
			log.Fatalf("SETEVENTS NETWORK_LIVENESS has failed: %v", err)
		}
		for {
			ev, err := c.NextEvent()
			if err != nil {
				break
			}
			livenessCh <- strings.TrimPrefix(ev.Reply, "NETWORK_LIVENESS ")
		}
		livenessCh <- "not running"
		time.Sleep(torReconnectDelay)
	}
}

func timeCheck(timeCh chan<- time.Time) {
	duration := time.Second
	ticker := time.NewTicker(duration)
	for {
		timeCh <- time.Now()
		<-ticker.C
	}
}

func updateStatus(statusChan chan<- string) {
	var status Status

	timeCh := make(chan time.Time)
	go timeCheck(timeCh)

	livenessCh := make(chan string)
	go livenessCheck(livenessCh)

	messageCh := make(chan string)
	go messageCheck(messageCh)
	messageTimer := time.NewTimer(time.Duration(0))
	messageTicker := time.NewTicker(time.Minute)
	messageTicker.Stop()

	for {
		select {
		case time := <-timeCh:
			status.Time = time
		case liveness := <-livenessCh:
			status.TorLiveness = liveness
		case msg := <-messageCh:
			messageTimer.Reset(messageTTL)
			messageTicker.Stop()
			messageTicker = time.NewTicker(time.Minute)
			status.Message = msg
			status.MessageTimestamp = time.Now()
		case <-messageTicker.C:
		case <-messageTimer.C:
			status.Message = ""
			messageTicker.Stop()
		}
		statusChan <- status.Format()
	}

}

func setStatus(x *xgb.Conn, status string) {
	root := xproto.Setup(x).DefaultScreen(x).Root
	xproto.ChangeProperty(x, xproto.PropModeReplace, root, xproto.AtomWmName, xproto.AtomString, 8, uint32(len(status)), []byte(status))
}

func main() {
	x, err := xgb.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	defer x.Close()

	defer setStatus(x, "urc died (-__-,) ")
	statusChan := make(chan string)
	go updateStatus(statusChan)
	for status := range statusChan {
		setStatus(x, status)
	}
}
