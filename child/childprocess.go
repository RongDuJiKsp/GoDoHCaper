package child

import (
	"io"
	"os/exec"
)

type IOStream struct {
	in  io.WriteCloser
	out io.ReadCloser
}

func (I *IOStream) In() io.WriteCloser {
	return I.in
}

func (I *IOStream) Out() io.ReadCloser {
	return I.out
}

type Process struct {
	stream *IOStream
	cmd    *exec.Cmd
}

func CreateChildProcess(commandName string, arg ...string) (*Process, error) {
	cmd := exec.Command(commandName, arg...)
	childIn, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	childOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &Process{&IOStream{childIn, childOut}, cmd}, nil
}
