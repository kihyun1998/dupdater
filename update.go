package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func getServerIP(serverName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting user home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".testfolder", "config.json")
	file, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("error reading config file: %w", err)
	}

	var data struct {
		ProfileList []struct {
			Name string `json:"name"`
			IP   string `json:"ip"`
		} `json:"profileList"`
	}

	if err := json.Unmarshal(file, &data); err != nil {
		return "", fmt.Errorf("error parsing config file: %w", err)
	}

	for _, profile := range data.ProfileList {
		if profile.Name == serverName {
			return profile.IP, nil
		}
	}

	return "", fmt.Errorf("server with name '%s' not found", serverName)
}

func getFileNameFromServer(serverIP string) (string, error) {
	url := fmt.Sprintf("%s/update/updatefilename", serverIP)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error failed to get file name from server :%v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Filename string `json:"filename"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error json decode :%v", err)
	}
	return result.Filename, nil
}

func updateProcess(serverName, fromVersion string, progress *widget.ProgressBar, status *widget.Label, window fyne.Window) {
	updateUI(window, func() {
		status.SetText("Starting update process...")
	})

	serverIP, err := getServerIP(serverName)
	if err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error getting server IP: %v", err))
		})
		return
	}

	waitForApplicationToClose(applicationName, status, window)

	if err := moveFiles(status, window); err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error moving files: %v", err))
		})
		return
	}

	fileName, err := getFileNameFromServer(serverIP)
	if err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("error get file name: %v", err))
		})
		return
	}

	// downloadFilePath := "downloaded_file.zip"
	downloadURL := fmt.Sprintf("%s/update/file", serverIP)
	if err := downloadFile(downloadURL, fileName, progress, status, window); err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error downloading file: %v", err))
		})
		return
	}

	// Unzip the downloaded file
	updateUI(window, func() {
		status.SetText("Extracting files...")
	})
	if err := unzipFile(fileName, "."); err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error extracting files: %v", err))
		})
		return
	}

	updateUI(window, func() {
		status.SetText("Cleaning up...")
	})

	if err := removeFile(fileName); err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Warning: Failed to delete file: %v", err))
		})
	}

	// Launch the application
	launchApplication(fromVersion, status, window)
}
