// battery.go - battery status updater.
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
)

type BatteryLifetime struct {
	Percent int
	Time    time.Duration
}

func (bs *BatteryLifetime) Format() string {
	fb := "no bat"
	if bs.Percent != -1 {
		fb = fmt.Sprintf("%d%% %s", bs.Percent, strings.TrimRight(bs.Time.String(), "0s"))
	}
	return fb
}
