package errs

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"reflect"
	"runtime/debug"
)

var Debug = true

func Panic(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func PrintWithDepthToLog(err error, depth int) bool {
	if err == nil {
		return false
	}
	return printWithDepth(err, depth+1, log.Writer(), log.Output)
}

func PrintWithDepthToLogger(err error, depth int, log *log.Logger) bool {
	if err == nil {
		return false
	}
	return printWithDepth(err, depth+1, log.Writer(), log.Output)
}

func PrintWithDepthToLogBuffer(err error, depth int) *bytes.Buffer {
	if err == nil {
		return nil
	}
	buf := &bytes.Buffer{}
	PrintWithDepthToLogger(err, depth+1, log.New(buf, log.Prefix(), log.Flags()))
	return buf
}

func printWithDepth(err error, depth int, writer io.Writer, output func(int, string) error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	depth += 2
	e := output(depth, s)
	if e != nil {
		log.Println(e)
		log.Println(s)
	}
	if //goland:noinspection GoBoolExpressions
	!Debug {
		return true
	}
	stack := debug.Stack()
	stackLine := stack
	i := bytes.IndexByte(stackLine, '\n') + 1
	stackLine = stack[:i]
	_, e = writer.Write(stackLine)
	if e != nil {
		log.Println(e)
		log.Println(string(stackLine))
	}
	stackLine = stack[i:]
	depth *= 2
	for line := 0; line < depth; line++ {
		i = bytes.IndexByte(stackLine, '\n') + 1
		if i == 0 {
			break
		}
		stackLine = stackLine[i:]
	}
	_, e = writer.Write(stackLine)
	if e != nil {
		log.Println(e)
		log.Println(string(stackLine))
	}
	return true
}

func Print(err error) bool {
	return PrintWithDepthToLog(err, 1)
}

func printErr(err error) bool {
	return PrintWithDepthToLog(err, 2)
}

func Close(closer io.Closer) {
	printErr(closer.Close())
}

func CloseResponse(resp *http.Response) {
	printErr(resp.Body.Close())
}

func Defer(fn func() error) {
	printErr(fn())
}

func deferCall(fn func([]reflect.Value) []reflect.Value, args ...interface{}) {
	var values []reflect.Value
	for _, arg := range args {
		values = append(values, reflect.ValueOf(arg))
	}
	values = fn(values)
	length := len(values)
	if length == 0 {
		return
	}
	value := values[length-1].Interface()
	if value == nil {
		return
	}
	err, ok := value.(error)
	if !ok {
		return
	}
	PrintWithDepthToLog(err, 2)
}

func DeferCall(fn interface{}, args ...interface{}) {
	deferCall(reflect.ValueOf(fn).Call, args...)
}

func DeferCallSlice(fn interface{}, args ...interface{}) {
	deferCall(reflect.ValueOf(fn).CallSlice, args...)
}

func IsClosed(err error) bool {
	return errors.Is(err, fs.ErrClosed) || errors.Is(err, net.ErrClosed)
}

func CloseIgnoreClosed(closer io.Closer) {
	err := closer.Close()
	if IsClosed(err) {
		return
	}
	printErr(err)
}

func DeferIgnoreClosed(fn func() error) {
	err := fn()
	if IsClosed(err) {
		return
	}
	printErr(err)
}

func SetLogFlag() {
	log.SetFlags(log.Flags() | log.Llongfile)
}
