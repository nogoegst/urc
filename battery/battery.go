// battery.go - battery status updater.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of urc, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package battery

import (
	"fmt"
	"strings"
	"time"
)

type Lifetime struct {
	Percent int
	Time    time.Duration
}

func (bs Lifetime) Format() string {
	fb := "no bat"
	if bs.Percent != -1 {
		percent := fmt.Sprintf("%d%%", bs.Percent)
		timeLeft := fmt.Sprintf("%s", strings.TrimRight(bs.Time.String(), "0s"))
		fb = fmt.Sprintf("%s %s", percent, timeLeft)
	}
	return fb
}

func WatchLifetime() <-chan Lifetime {
	ch := make(chan Lifetime)
	go LifetimeCheck(ch)
	return ch
}
