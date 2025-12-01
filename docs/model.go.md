# internal/app/model.go

## Overview

This file defines the core data structures and logic for the Bubble Tea application model. It handles the UI state, user input, playback control, and rendering.

## Types

### `Song`

Represents a music track.

- `metadata` (`engine.Metadata`): Metadata of the song.
- `lyrics` (`*lyrics.Lyrics`): Loaded lyrics for the song.

### `playerState`

Enum for playback state: `stateStopped`, `statePlaying`, `statePaused`.

### `model`

The main Bubble Tea model struct.

- `songs`: List of loaded songs.
- `list`: Playlist UI component.
- `progress`: Progress bar UI component.
- `state`: Current player state.
- `currentSong`: Pointer to the currently playing song.
- `currentTime`: Current playback position in seconds.
- `player`: Instance of `engine.Player`.
- `shuffle`: Boolean flag for shuffle mode.
- `lyricsLoading`: Boolean flag for lyrics loading status.
- `loading`: Boolean flag for initial loading state.

## Functions

### `initialModel`

```go
func initialModel(musicDir string) *model
```

Initializes the model with default values and UI components.

### `Init`

```go
func (m *model) Init() tea.Cmd
```

Returns the initial command to run (starts loading tick).

### `Update`

```go
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd)
```

The main update loop. Handles:

- `WindowSizeMsg`: Resizes UI components.
- `songsLoadedMsg`: Populates the song list.
- `loadingTickMsg`: Animates the loading screen.
- `lyricsLoadedMsg`: Updates lyrics for the current song.
- `tea.KeyMsg`: Handles keyboard input (Play, Pause, Stop, Next, Prev, Seek, Shuffle, Quit).
- `tickMsg`: Updates playback progress.

### `View`

```go
func (m *model) View() string
```

Renders the UI string.

- Shows loading screen if `m.loading` is true.
- Splits screen into Left Panel (Info, Controls, Lyrics) and Right Panel (Playlist).

### `playSong`

```go
func (m *model) playSong(song *Song)
```

Internal helper to start playback.

- Stops current playback.
- Resets state variables.
- Calls `m.player.Play()`.
- Starts a goroutine to wait for playback completion (`m.player.Wait()`).

### `stopPlayback`

```go
func (m *model) stopPlayback()
```

Stops the player and resets state.

### `pausePlayback` / `resumePlayback`

Pauses and resumes playback using `m.player`.

### `seek`

```go
func (m *model) seek(seconds float64)
```

Seeks forward or backward by the specified seconds.

### `loadLyricsAsync`

```go
func loadLyricsAsync(song *Song, musicDir string) tea.Cmd
```

Returns a `tea.Cmd` that loads lyrics in the background (from file or API).

### `getCurrentLyrics`

```go
func getCurrentLyrics(ly *lyrics.Lyrics, currentTime float64) (current, next string)
```

Helper to determine which lyric lines to display based on the current time.
