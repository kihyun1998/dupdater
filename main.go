package main

import (
	"flag"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
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
func showMainWindow() {
	ui := newUpdaterUI(fmt.Sprintf("ACRA Point Update - %s", fromVersion))

	go func() {
		for {
			if !healthCheck(applicationName) {
				updateProcess(serverName, fromVersion, ui)
				break
			}
			ui.UpdateDetail("Waiting for application to close...")
			time.Sleep(2 * time.Second)
		}
	}()

	ui.Run()
}

func main() {
	if fromVersion == "" || serverName == "" {
		ui := newUpdaterUI("Error")
		ui.ShowError(fmt.Errorf("Error: -fromVersion flag and -serverName flag are required"))
		ui.Run()
		return
	}

	if err := InitLogger(); err != nil {
		ui := newUpdaterUI("Error")
		ui.ShowError(fmt.Errorf("Failed to initialize logger: %v", err))
		ui.Run()
		return
	}
	defer CloseLogger()

	LogInfo("Application started. From Version: %s, Server Name: %s", fromVersion, serverName)

	showMainWindow()
}
