package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func verifyFileHash(filePath string, status *widget.Label, window fyne.Window) error {
	updateUI(window, func() {
		status.SetText("Verifying file Hash...")
	})

	file, err := os.Open(filePath)
	if err != nil {
		err = fmt.Errorf("error opening file: %v", err)
		<-showErrorDialog(err, window)
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		err = fmt.Errorf("error getting file info: %v", err)
		<-showErrorDialog(err, window)
		return err
	}

	hashData := make([]byte, 36)
	_, err = file.ReadAt(hashData, fileInfo.Size()-36)
	if err != nil {
		err = fmt.Errorf("error reading hash data: %v", err)
		<-showErrorDialog(err, window)
		return err

	}

	hashLength := binary.LittleEndian.Uint32(hashData[:4])
	storedHash := hashData[4:]

	hash := sha256.New()
	if _, err := io.CopyN(hash, file, fileInfo.Size()-36); err != nil {
		err = fmt.Errorf("error calculating file hash: %v", err)
		<-showErrorDialog(err, window)
		return err

	}

	calculateHash := hash.Sum(nil)

	if !compareHashes(calculateHash, storedHash[:hashLength]) {
		err := fmt.Errorf("file hash verification faield")
		<-showErrorDialog(err, window)
		return err
	}

	updateUI(window, func() {
		status.SetText("File hash verified successfully")
	})

	return nil
}

func verifyHashSum(status *widget.Label, window fyne.Window) error {
	updateUI(window, func() {
		status.SetText("Verifying hash.sum file...")
	})

	sumFilePath := filepath.Join(".", "hash_sum.txt")
	file, err := os.Open(sumFilePath)
	if err != nil {
		err = fmt.Errorf("error opening hash.sum file: %v", err)
		<-showErrorDialog(err, window)
		return err

	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")
		if len(parts) != 3 {
			err = fmt.Errorf("invalid line in hash.sum file: %s", line)
			<-showErrorDialog(err, window)
			return err

		}
		pathHash, dataHash := parts[1], parts[2]

		filePath := getFilePathFromHash(pathHash)
		if filePath == "" {
			err = fmt.Errorf("file not found for hash: %s", pathHash)
			<-showErrorDialog(err, window)
			return err

		}

		if err := verifyFileHashSum(filePath, dataHash); err != nil {
			err = fmt.Errorf("hash verification failed for %s:%v", filePath, err)
			<-showErrorDialog(err, window)
			return err

		}
	}

	if err := scanner.Err(); err != nil {
		err = fmt.Errorf("error reading hash.sum file: %v", err)
		<-showErrorDialog(err, window)
		return err

	}

	updateUI(window, func() {
		status.SetText("hash.sum verified successfully")
	})

	return nil
}

func verifyFileHashSum(filePath, expectedHash string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	calculatedHash := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	if calculatedHash != expectedHash {
		return fmt.Errorf("hash mismatch")
	}
	return nil
}

func getFilePathFromHash(pathHash string) string {
	files, err := filepath.Glob("*")
	if err != nil {
		return ""
	}

	for _, file := range files {
		hash := calculatePathHash(file)
		if hash == pathHash {
			return file
		}
	}
	return ""
}

func calculatePathHash(relPath string) string {
	hash := sha256.Sum256([]byte(relPath))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func compareHashes(hash1, hash2 []byte) bool {
	if len(hash1) != len(hash2) {
		return false
	}
	for i := range hash1 {
		if hash1[i] != hash2[i] {
			return false
		}
	}
	return true
}
