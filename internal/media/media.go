package media

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Metadata struct {
	Title      string
	Artist     string
	Album      string
	Duration   float64
	FilePath   string
	Bitrate    string
	Codec      string
	SampleRate string
}

func LoadFromDirectory(dir string) ([]Metadata, error) {
	var tracks []Metadata
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

		meta, err := extractMetadata(path)
		if err != nil {
			return nil
		}
		tracks = append(tracks, meta)
		return nil
	})

	return tracks, err
}

func extractMetadata(filePath string) (Metadata, error) {
	data, err := ffmpeg.Probe(filePath)
	if err != nil {
		return Metadata{}, err
	}

	var probeData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &probeData); err != nil {
		return Metadata{}, err
	}

	meta := Metadata{
		FilePath:   filePath,
		Title:      filepath.Base(filePath),
		Artist:     "Unknown Artist",
		Album:      "Unknown Album",
		Bitrate:    "N/A",
		Codec:      "N/A",
		SampleRate: "N/A",
	}

	if format, ok := probeData["format"].(map[string]interface{}); ok {
		if duration, ok := format["duration"].(string); ok {
			fmt.Sscanf(duration, "%f", &meta.Duration)
		}
		if bitrate, ok := format["bit_rate"].(string); ok {
			bitrateInt := 0
			fmt.Sscanf(bitrate, "%d", &bitrateInt)
			if bitrateInt > 0 {
				meta.Bitrate = fmt.Sprintf("%d kbps", bitrateInt/1000)
			}
		}
		if tags, ok := format["tags"].(map[string]interface{}); ok {
			if title, ok := tags["title"].(string); ok {
				meta.Title = title
			}
			if artist, ok := tags["artist"].(string); ok {
				meta.Artist = artist
			}
			if album, ok := tags["album"].(string); ok {
				meta.Album = album
			}
		}
	}

	if streams, ok := probeData["streams"].([]interface{}); ok {
		for _, stream := range streams {
			streamMap, ok := stream.(map[string]interface{})
			if !ok {
				continue
			}
			if codecType, ok := streamMap["codec_type"].(string); !ok || codecType != "audio" {
				continue
			}
			if codecName, ok := streamMap["codec_name"].(string); ok {
				meta.Codec = strings.ToUpper(codecName)
			}
			if sampleRate, ok := streamMap["sample_rate"].(string); ok {
				sampleRateInt := 0
				fmt.Sscanf(sampleRate, "%d", &sampleRateInt)
				if sampleRateInt > 0 {
					meta.SampleRate = fmt.Sprintf("%.1f kHz", float64(sampleRateInt)/1000.0)
				}
			}
			if meta.Bitrate == "N/A" {
				if bitrate, ok := streamMap["bit_rate"].(string); ok {
					bitrateInt := 0
					fmt.Sscanf(bitrate, "%d", &bitrateInt)
					if bitrateInt > 0 {
						meta.Bitrate = fmt.Sprintf("%d kbps", bitrateInt/1000)
					}
				}
			}
			if tags, ok := streamMap["tags"].(map[string]interface{}); ok {
				if title, ok := tags["title"].(string); ok && meta.Title == filepath.Base(filePath) {
					meta.Title = title
				}
				if artist, ok := tags["artist"].(string); ok && meta.Artist == "Unknown Artist" {
					meta.Artist = artist
				}
				if album, ok := tags["album"].(string); ok && meta.Album == "Unknown Album" {
					meta.Album = album
				}
			}
			break
		}
	}

	return meta, nil
}
