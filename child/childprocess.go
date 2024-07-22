package child

import (
	"io"
	"os"
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
	err    io.ReadCloser
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
	childErr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	go func() {
		_, _ = io.Copy(os.Stderr, childErr)
	}()
	return &Process{&IOStream{childIn, childOut}, childErr, cmd}, nil
}
