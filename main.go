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
var fromVersion string
var serverName string

func init() {
	// 실행할 어플리케이션 이름 정의
	applicationName = "update_test_app_1.exe"
	// from version 인자값
	flag.StringVar(&fromVersion, "fromVersion", "", "App's From Version (required)")
	// 프로파일명 인자값
	flag.StringVar(&serverName, "serverName", "", "Server Name (required)")
	flag.Parse()

}

// 인자값 올바르지 않은경우의 창
func showErrorAndExit(myApp fyne.App) {
	errorWindow := myApp.NewWindow("Error")
	errorLabel := widget.NewLabel("Error: -fromVersion flag is required and -serverName flag is required")
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

// 제대로 된 경우의 창
func showMainWindow(myApp fyne.App) {
	mainWindow := myApp.NewWindow(fmt.Sprintf("ACRA Point Update - %s", fromVersion))

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
				updateProcess(serverName, fromVersion, progress, status, mainWindow)
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

	if fromVersion == "" || serverName == "" {
		showErrorAndExit(myApp)
		return
	}

	if err := InitLogger(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}
	defer CloseLogger()

	LogInfo("Application started. From Version: %s, Server Name: %s", fromVersion, serverName)

	showMainWindow(myApp)
}
