// battery_other.go - fake battery status updater.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of urc, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

// +build !openbsd

package main

func LifetimeCheck(ch chan<- Lifetime) {
	ch <- Lifetime{Percent: -1}
}
