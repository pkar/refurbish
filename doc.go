/*
Package refurbish is a simple binary auto updater. Currently
it only works with accessible http urls.

As an example, on builds a file could be uploaded to S3 along with it's md5 file and
the command to run to auto update.

	package main

	import "github.com/pkar/refurbish"

	func main() {
		r := refurbish.New("http://example.com/path/to/binary", "http://example.com/path/to/md5", "sudo initctl restart binary")
		go r.run()
	}
*/
package refurbish
