# service/pickFolder.go

## Overview

Provides Windows-specific functionality to open a folder selection dialog using the Windows API (`shell32.dll`).

## Constants

- `MAX_PATH`: Maximum path length (260).
- `BIF_RETURNONLYFSDIRS`: Flag to return only file system directories.
- `BIF_NEWDIALOGSTYLE`: Flag to use the new dialog style.

## Types

### `BROWSEINFOW`

Struct representing the `BROWSEINFO` structure used by `SHBrowseForFolderW`.

- `hwndOwner`: Handle to the owner window.
- `pidlRoot`: PIDL specifying the location of the root folder.
- `pszDisplayName`: Address of a buffer to receive the display name.
- `lpszTitle`: Address of a null-terminated string that is displayed above the tree view control.
- `ulFlags`: Flags specifying the options for the dialog box.
- `lpfn`: Address of an application-defined callback function.
- `lParam`: Application-defined value that the dialog box passes to the callback function.
- `iImage`: Image associated with the selected folder.

## Functions

### `PickFolder`

```go
func PickFolder() (string, error)
```

Opens a native Windows folder picker dialog.

**Logic:**

1.  **Prepare Structures**: Sets up the `BROWSEINFOW` struct with appropriate flags and buffers.
2.  **Call API**: Invokes `SHBrowseForFolderW` via syscall.
3.  **Handle Result**:
    - If successful, it returns a PIDL (pointer to an item identifier list).
    - If failed or cancelled, returns an error.
4.  **Get Path**: Calls `SHGetPathFromIDListW` to convert the PIDL to a file system path.
5.  **Cleanup**: Frees the PIDL using `CoTaskMemFree`.
6.  **Return**: Returns the selected directory path as a string.
