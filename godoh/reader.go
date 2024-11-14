package godoh

import (
	"bufio"
	"go-godoh-damon/child"
	"go-godoh-damon/logger"
)

type LineReader interface {
	NextLine(line []byte)
	Close()
}

func SyncListen(stream *child.IOStream, listeners []LineReader) {
	scanner := bufio.NewScanner(stream.Out())
	for scanner.Scan() {
		line := scanner.Text()
		logger.Output(line)
		for _, l := range listeners {
			l.NextLine([]byte(line))
		}
	}
	for _, l := range listeners {
		l.Close()
	}
	logger.Log("std::out closed")
}
