package godoh

import (
	"bufio"
	"go-godoh-proxy/child"
)

type LineReader interface {
	NextLine(line []byte)
	Close()
}

func SyncListen(stream *child.IOStream, listeners []LineReader) {
	scanner := bufio.NewScanner(stream.Out())
	for scanner.Scan() {
		line := scanner.Bytes()
		for _, l := range listeners {
			l.NextLine(line)
		}
	}
	for _, l := range listeners {
		l.Close()
	}
}
