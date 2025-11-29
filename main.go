package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"Player/internal/app"
	"Player/service"
)

func main() {
	sourceDir := flag.String("sd", "", "Source directory to load music files")
	flag.Parse()

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
