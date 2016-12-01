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
	"time"
	"strings"

	"github.com/nogoegst/bulb"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

const timeLayout = "Mon 02.01 15:04:05"
const torReconnectDelay = 5 * time.Second

type Status struct {
	Time	time.Time
	TorLiveness	string
}

func (s *Status) Format() (string) {
	fTorLiveness := strings.ToLower(s.TorLiveness)
	fTime := s.Time.Format(timeLayout)
	return fmt.Sprintf("Λ > 0 | tor is %s | %s ", fTorLiveness, fTime)
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
		<-ticker.C
		timeCh <- time.Now()
	}
}

func updateStatus(statusChan chan<- string) {
	var status Status

	timeCh := make(chan time.Time)
	go timeCheck(timeCh)

	livenessCh := make(chan string)
	go livenessCheck(livenessCh)

	for {
		select {
		case time := <-timeCh:
			status.Time = time
		case liveness := <-livenessCh:
			status.TorLiveness = liveness
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
