package godoh

import (
	"errors"
	"fmt"
	"go-godoh-proxy/child"
	"io"
	"strings"
	"time"
)

type IdentityReader struct {
	stream           *child.IOStream
	registerIdentity []string
}

func NewIdentityReader(stream *child.IOStream) *IdentityReader {
	return &IdentityReader{stream, nil}
}
func (i *IdentityReader) RequestIdentity() {
	_, _ = io.WriteString(i.stream.In(), "agents\n")
}
func (i *IdentityReader) Use(identity string) {
	_, _ = io.WriteString(i.stream.In(), fmt.Sprintf("use %s \n", identity))
}
func (i *IdentityReader) Run(cmd string) {
	_, _ = io.WriteString(i.stream.In(), cmd+"\n")
}
func (i *IdentityReader) SyncTickHandle(duration time.Duration, fn func(identity []string), running *bool) {
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
			i.registerIdentity = append(i.registerIdentity, id)
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
