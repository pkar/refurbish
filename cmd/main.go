package main

import "github.com/pkar/refurbish"

func main() {
	r := refurbish.New("http://example.com/path/to/binary", "http://example.com/path/to/md5", "sudo initctl restart binary")
	go r.Run()
}
