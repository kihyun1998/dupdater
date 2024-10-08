package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var applicationName string

func init() {
	applicationName = "update_test_app_1.exe"
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("File Downloader")

	progress := widget.NewProgressBar()
	status := widget.NewLabel("Checking for updates...")

	content := container.NewVBox(
		progress,
		status,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(300, 100))

	go func() {
		for {
			if !healthCheck(applicationName) {
				updateProcess(progress, status, myWindow)
				break
			}
			updateUI(myWindow, func() {
				status.SetText("Waiting for application to close...")
			})
			time.Sleep(2 * time.Second)
		}
	}()

	myWindow.ShowAndRun()
}
