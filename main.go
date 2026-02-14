package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"Player/internal/app"
	ffmpeginstall "Player/internal/ffmpeg_install"
	"Player/service"
	"path/filepath"
)

var Version = "dev"

func main() {
	sourceDir := flag.String("sd", "", "Source directory to load music files")
	installFFmpegFlag := flag.Bool("install-ffmpeg", false, "Force FFmpeg installation prompt")
	customFFmpegDir := flag.String("use-custom-ffmpeg", "", "Path to directory containing custom FFmpeg binaries")
	versionFlag := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("StellePlayer version %s\n", Version)
		os.Exit(0)
	}

	// Prepend custom FFmpeg directory to PATH if provided
	if *customFFmpegDir != "" {
		absDir, err := filepath.Abs(*customFFmpegDir)
		if err != nil {
			fmt.Printf("Error resolving custom FFmpeg path: %v\n", err)
			os.Exit(1)
		}
		path := os.Getenv("PATH")
		os.Setenv("PATH", absDir+string(os.PathListSeparator)+path)
		fmt.Printf("Using custom FFmpeg from: %s\n", absDir)
	}

	// Force install mode for testing
	if *installFFmpegFlag {
		fmt.Println("Manual FFmpeg installation mode triggered.")
		runFFmpegInstall()
		return
	}

	// Check for FFmpeg dependencies
	if !isFFmpegInstalled() {
		handleFFmpegMissing()
	}

	musicDir := strings.TrimSpace(*sourceDir)
	if musicDir == "" {
		fmt.Println("No -sd flag provided. Please select a folder containing your music files.")
		selectedDir, err := service.PickFolder()
		if err != nil {
			fmt.Printf("Failed to pick folder: %v\n", err)
			os.Exit(1)
		}
		musicDir = strings.TrimSpace(selectedDir)
		if musicDir == "" {
			fmt.Println("No folder selected. Exiting.")
			os.Exit(1)
		}
	}

	if err := app.Run(musicDir); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// isFFmpegInstalled checks if ffplay and ffprobe are available in PATH.
func isFFmpegInstalled() bool {
	_, errPlay := exec.LookPath("ffplay")
	_, errProbe := exec.LookPath("ffprobe")
	return errPlay == nil && errProbe == nil
}

// handleFFmpegMissing prompts the user and attempts installation.
func handleFFmpegMissing() {
	fmt.Println("⚠️  FFmpeg (ffplay, ffprobe) not found in your PATH.")
	fmt.Println("This application requires FFmpeg to play music and read metadata.")
	fmt.Print("Would you like to attempt automatic installation? (y/n): ")

	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer != "y" && answer != "yes" {
		fmt.Println("Please install FFmpeg manually and add it to your PATH.")
		fmt.Println("Exiting.")
		os.Exit(1)
	}

	fmt.Println("Starting FFmpeg installation...")
	if err := ffmpeginstall.InstallFFmpeg(); err != nil {
		fmt.Printf("❌ Installation failed: %v\n", err)
		fmt.Println("Please install FFmpeg manually.")
		os.Exit(1)
	}

	fmt.Println("✅ FFmpeg installation complete.")
	fmt.Println("Please restart the application.")
	os.Exit(0)
}

// runFFmpegInstall runs the installation directly without prompts (for --install-ffmpeg flag).
func runFFmpegInstall() {
	fmt.Println("Starting FFmpeg installation...")
	if err := ffmpeginstall.InstallFFmpeg(); err != nil {
		fmt.Printf("❌ Installation failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ FFmpeg installation complete.")
}
