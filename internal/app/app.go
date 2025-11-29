package app

import (
	"fmt"
	"os"

	"Player/internal/media"

	tea "github.com/charmbracelet/bubbletea"
)

// Run starts the Bubble Tea program and loads songs from the provided directory.
func Run(musicDir string) error {
	program := tea.NewProgram(initialModel(musicDir), tea.WithAltScreen())

	go func() {
		metas, err := media.LoadFromDirectory(musicDir)
		if err != nil {
			fmt.Printf("\nError loading songs: %v\n", err)
			program.Quit()
			os.Exit(1)
		}

		if len(metas) == 0 {
			fmt.Println("\nNo music files found in the specified directory")
			fmt.Println("Supported formats: .mp3, .m4a, .flac, .wav, .ogg, .aac")
			program.Quit()
			os.Exit(1)
		}

		songs := make([]Song, len(metas))
		for i, meta := range metas {
			songs[i] = Song{metadata: meta}
		}

		program.Send(songsLoadedMsg{songs: songs, musicDir: musicDir})
	}()

	if _, err := program.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}
