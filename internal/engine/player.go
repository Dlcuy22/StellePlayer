package engine

import (
	"fmt"
	"os/exec"
	"sync"
)

type Player struct {
	cmd   *exec.Cmd
	cmdMu sync.Mutex
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) Play(path string) error {
	p.cmdMu.Lock()
	defer p.cmdMu.Unlock()

	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Kill()
		p.cmd.Process.Wait()
	}

	p.cmd = exec.Command("ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", path)
	return p.cmd.Start()
}

func (p *Player) Stop() {
	p.cmdMu.Lock()
	defer p.cmdMu.Unlock()

	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Kill()
		p.cmd.Process.Wait()
		p.cmd = nil
	}
}

func (p *Player) Pause() {
	p.Stop() // ffplay doesn't support pause via command line easily without complex IPC, so we stop for now.
	// In the original code, pause was implemented by killing the process.
	// We will keep the same behavior for now, but the caller needs to handle the state.
}

func (p *Player) Resume(path string, seekTime float64) error {
	p.cmdMu.Lock()
	defer p.cmdMu.Unlock()

	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Kill()
		p.cmd.Process.Wait()
	}

	p.cmd = exec.Command("ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet",
		"-ss", fmt.Sprintf("%.2f", seekTime), path)
	return p.cmd.Start()
}

func (p *Player) Seek(path string, seekTime float64) error {
	return p.Resume(path, seekTime)
}

func (p *Player) Wait() error {
	// We need to be careful here not to hold the lock while waiting,
	// but we need to know which command we are waiting for.
	p.cmdMu.Lock()
	cmd := p.cmd
	p.cmdMu.Unlock()

	if cmd == nil {
		return nil
	}

	return cmd.Wait()
}

func (p *Player) IsRunning() bool {
	p.cmdMu.Lock()
	defer p.cmdMu.Unlock()
	return p.cmd != nil && p.cmd.Process != nil && p.cmd.ProcessState == nil
}
