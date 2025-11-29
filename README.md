# Stelle Music Player TUI

A terminal-based music player with lyrics support built with Go and Bubble Tea.

## Features

- Play music files (MP3, M4A, FLAC, WAV, OGG, AAC)
- Display metadata (title, artist, album, duration, bitrate)
- Synchronized lyrics display with LRC file support
- Automatic lyrics fetching from LRCLIB API
- Basic playback controls (play, pause, stop, next, previous)
- Seek functionality (forward/backward 5 seconds)
- Shuffle mode
- File filtering and search

## Installation

1. Clone the repository
2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Build the application:

   ```bash
   go build -o player.exe
   ```

note: make sure you have ffplay in your PATH

## Usage

Run with a music directory:

```bash
./player.exe -sd /path/to/music/directory
```

Or run without flags to select a folder interactively:

```bash
./player.exe
```

## Controls

- `p` - Play
- `Space` - Pause/Resume
- `s` - Stop
- `n` / `→` - Next song
- `b` / `←` - Previous song
- `f` - Forward 5 seconds
- `r` - Rewind 5 seconds
- `h` - Toggle shuffle
- `q` / `Ctrl+C` - Quit

## Lyrics

The player automatically searches for LRC files in a `lyrics` subdirectory within your music folder. If no local lyrics are found, it attempts to fetch them from the LRCLIB API and saves them for future use.

## Dependencies

- Go 1.19+
- FFmpeg (for metadata extraction)

## License

MIT License
