// status.go - status maker.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of urc, using the creative
// commons "cc0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"fmt"
	"strings"
	"time"
	"unicode"
)

const timeLayout = "Mon 02.01 15:04:05"
const maxMsgLength = 64

type Status struct {
	Time        time.Time
	TorLiveness string
	Battery     batteryLifetime
	Message     Message
}

func (s *Status) Format() string {
	fMsg := strings.TrimRight(s.Message.Text, "\n\r")
	if len(fMsg) > maxMsgLength {
		fMsg = fMsg[:maxMsgLength] + "[...]"
	}
	if fMsg != "" {
		fMsg += fmt.Sprintf(" %dm", int(time.Since(s.Message.Timestamp).Minutes()))
	}
	fBattery := "no bat"
	if s.Battery.Percent != -1 {
		fBattery = fmt.Sprintf("%d%% %s", s.Battery.Percent, strings.TrimRight(s.Battery.Time.String(), "0s"))
	}
	fTorLiveness := strings.ToLower(s.TorLiveness)
	fTime := s.Time.Format(timeLayout)
	status := fmt.Sprintf(" %s | Λ > 0 | tor is %s | %s | %s ", fMsg, fTorLiveness, fBattery, fTime)
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

	batteryCh := make(chan batteryLifetime)
	go batteryCheck(batteryCh)

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
