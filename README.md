urc
===========
Makes stuff be displayed in upper right corner of DWM.
And does it asyncronously.

Currently only time/date, tor liveness, state of the Universe are implemented.

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
$GOPATH/bin/urc
.
w
q
```
That's it.
