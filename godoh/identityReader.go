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
	FirstStart     = 10 * time.Second
	FreeWait       = 5 * time.Second
	BallingTimeout = 5 * time.Second
)

type IdentityReader struct {
	stream         *child.IOStream
	connIdentities *arrayqueue.Queue[string]
	queueLock      *sync.Mutex
	cmd            *sync.Mutex
	turnNext       chan struct{}
	lastBalling    time.Time
	ballingLock    *sync.Mutex
}

func NewIdentityReader(stream *child.IOStream) *IdentityReader {
	return &IdentityReader{stream, arrayqueue.New[string](), &sync.Mutex{}, &sync.Mutex{}, make(chan struct{}, 500), time.Now(), &sync.Mutex{}}
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
}
func (i *IdentityReader) Balling(time time.Time) {
	i.ballingLock.Lock()
	defer i.ballingLock.Unlock()
	i.lastBalling = time
}
func (i *IdentityReader) IsTimeout() bool {
	i.ballingLock.Lock()
	defer i.ballingLock.Unlock()
	return time.Now().Sub(i.lastBalling) >= BallingTimeout
}
func (i *IdentityReader) NewClient(identity string) {
	i.queueLock.Lock()
	defer i.queueLock.Unlock()
	i.connIdentities.Enqueue(identity)
}
func (i *IdentityReader) NextClient(lastClient string) (string, error) {
	i.queueLock.Lock()
	defer i.queueLock.Unlock()
	if lastClient != "" {
		i.connIdentities.Enqueue(lastClient)
	}
	c, ok := i.connIdentities.Dequeue()
	if !ok {
		return "", errors.New("no client found")
	}
	return c, nil
}

func (i *IdentityReader) SyncHandleOnBallingOrTimeout(timeout time.Duration, fn func(identity string), running *bool) {
	timer := time.NewTimer(FirstStart) //首次启动等待10秒
	var ballingClient string
	for *running {
		select {
		case <-i.turnNext:
			timer.Stop()
			logger.Log("Transfer ", ballingClient, " ok")
			n, err := i.NextClient(ballingClient)
			if err == nil {
				ballingClient = n
			}
		case <-timer.C:
			if i.IsTimeout() {
				timer.Reset(timeout)
				continue
			}
			logger.Log("Transfer ", ballingClient, " timeout")
			n, err := i.NextClient("") // timeout bye~
			if err == nil {
				ballingClient = n
			}
		}
		logger.Log("Tick Transfer ", ballingClient)
		if ballingClient != "" {
			fn(ballingClient)
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
		i.Balling(time.Now())
		i.turnNext <- struct{}{}
	}
	if hasBalling(strLine) {
		i.Balling(time.Now())
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
