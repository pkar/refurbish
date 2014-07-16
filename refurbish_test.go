package refurbish

import (
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	r := New("binurl", "md5url", "update cmd")
	if r == nil {
		t.Fatal("no refurbish returned")
	}
}

func TestCalc(t *testing.T) {
	r := New("binurl", "md5url", "update cmd")
	m, err := r.Calc()
	if err != nil {
		t.Fatal(err)
	}
	if m == "" {
		t.Fatal("md5 not calculated")
	}
	if r.md5 == "" {
		t.Fatal("md5 not calculated")
	}
	if r.path == "" {
		t.Fatal("path not calculated")
	}
}

func TestMD5(t *testing.T) {
	in := strings.NewReader("a")
	r := New("binurl", "md5url", "update cmd")
	md5 := r.MD5(1, in)
	if md5 == "" {
		t.Fatal("md5 not calculated")
	}
}

func TestPath(t *testing.T) {
	r := New("binurl", "md5url", "update cmd")
	p, err := r.Path()
	if err != nil {
		t.Fatal(err)
	}
	if p == "" {
		t.Fatal("path not calculated")
	}
}
