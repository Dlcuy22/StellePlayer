# internal/engine/player.go

## Overview

Provides a wrapper around the `ffplay` command to control audio playback.

## Types

### `Player`

Struct managing the playback process.

- `cmd`: Pointer to the `exec.Cmd` of the running `ffplay` process.
- `cmdMu`: Mutex to protect access to `cmd`.

## Functions

### `NewPlayer`

```go
func NewPlayer() *Player
```

Creates a new `Player` instance.

### `Play`

```go
func (p *Player) Play(path string) error
```

Starts playing the file at `path`.

- Kills any existing process.
- Starts `ffplay` with flags: `-nodisp`, `-autoexit`, `-loglevel quiet`.

### `Stop`

```go
func (p *Player) Stop()
```

Stops playback by killing the `ffplay` process.

### `Pause`

```go
func (p *Player) Pause()
```

Pauses playback. Currently implemented as `Stop()` because `ffplay` does not support external pause control easily.

### `Resume`

```go
func (p *Player) Resume(path string, seekTime float64) error
```

Resumes playback from `seekTime`.

- Restarts `ffplay` with the `-ss` flag to seek to the specified time.

### `Seek`

```go
func (p *Player) Seek(path string, seekTime float64) error
```

Alias for `Resume`. Used to jump to a specific time.

### `Wait`

```go
func (p *Player) Wait() error
```

Waits for the current process to exit. Used to detect when a song finishes naturally.

### `IsRunning`

```go
func (p *Player) IsRunning() bool
```

Returns `true` if the player process is currently active.
