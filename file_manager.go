package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// 백업폴더로 이동함수
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

// 압축해제 함수
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

// 파일 이동 함수
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

// 파일 복사 함수
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
