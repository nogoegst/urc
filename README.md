urc
===========
Makes stuff be displayed in upper right corner of DWM.
And does it asyncronously.

Currently only time/date, tor liveness, incoming strings on unix socket, battery status (OpenBSD) and state of the Universe are implemented.

Install
-------
```
$ go get github.com/nogoegst/urc
```

Usage
-----
```
$ ed ~/.xinitrc
1i
$GOPATH/bin/urc &
.
w
q
```

Send messages to `$HOME/urc.sock`. `urc` is at service!

That's it.

Hacking
-------
It's easy to change your status line. Go ahead to `status.go` and change formatting, active modules or implement more of them!
