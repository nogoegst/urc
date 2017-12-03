// status.go - status maker.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of urc, using the creative
// commons "cc0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"strings"
	"unicode"
)

type Status struct {
	Clock     Clock
	TorStatus TorStatus
	Battery   BatteryLifetime
	Message   Message
}

func (s *Status) Format() string {
	msg := s.Message.Format()
	cosmoStatus := "Λ > 0"
	torStatus := s.TorStatus.Format()
	battery := s.Battery.Format()
	clockStatus := s.Clock.Format()

	status := Compose(msg, cosmoStatus, torStatus, battery, clockStatus)
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

	clockCh := make(chan Clock)
	go clockCheck(clockCh)

	torstatusCh := make(chan TorStatus)
	go torstatusCheck(torstatusCh)

	messageCh := make(chan Message)
	go messageBufferedCheck(messageCh, UnixSocketMessageCheck)

	batteryCh := make(chan BatteryLifetime)
	go BatteryCheck(batteryCh)

	for {
		select {
		case clock := <-clockCh:
			status.Clock = clock
		case torstatus := <-torstatusCh:
			status.TorStatus = torstatus
		case bs := <-batteryCh:
			status.Battery = bs
		case msg := <-messageCh:
			status.Message = msg
		}
		statusChan <- status.Format()
	}

}
