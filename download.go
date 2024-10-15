package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// 다운로드 함수
func downloadFile(url, filepath string, ui *UpdaterUI) error {
	ui.UpdateDetail("Starting download...")

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading file: %v", err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	counter := &WriteCounter{
		Total: resp.ContentLength,
		ui:    ui,
	}

	reader := io.TeeReader(resp.Body, counter)

	// Start time for speed calculation
	startTime := time.Now()

	// Use io.Copy to optimize file writing and update progress
	_, err = io.Copy(out, reader)
	if err != nil {
		return fmt.Errorf("error copying content: %v", err)
	}

	elapsedTime := time.Since(startTime).Seconds()
	speed := float64(counter.Total) / elapsedTime / 1024 / 1024 // MB/s

	ui.UpdateDetail(fmt.Sprintf("Download completed (%.2f MB/s)", speed))

	time.Sleep(time.Second)

	ui.UpdateDetail("Verifying file integrity...")
	if err := verifyFileHash(filepath, ui); err != nil {
		return fmt.Errorf("file verification failed: %v", err)
	}

	return nil
}
