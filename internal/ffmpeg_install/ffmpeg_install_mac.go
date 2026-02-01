//go:build darwin

package ffmpeginstall

import (
	"fmt"
)

// InstallFFmpeg handles the installation of FFmpeg on macOS systems.
func InstallFFmpeg() error {
	return fmt.Errorf("automatic installation is not yet supported on macOS. Please install FFmpeg manually using 'brew install ffmpeg'")
}
