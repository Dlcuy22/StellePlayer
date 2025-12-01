# Music Player Documentation

This document explains the high-level workflow of the Music Player application, detailing how it loads, plays, pauses, seeks, and manages lyrics.

## 1. Initialization

The application entry point is `main.go`.

1.  **Flag Parsing**: It parses command-line flags, specifically `-sd` for the source directory.
2.  **Directory Selection**: If no directory is provided, it prompts the user to select one using a file dialog (`service.PickFolder`).
3.  **App Start**: It calls `app.Run(musicDir)` to start the Bubble Tea program.

## 2. Application Loop (Bubble Tea)

The core logic resides in `internal/app/model.go`, which implements the Bubble Tea `Model` interface.

### `Init`

- Starts a loading tick (`loadingTickCmd`) to animate the loading screen.
- Enters the alternate screen buffer.
- Spawns a goroutine (in `app.Run`) to load songs asynchronously.

### `Update`

Handles messages and updates the model state.

#### Loading Phase

- **`songsLoadedMsg`**: Received when songs are loaded. Populates `m.songs`, creates list items, and auto-plays the first song if available.
- **`loadingTickMsg`**: Updates the loading animation dots.

#### Playback Phase

- **`tea.KeyMsg`**: Handles user input (Play, Pause, Stop, Next, Previous, Seek, Shuffle, Quit).
- **`tickMsg`**: Updates the playback progress (`m.currentTime`) every 100ms. If the song finishes, it triggers `playNextCmd`.
- **`lyricsLoadedMsg`**: Updates the current song with loaded lyrics.

### `View`

Renders the UI based on the current state (`loading` or `playback`).

- **Left Panel**: Displays current song info, status, progress bar, time, controls, and lyrics.
- **Right Panel**: Displays the playlist (using `bubbles/list`).

## 3. Audio Engine (`internal/engine`)

The audio processing logic is encapsulated in the `internal/engine` package.

### Metadata Extraction

`engine.ExtractMetadata(filePath)` uses `ffmpeg-go` (wrapping `ffprobe`) to extract:

- Title, Artist, Album
- Duration
- Bitrate, Codec, Sample Rate

### Playback Control (`engine.Player`)

The `Player` struct manages the `ffplay` process.

- **`Play(path)`**: Starts `ffplay` in a new process.
- **`Stop()`**: Kills the current `ffplay` process.
- **`Pause()`**: Currently implemented as `Stop()` (since `ffplay` doesn't support IPC pause easily).
- **`Resume(path, seekTime)`**: Starts `ffplay` at the specified `seekTime`.
- **`Seek(path, seekTime)`**: Same as `Resume`.
- **`Wait()`**: Waits for the `ffplay` process to finish.

## 4. Workflows

### Loading Songs

1.  `app.Run` calls `media.LoadFromDirectory`.
2.  `media.LoadFromDirectory` walks the directory, filtering for supported extensions.
3.  For each file, it calls `engine.ExtractMetadata`.
4.  The list of `engine.Metadata` is converted to `app.Song` structs.
5.  A `songsLoadedMsg` is sent to the Bubble Tea loop.

### Playing a Song

1.  User selects a song and presses `p` (or auto-play triggers).
2.  `playSongCmd` is called.
3.  `m.player.Play(path)` starts the audio.
4.  A goroutine is spawned to wait for the player to finish (`m.player.Wait()`).
5.  Lyrics are loaded asynchronously (`loadLyricsAsync`).

### Pausing/Resuming

- **Pause**: Calls `m.player.Pause()`, sets state to `statePaused`.
- **Resume**: Calls `m.player.Resume(path, currentTime)`, sets state to `statePlaying`.

### Seeking

1.  User presses `f` (forward) or `r` (backward).
2.  `m.seek(seconds)` calculates the new time.
3.  `m.player.Seek(path, newTime)` restarts playback at the new position.

### Lyrics

1.  `loadLyricsAsync` checks for local `.lrc` files.
2.  If not found, it fetches from `lrclib.net` API (`lyrics.FetchFromAPI`).
3.  Lyrics are saved locally (`lyrics.SaveToFile`) and parsed (`lyrics.Parse`).
4.  `lyricsLoadedMsg` updates the UI.
5.  `getCurrentLyrics` syncs lyrics with `m.currentTime` during rendering.
