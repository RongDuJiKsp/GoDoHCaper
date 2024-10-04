package tools

import "os/exec"

func NewProcessToRun() *exec.Cmd {
	cmd := exec.Command(`godoh`, "agent", "-d", "tunnel.safecv.cn", "-p", "cloudflare")
	return cmd
}
