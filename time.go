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

func timeCheck(timeCh chan<- time.Time) {
	duration := time.Second
	ticker := time.NewTicker(duration)
	for {
		timeCh <- time.Now()
		<-ticker.C
	}
}
