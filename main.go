package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("启动守护进程")
	for {
		fmt.Println("启动子进程")
		cmd := exec.Command(`godoh`, "agent", "-d", "tunnel.safecv.cn", "-p", "cloudflare")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
