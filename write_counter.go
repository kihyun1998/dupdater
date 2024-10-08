package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type WriteCounter struct {
	Total      int64
	Written    int64
	progress   *widget.ProgressBar
	status     *widget.Label
	window     fyne.Window
	lastUpdate time.Time
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Written += int64(n)
	percentage := float64(wc.Written) / float64(wc.Total)

	// Update UI every 100ms to reduce update frequency
	if time.Since(wc.lastUpdate) > 100*time.Millisecond {
		updateUI(wc.window, func() {
			wc.progress.SetValue(percentage)
			speed := float64(wc.Written) / time.Since(wc.lastUpdate).Seconds() / 1024 // KB/s
			wc.status.SetText(fmt.Sprintf("Downloading... %.2f%% (%.2f KB/s)", percentage*100, speed))
		})
		wc.lastUpdate = time.Now()
	}

	return n, nil
}
