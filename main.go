package main

import (
	"fmt"
	"go-godoh-damon/tools"
	"log"
	"os"
	"os/exec"
	"time"
)

func run(cmd *exec.Cmd) {
	fmt.Println("开始执行子进程")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("开始执行....")
	nowProcess := tools.NewProcessToRun()
	go run(nowProcess)
	for p := range makeCommandChan() {
		if err := nowProcess.Process.Kill(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("正在重启子进程")
		nowProcess = p
		go run(p)
	}

}
func makeCommandChan() chan *exec.Cmd {
	ch := make(chan *exec.Cmd)
	go func() {
		for range time.Tick(time.Minute * 2) {
			ch <- tools.NewProcessToRun()
		}
	}()
	return ch
}
