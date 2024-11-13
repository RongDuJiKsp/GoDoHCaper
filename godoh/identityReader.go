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
	MaxClient      = 5
	FirstStart     = 10 * time.Second
	FreeWait       = 5 * time.Second
	BallingTimeout = 5 * time.Second
)

type IdentityReader struct {
	stream           *child.IOStream
	registerIdentity string
	connIdentities   *arrayqueue.Queue[string]
	queueLock        *sync.Mutex
	cmd              *sync.Mutex
	turnNext         chan string
	lastBalling      time.Time
}

func NewIdentityReader(stream *child.IOStream) *IdentityReader {
	return &IdentityReader{stream, "", arrayqueue.New[string](), &sync.Mutex{}, &sync.Mutex{}, make(chan string, 500), time.Now()}
}

func (i *IdentityReader) Run(cmd string) {
	i.cmd.Lock()
	defer i.cmd.Unlock()
	logger.Log("Run Cmd: ", cmd)
	writer := bufio.NewWriter(i.stream.In())
	var err error
	_, err = writer.WriteString(cmd + "\n\n")
	if err != nil {
		logger.Log("Error writing cmd: ", err)
	}
	err = writer.Flush()
	if err != nil {
		logger.Log("Error writing cmd: ", err)
	}
}
func (i *IdentityReader) Use(identity string) {
	if identity == "" {
		return
	}
	i.Run("use " + identity)
	i.registerIdentity = identity
}
func (i *IdentityReader) NewClient(identity string) {
	i.queueLock.Lock()
	defer i.queueLock.Unlock()
	i.connIdentities.Enqueue(identity)
	for i.connIdentities.Size() > MaxClient {
		i.connIdentities.Dequeue()
	}

}
func (i *IdentityReader) NextClient(living bool, livingClient string) (string, error) {
	i.queueLock.Lock()
	defer i.queueLock.Unlock()
	if i.connIdentities.Empty() {
		if i.registerIdentity != "" {
			i.connIdentities.Enqueue(i.registerIdentity)
		}
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
	var nextClient string
	for *running {
		select {
		case livingClient := <-i.turnNext:
			timer.Stop()
			logger.Log("Transfer ", livingClient, " ok")
			n, err := i.NextClient(true, livingClient)
			if err != nil {
				nextClient = livingClient
			} else {
				nextClient = n
			}
		case <-timer.C:
			logger.Log("Transfer ", nextClient, " timeout")
			if time.Now().Sub(i.lastBalling) < BallingTimeout {
				timer.Reset(timeout)
				continue
			}
			n, err := i.NextClient(false, "")
			if err == nil {
				nextClient = n
			}
		}
		logger.Log("Tick Transfer ", nextClient)
		if nextClient != "" {
			fn(nextClient)
			timer.Reset(timeout)
		} else {
			timer.Reset(FreeWait)
		}

	}
}
func (i *IdentityReader) NextLine(line []byte) {
	strLine := string(line)
	id, err := hasNewClient(strLine)
	if err == nil {
		i.NewClient(id)
		logger.Log("New client ", id)
	}
	if hasFinished(strLine) {
		i.turnNext <- i.registerIdentity
	}
	if hasBalling(strLine) {
		i.lastBalling = time.Now()
	}
}

func (i *IdentityReader) Close() {

}
func hasFinished(strLine string) bool {
	if strings.Contains(strLine, "Writing file to disk") {
		return true
	}
	return false
}
func hasBalling(strLine string) bool {
	if strings.Contains(strLine, "Question had less than 9 labels, bailing") {
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
