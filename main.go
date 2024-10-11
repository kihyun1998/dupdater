package main

import (
	"flag"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var applicationName string
var currentVersion string
var serverName string

func init() {
	applicationName = "update_test_app_1.exe"
	flag.StringVar(&currentVersion, "currentVersion", "", "App's Current Version (required)")
	flag.StringVar(&serverName, "serverName", "", "Server Name (required)")
	flag.Parse()

}

func showErrorAndExit(myApp fyne.App) {
	errorWindow := myApp.NewWindow("Error")
	errorLabel := widget.NewLabel("Error: -currentVersion flag is required and -serverName flag is required")
	okButton := widget.NewButton("OK", func() {
		errorWindow.Close()
		myApp.Quit()
	})

	content := container.NewVBox(
		errorLabel,
		okButton,
	)

	errorWindow.SetContent(content)
	errorWindow.Resize(fyne.NewSize(500, 100))
	errorWindow.SetOnClosed(func() {
		myApp.Quit()
	})

	errorWindow.CenterOnScreen()
	errorWindow.Show()

	myApp.Run()
}

func showMainWindow(myApp fyne.App) {
	mainWindow := myApp.NewWindow(fmt.Sprintf("ACRA Point Update - %s", currentVersion))

	progress := widget.NewProgressBar()
	status := widget.NewLabel("Checking for updates...")

	content := container.NewVBox(
		progress,
		status,
	)

	mainWindow.SetContent(content)
	mainWindow.Resize(fyne.NewSize(500, 100))

	go func() {
		for {
			if !healthCheck(applicationName) {
				updateProcess(serverName, progress, status, mainWindow)
				break
			}
			updateUI(mainWindow, func() {
				status.SetText("Waiting for application to close...")
			})
		}
	}()

	mainWindow.CenterOnScreen()
	mainWindow.ShowAndRun()
}

func main() {
	myApp := app.New()

	if currentVersion == "" || serverName == "" {
		showErrorAndExit(myApp)
		return
	}

	showMainWindow(myApp)
}
