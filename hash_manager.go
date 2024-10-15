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
	"time"
)

func verifyFileHash(filePath string, ui *UpdaterUI) (err error) {
	defer func() {
		if r := recover(); r != nil {
			LogError("Panic in verifyFileHash: %v", r)
			err = fmt.Errorf("verifyFileHash failed unexpectedly: %v", r)
		}
	}()

	ui.UpdateDetail("Verifying file Hash...")

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}

	hashData := make([]byte, 36)
	_, err = file.ReadAt(hashData, fileInfo.Size()-36)
	if err != nil {
		return fmt.Errorf("error reading hash data: %v", err)
	}

	hashLength := binary.LittleEndian.Uint32(hashData[:4])
	storedHash := hashData[4:]

	hash := sha256.New()
	if _, err := io.CopyN(hash, file, fileInfo.Size()-36); err != nil {
		return fmt.Errorf("error calculating file hash: %v", err)
	}

	calculateHash := hash.Sum(nil)

	result, err := compareHashes(calculateHash, storedHash[:hashLength])
	if err != nil {
		return fmt.Errorf("file hash verification faield")
	}

	if !result {
		return fmt.Errorf("file hash verification faield")
	}
	ui.UpdateDetail("File hash verified successfully")

	time.Sleep(time.Second)
	return nil
}

func verifyHashSum(ui *UpdaterUI) (err error) {
	defer func() {
		if r := recover(); r != nil {
			LogError("Panic in verifyHashSum: %v", r)
			// err = fmt.Errorf("verifyHashSum failed unexpectedly: %v", r)
			ui.UpdateDetail("Error occurred. Attempting to restore files...")
			if err := restoreFiles(ui); err != nil {
				LogError("Failed to restore files: %v", err)
				ui.ShowError(fmt.Errorf("update and restore failed: %v", r))
			} else {
				ui.ShowError(fmt.Errorf("update failed, files restored: %v", r))
			}
		}
	}()

	ui.UpdateDetail("Verifying hash.sum file...")

	sumFilePath := filepath.Join(".", "hash_sum.txt")
	file, err := os.Open(sumFilePath)
	if err != nil {
		return fmt.Errorf("error opening hash.sum file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")
		if len(parts) != 3 {
			return fmt.Errorf("invalid line in hash.sum file: %s", line)
		}
		pathHash, dataHash := parts[1], parts[2]

		filePath, err := getFilePathFromHash(pathHash)
		if filePath == "" {
			return fmt.Errorf("file not found for hash: %s", pathHash)
		}
		if err != nil {
			return fmt.Errorf("getFilePathFromHash error: %s", pathHash)
		}

		if err := verifyFileHashSum(filePath, dataHash); err != nil {
			return fmt.Errorf("hash verification failed for %s:%v", filePath, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading hash.sum file: %v", err)
	}

	ui.UpdateDetail("hash.sum verified successfully")

	return nil
}

func verifyFileHashSum(filePath, expectedHashBase64 string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			LogError("Panic in verifyFileHashSum: %v", r)
			err = fmt.Errorf("verifyFileHashSum failed unexpectedly: %v", r)
		}
	}()

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	calculatedHash := hash.Sum(nil)
	calculatedHashBase64 := base64.StdEncoding.EncodeToString(calculatedHash)

	if calculatedHashBase64 != expectedHashBase64 {
		return fmt.Errorf("hash mismatch for %s", filePath)
	}
	return nil
}

func getFilePathFromHash(pathHash string) (result string, err error) {
	defer func() {
		if r := recover(); r != nil {
			LogError("Panic in getFilePathFromHash: %v", r)
			err = fmt.Errorf("getFilePathFromHash failed unexpectedly: %v", r)
		}
	}()

	var matchedPath string
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(".", path)
		if err != nil {
			return err
		}

		relPath = filepath.ToSlash(relPath)

		if info.IsDir() {
			return nil
		}

		hash, err := calculatePathHash(relPath)
		if err != nil {
			return fmt.Errorf("error calculate path hash: %v", err)
		}
		LogInfo("pathHash: %s, calculated hash: %s for path: %s", pathHash, hash, relPath)

		if hash == pathHash {
			matchedPath = relPath
			return filepath.SkipDir // 파일을 찾았으므로 검색 중단
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("error walking directory: %v", err)
	}

	if matchedPath == "" {
		return "", fmt.Errorf("no matching file found for hash: %s", pathHash)
	}

	return matchedPath, nil
}

func calculatePathHash(relPath string) (result string, err error) {
	defer func() {
		if r := recover(); r != nil {
			LogError("Panic in calculatePathHash: %v", r)
			err = fmt.Errorf("calculatePathHash failed unexpectedly: %v", r)
		}
	}()

	hash := sha256.Sum256([]byte(relPath))
	base64Hash := base64.StdEncoding.EncodeToString(hash[:])
	return base64Hash, nil
}

func compareHashes(hash1, hash2 []byte) (result bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			LogError("Panic in compareHashes: %v", r)
			err = fmt.Errorf("compareHashes failed unexpectedly: %v", r)
		}
	}()

	if len(hash1) != len(hash2) {
		return false, nil
	}
	for i := range hash1 {
		if hash1[i] != hash2[i] {
			return false, nil
		}
	}
	return true, nil
}
