//go:build !windows

package service

import "fmt"

func PickFolder() (string, error) {
	return "", fmt.Errorf("interactive folder picker is only available on Windows. Please use the -sd flag to specify the music directory")
}
