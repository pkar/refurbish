package refurbish

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// chunk is the chunk size to read while calculating
// MD5.
const (
	chunk = 8192
)

// Refurbish is a self contained struct that figures out
// the currently running executables md5 and updates when
// it detects changes.
type Refurbish struct {
	latestBinURL string // url to the latest binary.
	latestMD5URL string // url to the latest md5 file.
	updateCMD    string // the command to run when updates detected.
	md5          string // the current binary md5 checksum.
	path         string // path the the current running executable.
}

// New creates and runs an interval checking updater.
func New(binURL, md5URL, cmd string) *Refurbish {
	r := &Refurbish{binURL, md5URL, cmd, "", ""}
	return r
}

// run periodically checks the current binary checksum against
// the given remote binary and updates on change.
func (r *Refurbish) Run() {
	md5, err := r.Calc()
	if err != nil {
		return
	}
	r.md5 = md5

	ticker := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-ticker.C:
			remoteMD5, doUpdate, err := r.checkUpdate()
			if err != nil {
				continue
			}
			if doUpdate {
				err := r.update(remoteMD5)
				if err != nil {
					continue
				}
			}
		}
	}
}

// checkUpdate compares the remote md5 against the currently
// running binary and returns a bool whether or not to update.
func (r *Refurbish) checkUpdate() (md5 string, update bool, err error) {
	resp, err := http.Get(r.latestMD5URL)
	if err != nil {
	}
	defer resp.Body.Close()
	md5Bytes, err := ioutil.ReadAll(resp.Body)
	md5 = string(md5Bytes)
	if err == nil && md5 != r.md5 {
		update = true
	}

	return
}

// update downloads and then swaps the binary, then runs
// the reload cmd given.
func (r *Refurbish) update(remoteMD5 string) error {
	out, err := os.Create("/tmp/refurbishedbin")
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(r.latestBinURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	err = os.Rename("/tmp/refurbishedbin", r.path)
	if err != nil {
		return err
	}

	// split the update command on ; and then by space and execute each.
	parts := strings.Split(r.updateCMD, ";")
	for _, cmd := range parts {
		tokens := strings.Split(cmd, " ")
		var err error
		switch len(tokens) {
		case 0:
		case 1:
			_, err = exec.Command(tokens[0]).Output()
		default:
			_, err = exec.Command(tokens[0], tokens[1:]...).Output()
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// Calc gets the current file stats and calculates
// the md5 checksum of it.
func (r *Refurbish) Calc() (string, error) {
	path, err := r.Path()
	if err != nil {
		return "", err
	}

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	r.path = path

	// calculate the file size
	info, err := file.Stat()
	if err != nil {
		return "", err
	}

	md5 := r.MD5(info.Size(), file)
	r.md5 = md5
	return md5, nil
}

// MD5 calculates the MD5 checksum of the given io.Reader.
// Total size must be precomputed.
func (r *Refurbish) MD5(size int64, input io.Reader) string {
	blocks := uint64(math.Ceil(float64(size) / float64(chunk)))

	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(chunk, float64(size-int64(i*chunk))))
		buf := make([]byte, blocksize)

		input.Read(buf)
		io.WriteString(hash, string(buf)) // append into the hash
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Path calculates the path to the currently running executable.
func (r *Refurbish) Path() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	p := dir + "/" + filepath.Base(os.Args[0])
	return p, nil
}
