package godoh

import (
	"bufio"
	"errors"
	"go-godoh-proxy/child"
	"strings"
	"time"
)

type IdentityReader struct {
	stream           *child.IOStream
	registerIdentity string
}

func NewIdentityReader(stream *child.IOStream) *IdentityReader {
	return &IdentityReader{stream, ""}
}

func (i *IdentityReader) Run(cmd string) {
	writer := bufio.NewWriter(i.stream.In())
	_, _ = writer.WriteString(cmd + "\n")
	_ = writer.Flush()
}
func (i *IdentityReader) RequestIdentity() {
	i.Run("agents")
}
func (i *IdentityReader) Use(identity string) {
	i.Run("use " + identity)
}
func (i *IdentityReader) SyncTickHandle(duration time.Duration, fn func(identity string), running *bool) {
	for range time.Tick(duration) {
		if !*running {
			break
		}
		fn(i.registerIdentity)
	}
}
func (i *IdentityReader) NextLine(line []byte) {
	strLine := string(line)
	if strings.Contains(strLine, "First time checkin for agent") {
		id, err := getIdByRegisterLine(strLine)
		if err == nil {
			i.registerIdentity = id
			i.Use(id)
		}
	}
}

func (i *IdentityReader) Close() {

}
func getIdByRegisterLine(line string) (string, error) {
	sp := strings.Split(line, "ident=")
	if len(sp) < 2 || len(sp[1]) < 5 {
		return "", errors.New("异常：identity不对")
	}
	return sp[1][:5], nil
}
