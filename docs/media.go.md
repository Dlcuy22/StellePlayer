# internal/media/media.go

## Overview

Responsible for scanning directories and loading media files.

## Functions

### `LoadFromDirectory`

```go
func LoadFromDirectory(dir string) ([]engine.Metadata, error)
```

Scans the specified directory recursively for supported audio files.

**Supported Extensions:**

- `.mp3`, `.m4a`, `.flac`, `.wav`, `.ogg`, `.aac`, `.opus`

**Logic:**

1.  Walks the directory tree using `filepath.Walk`.
2.  Checks file extensions against the supported list.
3.  Calls `engine.ExtractMetadata` for each valid file.
4.  Accumulates and returns a slice of `engine.Metadata`.
