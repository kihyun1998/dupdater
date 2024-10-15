package main

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
)

// 프로그램 health check
func healthCheck(appName string) (result bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			LogError("Panic in healthCheck: %v", r)
			err = fmt.Errorf("healthCheck failed unexpectedly: %v", r)
		}
	}()

	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s", appName))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error checking running processes:", err)
		return false, nil
	}
	return strings.Contains(string(output), appName), nil
}

func waitForApplicationToClose(appName string, ui *UpdaterUI) (err error) {
	defer func() {
		if r := recover(); r != nil {
			LogError("Panic in waitForApplicationToClose: %v", r)
			err = fmt.Errorf("waitForApplicationToClose failed unexpectedly: %v", r)
		}
	}()

	const maxAttempts = 30
	attempt := 0

	for {
		isRunning, err := healthCheck(appName)
		if err != nil {
			return fmt.Errorf("error checking application status:%v", err)
		}

		if !isRunning {
			ui.UpdateDetail("Application has closed successfully.")
			return nil
		}

		attempt++
		if attempt >= maxAttempts {
			return fmt.Errorf("timeout waiting for application to close after %d seconds", maxAttempts*2)
		}

		ui.UpdateDetail(fmt.Sprintf("Waiting for application to close... (Attempt %d/%d)", attempt, maxAttempts))
	}
}

// 프로그램 실행 함수
func launchApplication(fromVersion string, ui *UpdaterUI) (err error) {
	defer func() {
		if r := recover(); r != nil {
			LogError("Panic in launchApplication: %v", r)
			err = fmt.Errorf("launchApplication failed unexpectedly: %v", r)
		}
	}()

	ui.UpdateDetail("Preparing to launch application...")

	cmd := exec.Command(fmt.Sprintf("./%s", applicationName), "--patch", "--fromVersion", fromVersion)

	// 창을 보이게 설정
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_NEW_CONSOLE,
	}

	ui.UpdateDetail("Launching client.exe...")
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to launch application: %v", err)
	}

	time.Sleep(2 * time.Second)
	ui.window.Close()

	return nil
}
