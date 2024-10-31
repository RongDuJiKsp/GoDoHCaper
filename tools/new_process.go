package tools

import (
	"fmt"
	"os/exec"
)

func NewProcessToRun() *exec.Cmd {
	cmd := exec.Command(`godoh`, "agent", "-d", "tunnel.safecv.cn", "-p", "cloudflare")
	fmt.Println("创建子进程成功")
	return cmd
}
