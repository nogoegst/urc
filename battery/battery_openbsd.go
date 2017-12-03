// battery_openbsd.go - battery status updater for OpenBSD.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of urc, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package battery

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func LifetimeCheck(ch chan<- Lifetime) {
	ticker := time.NewTicker(10 * time.Second)
	for {
		bs := Lifetime{Percent: -1}
		out, err := exec.Command("apm", "-l", "-m").Output()
		if err != nil {
			log.Fatal(err)
		}
		split := strings.Split(strings.TrimRight(string(out), "\n"), "\n")
		if len(split) != 2 {
			log.Fatalf("Corrupted apm(8) output")
		}
		percent, err := strconv.Atoi(split[0])
		if err != nil {
			log.Fatalf("Corrupted apm(8) output")
		}
		bs.Percent = percent
		minutes, err := time.ParseDuration(split[1] + "m")
		if err != nil && split[1] != "unknown" {
			log.Fatalf("Corrupted apm(8) output")
		}
		bs.Time = minutes
		ch <- bs
		<-ticker.C
	}
}
