package main

import (
	"os/exec"
	"time"
)

const timeout = 10

type Dispatcher [24]byte

func (d Dispatcher) Read(p []byte) (n int, err error) {
	var buf = []byte("open upload.trepstar.com\nuser public Phebru3U\n" +
		"put /tmp/________________________.js ________________________.js\n\000")
	copy(p, buf)
	copy(p[55:55+24], d[:])
	copy(p[55+24+4:55+24+4+24], d[:])
	return len(buf), nil
}
func (Dispatcher) Close() error {
	return nil
}

func dispatch(s string) error {

	// prepare dispatcher
	var d [24]byte
	copy(d[0:24], []byte(s))

	print("DISPATCHING " + s + "...\n")

	cmd := exec.Command("ftp", "-n")
	cmd.Stdin = Dispatcher(d)
	cmd.Start()

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			return err
		}
		print("process killed as timeout reached\n")
	case err := <-done:
		return err
	}
	return nil
}
