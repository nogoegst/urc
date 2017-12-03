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
	"github.com/nogoegst/urc/torstatus"
)

type Status struct {
	TorStatus       torstatus.TorStatus
	Clock           clock.Clock
	BatteryLifetime battery.Lifetime
	Message         message.Message
}

func (s Status) Format() string {
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
	return Sanitize(status)
}

func Sanitize(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' {
			return ' '
		} else if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, s)
}

func updateStatus(statusChan chan<- string) {
	var status Status

	clockCh := clock.WatchClock()
	messageCh := message.WatchMessages()
	batteryLifetimeCh := battery.WatchLifetime()
	torStatusCh := torstatus.WatchTorStatus()

	for {
		select {
		case v := <-clockCh:
			status.Clock = v
		case v := <-torStatusCh:
			status.TorStatus = v
		case v := <-batteryLifetimeCh:
			status.BatteryLifetime = v
		case v := <-messageCh:
			status.Message = v
		}
		statusChan <- status.Format()
	}

}
