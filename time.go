// time.go - clock status updater.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of urc, using the creative
// commons "cc0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"time"
)

const clockLayout = "Mon 02.01 15:04:05"

type Clock struct {
	Time time.Time
}

func (c Clock) Format() string {
	return c.Time.Format(clockLayout)
}

func clockCheck(clockCh chan<- Clock) {
	duration := time.Second
	ticker := time.NewTicker(duration)
	for {
		clock := Clock{
			Time: time.Now(),
		}
		clockCh <- clock
		<-ticker.C
	}
}
