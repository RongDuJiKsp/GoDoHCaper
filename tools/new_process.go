package tools

import (
	"fmt"
	"os/exec"
)

func NewProcessToRun() *exec.Cmd {
	cmd := exec.Command(`godoh`, "agent", "-d", "send.tunvision.work", "-p", "cloudflare", "-t", "15")
	fmt.Println("创建子进程成功")
	return cmd
}
