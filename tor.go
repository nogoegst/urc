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
