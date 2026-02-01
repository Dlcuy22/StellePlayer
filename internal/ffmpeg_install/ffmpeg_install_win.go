//go:build windows

package ffmpeginstall

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
)

const (
	ffmpegDlUrl = "https://www.gyan.dev/ffmpeg/builds/ffmpeg-git-essentials.7z"
)

// InstallFFmpeg handles the installation of FFmpeg on Windows systems.
// It first tries 'winget', then falls back to a manual download and PATH update.
func InstallFFmpeg() error {
	// 1. Try Winget first
	fmt.Println("Attempting to install FFmpeg via Winget...")
	if err := installViaWinget(); err == nil {
		fmt.Println("FFmpeg installed successfully via Winget.")
		return nil
	} else {
		fmt.Printf("Winget installation failed: %v. Falling back to manual download...\n", err)
	}

	// 2. Manual Download & Install
	return installManual()
}

func installViaWinget() error {
	cmd := exec.Command("winget", "install", "ffmpeg", "--accept-source-agreements", "--accept-package-agreements")
	return cmd.Run()
}

func installManual() error {
	// Define install directory: %LOCALAPPDATA%\StellePlayer\ffmpeg
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return fmt.Errorf("LOCALAPPDATA environment variable not found")
	}
	installDir := filepath.Join(localAppData, "StellePlayer", "ffmpeg")
	binDir := filepath.Join(installDir, "bin")

	// Create directory
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	// Download
	archivePath := filepath.Join(installDir, "ffmpeg.7z")
	fmt.Println("Downloading FFmpeg from gyan.dev...")
	if err := downloadFile(ffmpegDlUrl, archivePath); err != nil {
		return fmt.Errorf("failed to download FFmpeg: %w", err)
	}
	defer os.Remove(archivePath) // Cleanup archive

	// Extract
	fmt.Println("Extracting FFmpeg...")
	if err := extract7z(archivePath, installDir); err != nil {
		return fmt.Errorf("failed to extract FFmpeg: %w", err)
	}

	// Locate the bin folder inside the extracted structure
	// The archive usually has a root folder like "ffmpeg-2023-..."
	var foundBin bool
	err := filepath.Walk(installDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "bin" {
			// Check if ffmpeg.exe is inside
			if _, err := os.Stat(filepath.Join(path, "ffmpeg.exe")); err == nil {
				// Move contents to our main binDir or just update PATH to this?
				// Let's simpler: Update PATH to this detected bin directory.
				binDir = path
				foundBin = true
				return io.EOF // Stop walking
			}
		}
		return nil
	})

	if !foundBin && err != io.EOF {
		return fmt.Errorf("could not find 'bin' directory in extracted files")
	}

	// Add to User PATH
	fmt.Println("Adding FFmpeg to User PATH...")
	if err := addToPath(binDir); err != nil {
		return fmt.Errorf("failed to update PATH: %w", err)
	}

	fmt.Println("FFmpeg installed manually. Please restart the application.")
	return nil
}

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func extract7z(archivePath, dest string) error {
	r, err := sevenzip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func addToPath(newPath string) error {
	// Get current User PATH
	cmd := exec.Command("powershell", "-Command", "[Environment]::GetEnvironmentVariable('PATH', 'User')")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	currentPath := strings.TrimSpace(string(output))

	// Check if already exists
	if strings.Contains(currentPath, newPath) {
		return nil
	}

	// Update PATH
	newPathVal := currentPath + ";" + newPath
	// Escape quotes for Powershell command safety if needed, though simple paths are usually fine.
	updateCmd := fmt.Sprintf("[Environment]::SetEnvironmentVariable('PATH', '%s', 'User')", newPathVal)

	cmd = exec.Command("powershell", "-Command", updateCmd)
	return cmd.Run()
}
