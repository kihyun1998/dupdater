package main

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/sys/windows"
)

// / 프로그램 health check
func healthCheck(appName string) bool {
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s", appName))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error checking running processes:", err)
		return false
	}
	return strings.Contains(string(output), appName)
}

func waitForApplicationToClose(appName string, status *widget.Label, window fyne.Window) {
	for healthCheck(appName) {
		updateUI(window, func() {
			status.SetText("Waiting for application to close...")
		})
		time.Sleep(2 * time.Second)
	}
}

// 프로그램 실행 함수
func launchApplication(fromVersion string, status *widget.Label, window fyne.Window) {
	updateUI(window, func() {
		status.SetText("Launching client.exe...")
	})

	cmd := exec.Command(fmt.Sprintf("./%s", applicationName), "--patch", "--fromVersion", fromVersion)

	// 창을 보이게 설정
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_NEW_CONSOLE,
	}

	err := cmd.Start()
	if err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error launching client.exe: %v", err))
		})
	} else {
		updateUI(window, func() {
			status.SetText("client.exe launched successfully")
		})
	}

	time.Sleep(2 * time.Second)
	window.Close()
}

func updateUI(window fyne.Window, f func()) {
	window.Canvas().Refresh(window.Content())
	f()
}

func showErrorDialog(err error, window fyne.Window) <-chan struct{} {
	done := make(chan struct{})
	content := widget.NewLabel(fmt.Sprintf("An error occurred: %v\n\nClick OK to close the updater.", err))
	errorDialog := dialog.NewCustom("Error", "OK", content, window)
	errorDialog.SetOnClosed(func() {
		close(done)
	})
	errorDialog.Show()
	return done
}
