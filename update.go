package main

import (
	"fmt"
)

func updateProcess(serverName, fromVersion string, ui *UpdaterUI) {
	backupCompleted := false

	defer func() {
		if r := recover(); r != nil {
			LogError("panic in updateProcess: %v", r)
			if backupCompleted {
				ui.UpdateDetail("Error occurred. Attempting to restore files...")
				if err := restoreFiles(ui); err != nil {
					LogError("Failed to restore files: %v", err)
					ui.ShowError(fmt.Errorf("update and restore failed: %v", r))
				} else {
					ui.ShowError(fmt.Errorf("update failed, files restored: %v", r))
				}
			} else {
				ui.ShowError(fmt.Errorf("update process failed : %v", r))
			}
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
	// 백업 성공 표시
	backupCompleted = true

	// 다운로드 받을 파일명 가져오기
	ui.SetCurrentStep(3)
	ui.UpdateDetail("Getting update file name...")
	fileName, err := getFileNameFromServer(serverIP)
	if err != nil {
		LogError(err.Error())
		restore(ui)
		return
	}

	// 파일 다운로드
	ui.SetCurrentStep(4)
	ui.UpdateDetail("Downloading update...")
	downloadFilePath := fileName
	downloadURL := fmt.Sprintf("%s/update/file", serverIP)
	if err := downloadFile(downloadURL, downloadFilePath, ui); err != nil {
		LogError(err.Error())
		restore(ui)
		return
	}

	// Unzip the downloaded file
	ui.SetCurrentStep(5)
	ui.UpdateDetail("Extracting files...")
	if err := unzipFile(downloadFilePath, "."); err != nil {
		LogError(err.Error())
		restore(ui)
		return
	}

	if err := verifyHashSum(ui); err != nil {
		LogError(err.Error())
		restore(ui)
		return
	}

	ui.SetCurrentStep(6)
	ui.UpdateDetail("Cleaning up...")

	if err := removeFile(downloadFilePath, ui); err != nil {
		LogError(err.Error())
		restore(ui)
		return
	}

	// Launch the application
	if err := launchApplication(fromVersion, ui); err != nil {
		LogError(err.Error())
		restore(ui)
		return
	}

}
