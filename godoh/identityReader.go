package godoh

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/emirpasic/gods/v2/queues/arrayqueue"
	"go-godoh-proxy/child"
	"go-godoh-proxy/logger"
	"strings"
	"sync"
	"time"
)

const (
	MaxClient  = 5
	FirstStart = 10 * time.Second
)

type IdentityReader struct {
	stream           *child.IOStream
	registerIdentity string
	connIdentities   *arrayqueue.Queue[string]
	mutex            sync.Mutex
	turnNext         chan string
}

func NewIdentityReader(stream *child.IOStream) *IdentityReader {
	return &IdentityReader{stream, "", arrayqueue.New[string](), sync.Mutex{}, make(chan string, 500)}
}

func (i *IdentityReader) Run(cmd string) {
	writer := bufio.NewWriter(i.stream.In())
	_, _ = writer.WriteString(cmd + "\n")
	_ = writer.Flush()
}
func (i *IdentityReader) Use(identity string) {
	if identity == "" {
		return
	}
	i.Run("use " + identity)
	i.registerIdentity = identity
}
func (i *IdentityReader) NewClient(identity string) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.connIdentities.Enqueue(identity)
	for i.connIdentities.Size() > MaxClient {
		i.connIdentities.Dequeue()
	}

}
func (i *IdentityReader) NextClient(living bool, livingClient string) (string, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if i.connIdentities.Empty() {
		return "", errors.New("no other connected identities")
	}
	if living {
		i.connIdentities.Enqueue(livingClient)
	}
	c, _ := i.connIdentities.Dequeue()
	return c, nil
}

func (i *IdentityReader) SyncHandleOnBallingOrTimeout(timeout time.Duration, fn func(identity string), running *bool) {
	timer := time.NewTimer(FirstStart) //首次启动等待10秒
	for *running {
		var nextClient string
		select {
		case livingClient := <-i.turnNext:
			timer.Stop()
			logger.Log("Transfer ", livingClient, " ok")
			n, err := i.NextClient(true, livingClient)
			if err != nil {
				nextClient = livingClient
			} else {
				nextClient = n
				i.Use(nextClient)
			}
		case <-timer.C:
			logger.Log("Transfer ", nextClient, " timeout")
			n, err := i.NextClient(false, "")
			if err == nil {
				nextClient = n
				i.Use(nextClient)
			}
		}
		logger.Log("Tick Transfer ", nextClient)
		if nextClient != "" {
			fn(nextClient)
		}
		timer.Reset(timeout)
	}
}
func (i *IdentityReader) NextLine(line []byte) {
	strLine := string(line)
	id, err := hasNewClient(strLine)
	if err != nil {
		i.NewClient(id)
		logger.Log("New client ", id)
	}
	if hasFinished(strLine) {
		i.turnNext <- i.registerIdentity
	}
}

func (i *IdentityReader) Close() {

}
func hasFinished(strLine string) bool {
	if strings.Contains(strLine, "Writing file to desk") {
		return true
	}
	return false
}
func hasNewClient(strLine string) (string, error) {
	if strings.Contains(strLine, "First time checkin for agent") {
		id, err := getIdByRegisterLine(strLine)
		if err == nil {
			return id, nil
		} else {
			fmt.Println(err)
			return "", err
		}
	}
	return "", errors.New("no agent found")
}
func getIdByRegisterLine(line string) (string, error) {
	sp := strings.Split(line, "ident=")
	if len(sp) < 2 || len(sp[1]) < 5 {
		return "", errors.New("异常：identity不对")
	}
	return sp[1][:5], nil
}
