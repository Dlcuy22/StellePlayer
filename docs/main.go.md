# main.go

## Overview

The entry point of the application. It handles command-line arguments, ensures a music directory is selected, and starts the application.

## Functions

### `main`

```go
func main()
```

The main function of the program.

**Steps:**

1.  **Flag Parsing**: Defines and parses the `-sd` flag for the source directory.
2.  **Directory Validation**:
    - Checks if `sourceDir` is provided.
    - If not, calls `service.PickFolder()` to prompt the user to select a folder.
    - Validates that a directory is selected; exits if not.
3.  **Application Launch**: Calls `app.Run(musicDir)` to start the music player.
4.  **Error Handling**: Prints any errors returned by `app.Run` and exits with status 1.
