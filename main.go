package main

import (
	"fmt"
	"go-godoh-damon/tools"
	"log"
	"os"
	"os/exec"
)

func run(cmd *exec.Cmd) {
	fmt.Println("开始执行子进程")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}
}

func main() {
	fmt.Println("开始执行....")
	for {
		run(tools.NewProcessToRun())
	}
}
