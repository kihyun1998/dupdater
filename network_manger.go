package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// 프로파일명을 통해 파일을 읽고 서버 주소를 가져오는 함수
func getServerIP(serverName string) (serverIP string, err error) {
	defer func() {
		if r := recover(); r != nil {
			LogError("Panic in getServerIP: %v", r)
			err = fmt.Errorf("getServerIP failed unexpectedly: %v", r)
			// ui.UpdateDetail("Error occurred. Attempting to restore files...")
			// if err := restoreFiles(ui); err != nil {
			// 	LogError("Failed to restore files: %v", err)
			// 	ui.ShowError(fmt.Errorf("update and restore failed: %v", r))
			// } else {
			// 	ui.ShowError(fmt.Errorf("update failed, files restored: %v", r))
			// }
		}
	}()

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

// 서버 아이피를 통해 GET 요청을 통해 다운로드 받을 파일명을 가져오는 함수
func getFileNameFromServer(serverIP string) (fileName string, err error) {
	defer func() {
		if r := recover(); r != nil {
			LogError("Panic in getFileNameFromServer: %v", r)
			err = fmt.Errorf("getFileNameFromServer failed unexpectedly: %v", r)
		}
	}()

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
