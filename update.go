package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// 업데이트 시작
func updateProcess(serverName, fromVersion string, progress *widget.ProgressBar, status *widget.Label, window fyne.Window) {
	updateUI(window, func() {
		status.SetText("Starting update process...")
	})

	// 서버 주소 가져오기
	serverIP, err := getServerIP(serverName)
	if err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error getting server IP: %v", err))
		})
		return
	}

	// 어플리케이션 닫힐 때 까지 기다리기
	waitForApplicationToClose(applicationName, status, window)

	// 파일 백업
	if err := moveFiles(status, window); err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error moving files: %v", err))
		})
		return
	}

	// 다운로드 받을 파일명 가져오기
	fileName, err := getFileNameFromServer(serverIP)
	if err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("error get file name: %v", err))
		})
		return
	}

	// 파일 다운로드
	downloadFilePath := fileName
	downloadURL := fmt.Sprintf("%s/update/file", serverIP)
	if err := downloadFile(downloadURL, downloadFilePath, progress, status, window); err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error downloading file: %v", err))
		})
		return
	}

	// Unzip the downloaded file
	updateUI(window, func() {
		status.SetText("Extracting files...")
	})
	if err := unzipFile(downloadFilePath, "."); err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error extracting files: %v", err))
		})
		return
	}

	if err := verifyHashSum(status, window); err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error verifyHashSum: %v", err))
		})
		return
	}

	updateUI(window, func() {
		status.SetText("Cleaning up...")
	})

	if err := removeFile(downloadFilePath); err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Warning: Failed to delete file: %v", err))
		})
	}

	// Launch the application
	launchApplication(fromVersion, status, window)
}
