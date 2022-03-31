package errs

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"runtime/debug"
)

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

func Close(closer io.Closer) {
	Print(closer.Close())
}

func CloseResponse(resp *http.Response) {
	Print(resp.Body.Close())
}
