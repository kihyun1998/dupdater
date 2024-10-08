package main

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
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
func launchApplication(status *widget.Label, window fyne.Window) {
	updateUI(window, func() {
		status.SetText("Launching client.exe...")
	})

	cmd := exec.Command(fmt.Sprintf("./%s", applicationName))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
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
