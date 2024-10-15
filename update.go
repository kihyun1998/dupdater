package main

import (
	"fmt"
)

func updateProcess(serverName, fromVersion string, ui *UpdaterUI) {
	defer func() {
		if r := recover(); r != nil {
			LogError("panic in updateProcess: %v", r)
			ui.ShowError(fmt.Errorf("update process failed : %v", r))
		}
	}()

	ui.SetCurrentStep(0)
	ui.UpdateDetail("Starting update process...")

	// 서버 주소 가져오기
	serverIP, err := getServerIP(serverName)
	if err != nil {
		ui.ShowError(fmt.Errorf("error getting server IP: %v", err))
		return
	}

	// 어플리케이션 닫힐 때까지 기다리기
	// TODO: 몇초 기다리고 안꺼진거 확인되면 꺼달라고 하기
	ui.SetCurrentStep(1)
	ui.UpdateDetail("Waiting for application to close...")
	waitForApplicationToClose(applicationName, ui)

	// 파일 백업
	ui.SetCurrentStep(2)
	ui.UpdateDetail("Backing up files...")
	if err := moveFiles(ui); err != nil {
		ui.ShowError(fmt.Errorf("error moving files: %v", err))
		return
	}

	// 다운로드 받을 파일명 가져오기
	ui.SetCurrentStep(3)
	ui.UpdateDetail("Getting update file name...")
	fileName, err := getFileNameFromServer(serverIP)
	if err != nil {
		ui.ShowError(fmt.Errorf("error get file name: %v", err))
		return
	}

	// 파일 다운로드
	ui.SetCurrentStep(4)
	ui.UpdateDetail("Downloading update...")
	downloadFilePath := fileName
	downloadURL := fmt.Sprintf("%s/update/file", serverIP)
	if err := downloadFile(downloadURL, downloadFilePath, ui); err != nil {
		ui.ShowError(fmt.Errorf("error downloading file: %v", err))
		return
	}

	// Unzip the downloaded file
	ui.SetCurrentStep(5)
	ui.UpdateDetail("Extracting files...")
	if err := unzipFile(downloadFilePath, "."); err != nil {
		ui.ShowError(fmt.Errorf("error extracting files: %v", err))
		return
	}

	if err := verifyHashSum(ui); err != nil {
		ui.ShowError(fmt.Errorf("error verifying hash sum: %v", err))
		return
	}

	ui.SetCurrentStep(6)
	ui.UpdateDetail("Cleaning up...")

	if err := removeFile(downloadFilePath); err != nil {
		ui.ShowError(fmt.Errorf("Warning: Failed to delete file: %v", err))
		return
	}

	// Launch the application
	if err := launchApplication(fromVersion, ui); err != nil {
		ui.ShowError(fmt.Errorf("%v", err))
		return
	}

}

// // 업데이트 시작
// func updateProcess(serverName, fromVersion string, progress *widget.ProgressBar, status *widget.Label, window fyne.Window) {
// 	updateUI(window, func() {
// 		status.SetText("Starting update process...")
// 	})

// 	// 서버 주소 가져오기
// 	serverIP, err := getServerIP(serverName)
// 	if err != nil {
// 		updateUI(window, func() {
// 			status.SetText(fmt.Sprintf("Error getting server IP: %v", err))
// 		})
// 		return
// 	}

// 	// 어플리케이션 닫힐 때 까지 기다리기
// 	waitForApplicationToClose(applicationName, status, window)

// 	// 파일 백업
// 	if err := moveFiles(status, window); err != nil {
// 		updateUI(window, func() {
// 			status.SetText(fmt.Sprintf("Error moving files: %v", err))
// 		})
// 		return
// 	}

// 	// 다운로드 받을 파일명 가져오기
// 	fileName, err := getFileNameFromServer(serverIP)
// 	if err != nil {
// 		updateUI(window, func() {
// 			status.SetText(fmt.Sprintf("error get file name: %v", err))
// 		})
// 		return
// 	}

// 	// 파일 다운로드
// 	downloadFilePath := fileName
// 	downloadURL := fmt.Sprintf("%s/update/file", serverIP)
// 	if err := downloadFile(downloadURL, downloadFilePath, progress, status, window); err != nil {
// 		updateUI(window, func() {
// 			status.SetText(fmt.Sprintf("Error downloading file: %v", err))
// 		})
// 		return
// 	}

// 	// Unzip the downloaded file
// 	updateUI(window, func() {
// 		status.SetText("Extracting files...")
// 	})
// 	if err := unzipFile(downloadFilePath, "."); err != nil {
// 		updateUI(window, func() {
// 			status.SetText(fmt.Sprintf("Error extracting files: %v", err))
// 		})
// 		return
// 	}

// 	if err := verifyHashSum(status, window); err != nil {
// 		updateUI(window, func() {
// 			status.SetText(fmt.Sprintf("Error verifyHashSum: %v", err))
// 		})
// 		return
// 	}

// 	updateUI(window, func() {
// 		status.SetText("Cleaning up...")
// 	})

// 	if err := removeFile(downloadFilePath); err != nil {
// 		updateUI(window, func() {
// 			status.SetText(fmt.Sprintf("Warning: Failed to delete file: %v", err))
// 		})
// 	}

// 	// Launch the application
// 	launchApplication(fromVersion, status, window)
// }
