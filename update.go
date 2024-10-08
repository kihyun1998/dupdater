package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func updateProcess(progress *widget.ProgressBar, status *widget.Label, window fyne.Window) {
	updateUI(window, func() {
		status.SetText("Starting update process...")
	})

	waitForApplicationToClose(applicationName, status, window)

	if err := moveFiles(status, window); err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error moving files: %v", err))
		})
		return
	}

	downloadFilePath := "downloaded_file.zip"
	if err := downloadFile("http://localhost:8000/update/file", downloadFilePath, progress, status, window); err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error downloading file: %v", err))
		})
		return
	}

	// Unzip the downloaded file
	updateUI(window, func() {
		status.SetText("Extracting files...")
	})
	if err := unzipFile(downloadFilePath, "."); err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error extracting files: %v", err))
		})
		return
	}

	// Launch the application
	launchApplication(status, window)
}
