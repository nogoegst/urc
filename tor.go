// tor.go - tor-related status updater.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of urc, using the creative
// commons "cc0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"log"
	"strings"
	"time"

	"github.com/nogoegst/bulb"
)

const torReconnectDelay = 5 * time.Second

type TorStatus struct {
	Liveness string
}

func (ts TorStatus) Format() string {
	return "tor is " + ts.Liveness
}

func torstatusCheck(torstatusCh chan<- TorStatus) {
	defer close(torstatusCh)
	for {
		ts := TorStatus{}
		c, err := bulb.DialURL("default://")
		if err != nil {
			log.Printf("Failed to connect to control socket: %v", err)
			ts.Liveness = "not running"
			torstatusCh <- ts
			time.Sleep(torReconnectDelay)
			continue
		}
		defer c.Close()
		if err := c.Authenticate("ExamplePassword"); err != nil {
			log.Printf("Authentication failed: %v", err)
			return
		}
		ts.Liveness = "running"
		torstatusCh <- ts
		c.StartAsyncReader()
		resp, err := c.Request("GETINFO network-liveness")
		if err != nil {
			log.Printf("GETINFO failed: %v", err)
			return
		}

		ts.Liveness = strings.TrimPrefix(resp.Data[0], "network-liveness=")
		torstatusCh <- ts
		if _, err := c.Request("SETEVENTS NETWORK_LIVENESS"); err != nil {
			log.Printf("SETEVENTS NETWORK_LIVENESS has failed: %v", err)
			return
		}
		for {
			ev, err := c.NextEvent()
			if err != nil {
				log.Printf("NextEvent error: %v", err)
				break
			}
			ts.Liveness = strings.TrimPrefix(ev.Reply, "NETWORK_LIVENESS ")
			torstatusCh <- ts
		}
		ts.Liveness = "not running"
		torstatusCh <- ts
		time.Sleep(torReconnectDelay)
	}
}
