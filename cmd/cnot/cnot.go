// cnot.go - utility to clear notifications in urc.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to this module of urc, using the creative
// commons "CC0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package main

import (
	"log"
	"net"
	"os/user"
	"path/filepath"
)

func clearNotifications() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	urcSocket := filepath.Join(usr.HomeDir, "urc.sock")
	c, err := net.Dial("unix", urcSocket)
	if err != nil {
		return err
	}
	defer c.Close()
	_, err = c.Write([]byte("\n"))
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if err := clearNotifications(); err != nil {
		log.Fatal(err)
	}
}
