# Function Reference

This document provides a detailed reference for the functions in the codebase.

## `internal/engine`

### `metadata.go`

#### `ExtractMetadata(filePath string) (Metadata, error)`

Extracts audio metadata from the given file using `ffprobe`.

- **Parameters**:
  - `filePath`: Absolute path to the audio file.
- **Returns**:
  - `Metadata`: Struct containing Title, Artist, Album, Duration, etc.
  - `error`: Error if extraction fails.

### `player.go`

#### `NewPlayer() *Player`

Creates a new instance of `Player`.

#### `(p *Player) Play(path string) error`

Starts playback of the specified file using `ffplay`.

- **Parameters**:
  - `path`: Path to the audio file.

#### `(p *Player) Stop()`

Stops the current playback by killing the process.

#### `(p *Player) Pause()`

Pauses playback (currently implements Stop).

#### `(p *Player) Resume(path string, seekTime float64) error`

Resumes playback from a specific time.

- **Parameters**:
  - `path`: Path to the audio file.
  - `seekTime`: Time in seconds to start playing from.

#### `(p *Player) Seek(path string, seekTime float64) error`

Seeks to a specific time (alias for Resume).

#### `(p *Player) Wait() error`

Waits for the playback process to complete. Useful for detecting when a song finishes.

#### `(p *Player) IsRunning() bool`

Checks if the player process is currently running.

## `internal/media`

### `media.go`

#### `LoadFromDirectory(dir string) ([]engine.Metadata, error)`

Scans a directory for supported audio files and extracts their metadata.

- **Parameters**:
  - `dir`: Directory path to scan.
- **Returns**:
  - `[]engine.Metadata`: List of metadata for found songs.

## `internal/lyrics`

### `lyrics.go`

#### `Parse(content string) []Line`

Parses LRC format string into a slice of `Line` structs.

#### `LoadFromFile(songPath, musicDir string) (Lyrics, bool)`

Attempts to load lyrics from a local `.lrc` file matching the song name.

#### `FetchFromAPI(artist, title, album string) (string, error)`

Fetches lyrics from the `lrclib.net` API.

#### `SaveToFile(songPath, musicDir, content string) (string, error)`

Saves fetched lyrics to a `.lrc` file in a `lyrics` subdirectory.

## `internal/app`

### `app.go`

#### `Run(musicDir string) error`

Initializes and runs the Bubble Tea program.

### `model.go`

#### `initialModel(musicDir string) *model`

Initializes the application state.

#### `(m *model) Init() tea.Cmd`

Bubble Tea Init function. Starts the loading tick.

#### `(m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd)`

Bubble Tea Update function. Handles all state changes and events.

#### `(m *model) View() string`

Bubble Tea View function. Renders the UI.

#### `(m *model) playSong(song *Song)`

Internal helper to start playback of a song and handle state updates.

#### `(m *model) stopPlayback()`

Internal helper to stop playback and reset state.

#### `(m *model) seek(seconds float64)`

Internal helper to handle seeking logic (calculating new time and calling player).

#### `loadLyricsAsync(song *Song, musicDir string) tea.Cmd`

Bubble Tea command to load lyrics in the background.

#### `getCurrentLyrics(ly *lyrics.Lyrics, currentTime float64) (current, next string)`

Helper to find the current and next lyric lines based on the timestamp.
