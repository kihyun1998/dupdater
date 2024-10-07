package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("File Downloader")

	progress := widget.NewProgressBar()
	status := widget.NewLabel("Ready to start")

	button := widget.NewButton("Start Update", func() {
		go updateProcess(progress, status, myWindow)
	})

	content := container.NewVBox(
		progress,
		status,
		button,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(300, 100))
	myWindow.ShowAndRun()
}

func updateProcess(progress *widget.ProgressBar, status *widget.Label, window fyne.Window) {
	if err := moveFiles(status, window); err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error moving files: %v", err))
		})
		return
	}

	downloadFilePath := "downloaded_file.zip"
	if err := downloadFile("http://localhost:8000/update/file", downloadFilePath, progress, status, window); err != nil {
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

	// Launch the application
	launchApplication(status, window)
}

// 파일 이동 함수
func moveFiles(status *widget.Label, window fyne.Window) error {
	updateUI(window, func() {
		status.SetText("Moving files to backup directory...")
	})

	backupDir := "../BACK"
	if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating backup directory: %v", err)
	}
	files, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("error reading current directory: %v", err)
	}
	for _, file := range files {
		oldPath := file.Name()
		newPath := filepath.Join(backupDir, file.Name())
		if err := moveFile(oldPath, newPath); err != nil {
			return fmt.Errorf("error moving file %s: %v", file.Name(), err)
		}
		// }
	}
	return nil
}

// 다운로드 함수
func downloadFile(url, downloadFilePath string, progress *widget.ProgressBar, status *widget.Label, window fyne.Window) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error downloading file: %v", err))
		})
		return &exec.ExitError{}
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Bad status: %s", resp.Status))
		})
		return nil
	}

	// Create the file
	out, err := os.Create(downloadFilePath)
	if err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error creating file: %v", err))
		})
		return err
	}
	defer out.Close()

	// Create a custom io.Writer to track progress
	counter := &WriteCounter{
		Total:    resp.ContentLength,
		progress: progress,
		status:   status,
		window:   window,
	}

	// Start time for speed calculation
	startTime := time.Now()

	// Use io.Copy to optimize file writing and update progress
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		updateUI(window, func() {
			status.SetText(fmt.Sprintf("Error writing to file: %v", err))
		})
		return err
	}

	elapsedTime := time.Since(startTime).Seconds()
	speed := float64(counter.Written) / elapsedTime / 1024 // KB/s

	updateUI(window, func() {
		progress.SetValue(1)
		status.SetText(fmt.Sprintf("Download completed (%.2f KB/s)", speed))
	})

	time.Sleep(2 * time.Second)
	return nil
}

func unzipFile(zipFile, destDir string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		filePath := filepath.Join(destDir, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// 프로그램 실행 함수
func launchApplication(status *widget.Label, window fyne.Window) {
	updateUI(window, func() {
		status.SetText("Launching client.exe...")
	})

	cmd := exec.Command("./update_test_app_1.exe")
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

func moveFile(sourcePath, destPath string) error {
	// Try to move the file
	err := os.Rename(sourcePath, destPath)
	if err == nil {
		return nil
	}

	// If moving fails, try to copy and then delete
	err = copyFile(sourcePath, destPath)
	if err != nil {
		return err
	}

	// After successful copy, try to delete the original file
	return os.Remove(sourcePath)
}

func copyFile(sourcePath, destPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

type WriteCounter struct {
	Total      int64
	Written    int64
	progress   *widget.ProgressBar
	status     *widget.Label
	window     fyne.Window
	lastUpdate time.Time
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Written += int64(n)
	percentage := float64(wc.Written) / float64(wc.Total)

	// Update UI every 100ms to reduce update frequency
	if time.Since(wc.lastUpdate) > 100*time.Millisecond {
		updateUI(wc.window, func() {
			wc.progress.SetValue(percentage)
			speed := float64(wc.Written) / time.Since(wc.lastUpdate).Seconds() / 1024 // KB/s
			wc.status.SetText(fmt.Sprintf("Downloading... %.2f%% (%.2f KB/s)", percentage*100, speed))
		})
		wc.lastUpdate = time.Now()
	}

	return n, nil
}

func updateUI(window fyne.Window, f func()) {
	window.Canvas().Refresh(window.Content())
	f()
}
