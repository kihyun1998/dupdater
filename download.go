package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// 다운로드 함수
func downloadFile(url, downloadFilePath string, progress *widget.ProgressBar, status *widget.Label, window fyne.Window) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error downloading file: %v", err))
		})
		return &exec.ExitError{}
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Bad status: %s", resp.Status))
		})
		return nil
	}

	// Create the file
	out, err := os.Create(downloadFilePath)
	if err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error creating file: %v", err))
		})
		return err
	}
	defer out.Close()

	// Create a custom io.Writer to track progress
	counter := &WriteCounter{
		Total:    resp.ContentLength,
		progress: progress,
		status:   status,
		window:   window,
	}

	// Start time for speed calculation
	startTime := time.Now()

	// Use io.Copy to optimize file writing and update progress
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error writing to file: %v", err))
		})
		return err
	}

	elapsedTime := time.Since(startTime).Seconds()
	speed := float64(counter.Written) / elapsedTime / 1024 // KB/s

	updateUI(window, func() {
		progress.SetValue(1)
		status.SetText(fmt.Sprintf("Download completed (%.2f KB/s)", speed))
	})

	time.Sleep(2 * time.Second)
	return nil
}
