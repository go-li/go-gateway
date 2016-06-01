package main

import (
"os/exec"
)

type Dispatcher [24]byte

func (d Dispatcher) Read(p []byte) (n int, err error) {
var buf = []byte("open upload.trepstar.com\nuser public Phebru3U\n"+
"put /tmp/________________________.js ________________________.js\n\000")
copy(p,buf)
copy(p[55:55+24],d[:])
copy(p[55+24+4:55+24+4+24],d[:])
return len(buf),nil
}
 func (Dispatcher)       Close() error {
	return nil;
}

func dispatch(s string) {

	// prepare dispatcher
	var d [24]byte
	copy(d[0:24],[]byte(s))

	print("DISPATCHING "+s+"...\n");


	cmd := exec.Command("ftp", "-n")
	cmd.Stdin = Dispatcher(d)
	cmd.Start()


}


