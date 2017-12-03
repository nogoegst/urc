// status.go - status maker.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of urc, using the creative
// commons "cc0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"strings"
	"time"
	"unicode"
)

const timeLayout = "Mon 02.01 15:04:05"

type Status struct {
	Time        time.Time
	TorLiveness string
	Battery     BatteryLifetime
	Message     Message
}

func (s *Status) Format() string {
	fMsg := s.Message.Format()
	battery := s.Battery.Format()
	fTorLiveness := "tor is " + strings.ToLower(s.TorLiveness)
	fTime := s.Time.Format(timeLayout)
	cosmoStatus := "Λ > 0"

	status := Compose(fMsg, cosmoStatus, fTorLiveness, battery, fTime)
	return " " + status + " "
}

func Compose(statuses ...string) string {
	status := strings.Join(statuses, " | ")
	return strings.Map(func(r rune) rune {
		if r == '\n' {
			return ' '
		} else if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, status)
}

func updateStatus(statusChan chan<- string) {
	var status Status

	timeCh := make(chan time.Time)
	go timeCheck(timeCh)

	livenessCh := make(chan string)
	go livenessCheck(livenessCh)

	messageCh := make(chan Message)
	go messageBufferedCheck(messageCh, UnixSocketMessageCheck)

	batteryCh := make(chan BatteryLifetime)
	go BatteryCheck(batteryCh)

	for {
		select {
		case time := <-timeCh:
			status.Time = time
		case liveness := <-livenessCh:
			status.TorLiveness = liveness
		case bs := <-batteryCh:
			status.Battery = bs
		case msg := <-messageCh:
			status.Message = msg
		}
		statusChan <- status.Format()
	}

}
