# internal/engine/metadata.go

## Overview

Handles the extraction of metadata from audio files using `ffmpeg-go` (ffprobe).

## Types

### `Metadata`

Struct holding audio file information.

- `Title`, `Artist`, `Album` (string)
- `Duration` (float64): Duration in seconds.
- `FilePath` (string): Absolute path to the file.
- `Bitrate`, `Codec`, `SampleRate` (string): Audio technical details.

## Functions

### `ExtractMetadata`

```go
func ExtractMetadata(filePath string) (Metadata, error)
```

Extracts metadata for the given file.

**Logic:**

1.  Calls `ffmpeg.Probe(filePath)` to get JSON metadata.
2.  Unmarshals JSON into a map.
3.  Parses `format` section for Duration, Bitrate, and Tags (Title, Artist, Album).
4.  Parses `streams` section for Codec, Sample Rate, and fallback Bitrate.
5.  Returns a populated `Metadata` struct.
