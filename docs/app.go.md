# internal/app/app.go

## Overview

This file contains the entry point for the `app` package, responsible for initializing and running the Bubble Tea program.

## Functions

### `Run`

```go
func Run(musicDir string) error
```

Starts the Bubble Tea program and initiates the song loading process.

**Parameters:**

- `musicDir` (string): The path to the directory containing music files.

**Returns:**

- `error`: Returns an error if the program fails to run.

**Logic:**

1.  **Initialize Program**: Creates a new `tea.Program` using `initialModel(musicDir)` and enables the alternate screen.
2.  **Async Loading**: Starts a goroutine to load songs in the background:
    - Calls `media.LoadFromDirectory(musicDir)` to get metadata.
    - Handles errors (prints error and quits).
    - Checks if no songs were found (prints message and quits).
    - Converts `engine.Metadata` to `Song` structs.
    - Sends a `songsLoadedMsg` to the program.
3.  **Run Program**: Calls `program.Run()` to start the UI loop.
