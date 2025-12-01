package media

import (
	"Player/internal/engine"
	"os"
	"path/filepath"
	"strings"
)

func LoadFromDirectory(dir string) ([]engine.Metadata, error) {
	var tracks []engine.Metadata
	supportedExts := map[string]bool{
		".mp3":  true,
		".m4a":  true,
		".flac": true,
		".wav":  true,
		".ogg":  true,
		".aac":  true,
		".opus": true,
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !supportedExts[ext] {
			return nil
		}

		meta, err := engine.ExtractMetadata(path)
		if err != nil {
			return nil
		}
		tracks = append(tracks, meta)
		return nil
	})

	return tracks, err
}
