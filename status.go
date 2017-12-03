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

	"github.com/nogoegst/urc/battery"
	"github.com/nogoegst/urc/clock"
	"github.com/nogoegst/urc/message"
)

type Status struct {
	TorStatus       TorStatus
	Clock           clock.Clock
	BatteryLifetime battery.Lifetime
	Message         message.Message
}

func (s *Status) Format() string {
	msg := s.Message.Format()
	cosmoStatus := "Λ > 0"
	torStatus := s.TorStatus.Format()
	batteryLifetime := s.BatteryLifetime.Format()
	clockStatus := s.Clock.Format()

	status := Compose(msg, cosmoStatus, torStatus, batteryLifetime, clockStatus)
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

	torstatusCh := make(chan TorStatus)
	go torstatusCheck(torstatusCh)

	clockCh := clock.WatchClock()
	messageCh := message.WatchMessages()
	batteryLifetimeCh := battery.WatchLifetime()

	for {
		select {
		case clock := <-clockCh:
			status.Clock = clock
		case torstatus := <-torstatusCh:
			status.TorStatus = torstatus
		case v := <-batteryLifetimeCh:
			status.BatteryLifetime = v
		case v := <-messageCh:
			status.Message = v
		}
		statusChan <- status.Format()
	}

}
