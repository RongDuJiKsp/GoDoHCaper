package main

import (
	"go-godoh-damon/tools"
	"os/exec"
	"time"
)

func run(cmd *exec.Cmd) {
	_ = cmd.Run()
}

func main() {
	nowProcess := tools.NewProcessToRun()
	go run(nowProcess)
	for p := range makeCommandChan() {
		_ = nowProcess.Process.Kill()
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
