// main.go - get status and set it to upper right corner.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of urc, using the creative
// commons "cc0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"log"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

func setStatus(x *xgb.Conn, status string) {
	root := xproto.Setup(x).DefaultScreen(x).Root
	xproto.ChangeProperty(x, xproto.PropModeReplace, root, xproto.AtomWmName, xproto.AtomString, 8, uint32(len(status)), []byte(status))
}

func main() {
	x, err := xgb.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	defer x.Close()

	defer setStatus(x, "urc died")
	statusChan := make(chan string)
	go UpdateStatus(statusChan)
	for status := range statusChan {
		setStatus(x, status)
	}
}
