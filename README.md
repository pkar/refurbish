# refurbish
[![wercker status](https://app.wercker.com/status/810483e0fb7e4e6526272e827d7d5ef9/m "wercker status")](https://app.wercker.com/project/bykey/810483e0fb7e4e6526272e827d7d5ef9)

A simple binary auto updater. Currently it only works with accessible http urls.

```go
	package main

	import "github.com/pkar/refurbish"

	func main() {
		r := refurbish.New("http://example.com/path/to/binary", "http://example.com/path/to/md5", "sudo initctl restart binary")
		go r.run()
	}
```

#### todo
- alot
