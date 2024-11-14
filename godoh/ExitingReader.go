package godoh

import (
	"go-godoh-damon/child"
	"go-godoh-damon/logger"
	"os/exec"
	"strings"
)

type ExitingReader struct {
	stream *child.IOStream
	exit   chan struct{}
}

func (e *ExitingReader) NextLine(line []byte) {
	str := string(line)
	if hasExiting(str) {
		e.exit <- struct{}{}
	}
}

func (e *ExitingReader) Close() {

}
func (e *ExitingReader) SyncWaitKill(press *exec.Cmd) {
	<-e.exit
	if err := press.Process.Kill(); err != nil {
		logger.Log("Kill with err:", err.Error())
	}
}
func NewExitingReader(stream *child.IOStream) *ExitingReader {
	return &ExitingReader{stream: stream, exit: make(chan struct{})}
}
func hasExiting(s string) bool {
	return strings.Contains(s, "Exiting")
}
