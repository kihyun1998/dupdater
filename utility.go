package main

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"

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

// func waitForApplicationToClose(appName string, status *widget.Label, window fyne.Window) {
// 	for healthCheck(appName) {
// 		updateUI(window, func() {
// 			status.SetText("Waiting for application to close...")
// 		})
// 		time.Sleep(2 * time.Second)
// 	}
// }

func waitForApplicationToClose(appName string, ui *UpdaterUI) {
	for healthCheck(appName) {
		ui.UpdateStatus("Waiting for application to close...")
		time.Sleep(2 * time.Second)
	}
}

// 프로그램 실행 함수
func launchApplication(fromVersion string, ui *UpdaterUI) {

	ui.UpdateStatus("Preparing to launch application...")

	cmd := exec.Command(fmt.Sprintf("./%s", applicationName), "--patch", "--fromVersion", fromVersion)

	// 창을 보이게 설정
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_NEW_CONSOLE,
	}

	ui.UpdateStatus("Launching client.exe...")
	err := cmd.Start()
	if err != nil {
		ui.ShowError(fmt.Errorf("failed to launch application: %v", err))
		return
	}

	time.Sleep(2 * time.Second)
	ui.window.Close()
}

// func updateUI(window fyne.Window, f func()) {
// 	window.Canvas().Refresh(window.Content())
// 	f()
// }

// func showErrorDialog(err error, window fyne.Window) <-chan struct{} {
// 	done := make(chan struct{})
// 	content := widget.NewLabel(fmt.Sprintf("An error occurred: %v\n\nClick OK to close the updater.", err))
// 	errorDialog := dialog.NewCustom("Error", "OK", content, window)
// 	errorDialog.SetOnClosed(func() {
// 		close(done)
// 	})
// 	errorDialog.Show()
// 	return done
// }
