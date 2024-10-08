package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func healthCheck(appName string) bool {
	cmd := exec.Command("tasklist")
	output, err := cmd.Output() 
	if err != nil {
		fmt.Println("Error checking running processes:", err)
		return false
	}
	return strings.Contains(string(output), appName)
}

// 프로그램 실행 함수
func launchApplication(status *widget.Label, window fyne.Window) {
	updateUI(window, func() {
		status.SetText("Launching client.exe...")
	})

	cmd := exec.Command(fmt.Sprintf("./%s", applicationName))
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
