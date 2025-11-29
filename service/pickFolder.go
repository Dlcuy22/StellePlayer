package service

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	MAX_PATH             = 260
	BIF_RETURNONLYFSDIRS = 0x0001
	BIF_NEWDIALOGSTYLE   = 0x0040
)

var (
	shell32                 = windows.NewLazySystemDLL("shell32.dll")
	procSHBrowseForFolderW  = shell32.NewProc("SHBrowseForFolderW")
	procSHGetPathFromIDList = shell32.NewProc("SHGetPathFromIDListW")

	ole32             = windows.NewLazySystemDLL("ole32.dll")
	procCoTaskMemFree = ole32.NewProc("CoTaskMemFree")
)

// BROWSEINFOW marshalling layout for SHBrowseForFolderW
type BROWSEINFOW struct {
	hwndOwner      uintptr
	pidlRoot       uintptr
	pszDisplayName uintptr
	lpszTitle      uintptr
	ulFlags        uint32
	lpfn           uintptr
	lParam         uintptr
	iImage         int32
}

func PickFolder() (string, error) {
	// buffer for display name & returned path
	buf := make([]uint16, MAX_PATH)

	// title for dialog
	title := syscall.StringToUTF16Ptr("Select a folder")

	bi := BROWSEINFOW{
		hwndOwner:      0,
		pidlRoot:       0,
		pszDisplayName: uintptr(unsafe.Pointer(&buf[0])),
		lpszTitle:      uintptr(unsafe.Pointer(title)),
		ulFlags:        BIF_RETURNONLYFSDIRS | BIF_NEWDIALOGSTYLE,
		lpfn:           0,
		lParam:         0,
		iImage:         0,
	}

	// Call SHBrowseForFolderW
	pidl, _, err := procSHBrowseForFolderW.Call(uintptr(unsafe.Pointer(&bi)))
	if pidl == 0 {
		return "", fmt.Errorf("cancelled or no folder selected")
	}
	if err != nil && err != syscall.Errno(0) {
		// still try to free pidl
		procCoTaskMemFree.Call(pidl)
		return "", fmt.Errorf("SHBrowseForFolderW error: %w", err)
	}

	// Prepare buffer for path
	pathBuf := make([]uint16, MAX_PATH)

	ret, _, err := procSHGetPathFromIDList.Call(pidl, uintptr(unsafe.Pointer(&pathBuf[0])))
	if ret == 0 {
		procCoTaskMemFree.Call(pidl)
		if err != nil && err != syscall.Errno(0) {
			return "", fmt.Errorf("SHGetPathFromIDListW failed: %w", err)
		}
		return "", fmt.Errorf("SHGetPathFromIDListW failed")
	}

	// Convert to Go string
	path := syscall.UTF16ToString(pathBuf)
	// free pidl
	procCoTaskMemFree.Call(pidl)

	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("no folder selected")
	}

	return path, nil
}
