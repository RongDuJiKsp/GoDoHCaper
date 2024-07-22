package main

import (
	"go-godoh-proxy/child"
	"go-godoh-proxy/godoh"
	"time"
)

func main() {
	p, err := child.CreateChildProcess(`./godoh`)
	if err != nil {
		panic(err)
	}
	p.Handle(func(stream *child.IOStream) {
		isRunning := true
		sendDuration := 3000 * time.Second
		sendCommands := []string{"ls -a"}

		i := godoh.NewIdentityReader(stream)
		i.RequestIdentity()
		go i.SyncTickHandle(sendDuration, func(identity []string) {
			for _, id := range identity {
				i.Use(id)
				for _, cmd := range sendCommands {
					i.Run(cmd)
				}
			}
		}, &isRunning)
		go godoh.SyncListen(stream, []godoh.LineReader{i})
	})
	p.Wait()
}
