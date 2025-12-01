# internal/lyrics/lyrics.go

## Overview

Handles parsing, loading, fetching, and saving of synchronized lyrics (.lrc files).

## Types

### `Line`

Represents a single line of lyrics.

- `Time` (float64): Timestamp in seconds.
- `Text` (string): Lyric text.

### `Lyrics`

Container for lyrics.

- `Lines` ([]Line): Sorted list of lyric lines.
- `Loaded` (bool): Whether lyrics were successfully loaded.

## Functions

### `Parse`

```go
func Parse(content string) []Line
```

Parses a raw LRC string into a slice of `Line` structs.

- Supports standard `[mm:ss.xx]` format.
- Sorts lines by time.

### `LoadFromFile`

```go
func LoadFromFile(songPath, musicDir string) (Lyrics, bool)
```

Attempts to load a local `.lrc` file.

- Looks in a `lyrics` subdirectory within `musicDir`.
- Tries multiple naming variations (exact match, clean name, lowercase).

### `FetchFromAPI`

```go
func FetchFromAPI(artist, title, album string) (string, error)
```

Fetches synced lyrics from `https://lrclib.net/api/search`.

- Tries precise search with artist, title, and album.
- Fallback to general query search if precise search fails.

### `SaveToFile`

```go
func SaveToFile(songPath, musicDir, content string) (string, error)
```

Saves the provided content to a `.lrc` file in the `lyrics` subdirectory.

- Creates the directory if it doesn't exist.
