package errs

import (
	"errors"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var TestErr = errors.New("test")

func TestPrint(t *testing.T) {
	log.Println(Print(TestErr))
}

func TestPrintWithDepthToLogBuffer(t *testing.T) {
	buf := PrintWithDepthToLogBuffer(TestErr, 0)
	_, _ = log.Writer().Write(buf.Bytes())
}

func TestClose(t *testing.T) {
	fp, err := os.Open("LICENSE")
	Panic(err)
	defer Close(fp)
	info, err := fp.Stat()
	Panic(err)
	modTime := info.ModTime().UTC()
	t.Log(modTime)
	//t.Log(modTime.GoString())
	if !modTime.Equal(time.Date(2022, time.March, 30, 0, 27, 25, 368677500, time.UTC)) {
		t.Fail()
	}
}

func TestCloseResponse(t *testing.T) {
	resp, err := http.Head("https://www.baidu.com/")
	Panic(err)
	defer CloseResponse(resp)
	t.Log("ContentLength", resp.ContentLength)
}
