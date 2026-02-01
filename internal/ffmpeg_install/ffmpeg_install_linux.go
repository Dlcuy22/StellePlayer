//go:build linux

package ffmpeginstall

import (
	"fmt"
	"os/exec"
)

// InstallFFmpeg handles the installation of FFmpeg on Linux systems.
// It detects the distribution and runs the appropriate package manager command.
func InstallFFmpeg() error {
	distro, err := getDistroID()
	if err != nil {
		return fmt.Errorf("failed to detect Linux distribution: %w", err)
	}

	var cmd *exec.Cmd

	switch distro {
	case "ubuntu", "debian", "linuxmint", "pop":
		cmd = exec.Command("sudo", "apt-get", "update")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to update apt repositories: %w", err)
		}
		cmd = exec.Command("sudo", "apt-get", "install", "-y", "ffmpeg")
	case "fedora":
		cmd = exec.Command("sudo", "dnf", "install", "-y", "ffmpeg")
	case "arch", "manjaro":
		cmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", "ffmpeg")
	case "opensuse-leap", "opensuse-tumbleweed":
		cmd = exec.Command("sudo", "zypper", "install", "-y", "ffmpeg")
	default:
		return fmt.Errorf("unsupported distribution: %s. Please install FFmpeg manually", distro)
	}

	fmt.Printf("Installing FFmpeg for %s...\n", distro)
	cmd.Stdout = exec.Command("cat").Stdout // Pipe to stdout for user visibility (simplified)
	cmd.Stderr = exec.Command("cat").Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install FFmpeg: %w", err)
	}

	return nil
}

// getDistroID attempts to read the ID field from /etc/os-release
func getDistroID() (string, error) {
	cmd := exec.Command("sh", "-c", "source /etc/os-release && echo $ID")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output)[:len(string(output))-1], nil // Trim newline
}
