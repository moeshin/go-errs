package errs

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestPrint(t *testing.T) {
	log.Println(Print(errors.New("TestPrint")))
}

func TestPrintWithDepthToLogBuffer(t *testing.T) {
	buf := PrintWithDepthToLogBuffer(errors.New("TestPrintWithDepthToLogBuffer"), 0)
	_, _ = log.Writer().Write(buf.Bytes())
}

func TestClose(t *testing.T) {
	fp, err := os.Open("LICENSE")
	Panic(err)
	defer Close(fp)
	hash := sha256.New()
	_, err = io.Copy(hash, fp)
	Panic(err)
	s := hex.EncodeToString(hash.Sum(nil))
	t.Log("sha256", s)
	if s != "c0ce2cbd8203985284580cf1ef071bba769e60584c965e3381afe342a3851ef4" {
		t.Fail()
	}
}

func TestCloseResponse(t *testing.T) {
	resp, err := http.Head("https://www.baidu.com/")
	Panic(err)
	defer CloseResponse(resp)
	t.Log("ContentLength", resp.ContentLength)
}

func TestDefer(t *testing.T) {
	defer Defer(func() error {
		return errors.New("TestDefer")
	})
}

func TestDeferCall(t *testing.T) {
	defer DeferCall(errors.New, "TestDeferCall")
}

func TestDeferCallSlice(t *testing.T) {
	defer DeferCallSlice(func(arr ...string) error {
		if len(arr) == 0 {
			return nil
		}
		return errors.New(arr[0])
	}, []string{"TestDeferCall"})
}
