package main

import (
	"fmt"
)

type WriteCounter struct {
	Total      int64
	Downloaded int64
	ui         *UpdaterUI
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Downloaded += int64(n)
	percentage := float64(wc.Downloaded) / float64(wc.Total)

	wc.ui.progressBar.SetValue(percentage)
	wc.ui.UpdateStatus(fmt.Sprintf("Downloading... %.2f%%", percentage*100))

	return n, nil
}
