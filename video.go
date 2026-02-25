package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
)

func getVideoAspectRatio(filePath string) (string, error) {
	// 1. Setup the command
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)

	// 2. Capture the output
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	// 3. Run the command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffprobe error: %v", err)
	}

	// 4. Define a struct to match the JSON output
	type ffprobeOutput struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}

	var data ffprobeOutput
	if err := json.Unmarshal(stdout.Bytes(), &data); err != nil {
		return "", fmt.Errorf("json unmarshal error: %v", err)
	}
	if len(data.Streams) == 0 {
		return "", fmt.Errorf("no streams found in video")
	}

	// 5. Calculate the ratio
	width := data.Streams[0].Width
	height := data.Streams[0].Height

	// use float division to check ratios
	ratio := float64(width) / float64(height)

	// Tolerance for 16:9 (~1.77) and 9:16 (~0.56)
	const tolerance = 0.01
	if math.Abs(ratio-(16.0/9.0)) < tolerance {
		return "landscape", nil
	} else if math.Abs(ratio-(9.0/16.0)) < tolerance {
		return "portrait", nil
	}
	return "other", nil
}
