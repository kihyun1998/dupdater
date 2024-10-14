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
	startTime  time.Time
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Written += int64(n)
	percentage := float64(wc.Written) / float64(wc.Total)

	// 100ms마다 또는 100KiB마다 UI 업데이트 (둘 중 먼저 도달하는 조건)
	if time.Since(wc.lastUpdate) > 100*time.Millisecond || wc.Written-int64(wc.progress.Value*float64(wc.Total)) > 102400 {
		elapsedTime := time.Since(wc.startTime).Seconds()
		speed := float64(wc.Written) / elapsedTime / 1024 // KiB/s

		updateUI(wc.window, func() {
			wc.progress.SetValue(percentage)
			wc.status.SetText(fmt.Sprintf("다운로드 중... %.2f%% (%.2f KiB/s)", percentage*100, speed))
		})
		wc.lastUpdate = time.Now()
	}

	return n, nil
}
