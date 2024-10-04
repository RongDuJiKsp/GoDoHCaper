package main

import (
	"fmt"
	"go-godoh-damon/tools"
	"os/exec"
	"time"
)

func run(cmd *exec.Cmd) {
	fmt.Println("开始执行子进程")
	_ = cmd.Run()
}

func main() {
	fmt.Println("开始执行....")
	nowProcess := tools.NewProcessToRun()
	go run(nowProcess)
	for p := range makeCommandChan() {
		_ = nowProcess.Process.Kill()
		fmt.Println("正在重启子进程")
		nowProcess = p
		go run(p)
	}

}
func makeCommandChan() chan *exec.Cmd {
	ch := make(chan *exec.Cmd)
	go func() {
		for range time.Tick(time.Minute * 8) {
			ch <- tools.NewProcessToRun()
		}
	}()
	return ch
}
