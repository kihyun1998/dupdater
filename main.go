package main

import (
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
	status := widget.NewLabel("Ready to start")

	button := widget.NewButton("Start Update", func() {
		go updateProcess(progress, status, myWindow)
	})

	content := container.NewVBox(
		progress,
		status,
		button,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(300, 100))
	myWindow.ShowAndRun()
}
