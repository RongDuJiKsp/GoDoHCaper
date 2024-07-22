package godoh

import (
	"bufio"
	"go-godoh-proxy/child"
	"go-godoh-proxy/logger"
)

type LineReader interface {
	NextLine(line []byte)
	Close()
}

func SyncListen(stream *child.IOStream, listeners []LineReader) {
	scanner := bufio.NewScanner(stream.Out())
	for scanner.Scan() {
		line := scanner.Bytes()
		logger.Output(string(line))
		for _, l := range listeners {
			l.NextLine(line)
		}
	}
	for _, l := range listeners {
		l.Close()
	}
}
