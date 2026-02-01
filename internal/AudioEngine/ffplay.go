// AudioEngine/ffplay.go
// FFplay-based audio engine implementation.
//
// Types:
//   - FFplayEngine: implements Engine interface using ffplay subprocess
//
// Functions:
//   - NewFFplayEngine: creates a new FFplay engine instance
//   - Play, Stop, Pause, Resume, Seek: playback control methods
//   - GetState, SetOnComplete: state management methods

package AudioEngine

import (
	"fmt"
	"os/exec"
	"strconv"
	"sync"
)

// FFplayEngine implements the Engine interface using ffplay.
type FFplayEngine struct {
	cmd        *exec.Cmd
	mu         sync.Mutex
	state      PlaybackState
	onComplete func()
	filePath   string
}

// NewFFplayEngine creates a new FFplay-based audio engine.
func NewFFplayEngine() *FFplayEngine {
	return &FFplayEngine{
		state: StateStopped,
	}
}

// Play starts playing the audio file.
func (e *FFplayEngine) Play(filePath string, seekTo float64, volume int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.stopInternal()

	e.filePath = filePath
	e.state = StatePlaying

	args := []string{"-nodisp", "-autoexit", "-loglevel", "quiet", "-volume", strconv.Itoa(volume)}
	if seekTo > 0 {
		args = append(args, "-ss", fmt.Sprintf("%.2f", seekTo))
	}
	args = append(args, filePath)

	e.cmd = exec.Command("ffplay", args...)
	cmd := e.cmd

	if err := cmd.Start(); err != nil {
		e.state = StateStopped
		return err
	}

	go func() {
		cmd.Wait()
		e.mu.Lock()
		if e.cmd == cmd && e.state == StatePlaying {
			e.state = StateStopped
			if e.onComplete != nil {
				callback := e.onComplete
				e.mu.Unlock()
				callback()
				return
			}
		}
		e.mu.Unlock()
	}()

	return nil
}

// Stop stops playback and resets state.
func (e *FFplayEngine) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.stopInternal()
	e.state = StateStopped
}

// stopInternal kills the current process without locking.
func (e *FFplayEngine) stopInternal() {
	if e.cmd != nil && e.cmd.Process != nil {
		e.cmd.Process.Kill()
		e.cmd.Process.Wait()
		e.cmd = nil
	}
}

// Pause pauses playback.
func (e *FFplayEngine) Pause() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.stopInternal()
	e.state = StatePaused
}

// Resume resumes playback from the given position.
func (e *FFplayEngine) Resume(seekTo float64, volume int) error {
	e.mu.Lock()
	if e.state != StatePaused || e.filePath == "" {
		e.mu.Unlock()
		return nil
	}
	e.mu.Unlock()

	return e.Play(e.filePath, seekTo, volume)
}

// Seek jumps to the specified position.
func (e *FFplayEngine) Seek(position float64, volume int) error {
	e.mu.Lock()
	filePath := e.filePath
	wasPlaying := e.state == StatePlaying
	e.mu.Unlock()

	if filePath == "" {
		return nil
	}

	if wasPlaying {
		return e.Play(filePath, position, volume)
	}
	return nil
}

// GetState returns the current playback state.
func (e *FFplayEngine) GetState() PlaybackState {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.state
}

// SetOnComplete sets the callback for when playback finishes.
func (e *FFplayEngine) SetOnComplete(callback func()) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.onComplete = callback
}

// SetFilePath sets the current file path (used when loading a new song).
func (e *FFplayEngine) SetFilePath(filePath string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.filePath = filePath
}
