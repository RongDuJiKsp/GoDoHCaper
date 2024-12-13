package main

import (
	"go-godoh-damon/child"
	"go-godoh-damon/godoh"
	"go-godoh-damon/logger"
	"os"
	"time"
)

func main() {
	logger.Log("守护程序 开始执行....")
	for {
		p, err := child.CreateChildProcess(`godoh`, append([]string{"agent"}, os.Args[1:]...)...)
		if err != nil {
			logger.Log("What happened? ", err, "WaitExit")
			time.Sleep(3 * time.Second)
			continue
		}
		logger.Log("子进程创建成功")
		p.Run(func(stream *child.IOStream) {
			logger.Log("系统启动")
			i := godoh.NewExitingReader(stream)
			go i.SyncWaitKill(p.Cmd())
			godoh.SyncListen(stream, []godoh.LineReader{i})
		})
		p.WaitExit()
		time.Sleep(1 * time.Second)
		logger.Log("进程退出，重启中...")
	}
}
