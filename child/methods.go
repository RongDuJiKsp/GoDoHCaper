package child

func (p *Process) closeAll() {
	_ = p.stream.in.Close()
	_ = p.stream.out.Close()
}
func (p *Process) Wait() {
	_ = p.cmd.Wait()
	p.closeAll()
}
func (p *Process) Init(fn func(stream *IOStream)) {
	fn(p.stream)
}
