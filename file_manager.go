package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// 백업폴더로 이동함수
func moveFiles(ui *UpdaterUI) (err error) {
	defer func() {
		if r := recover(); r != nil {
			LogError("Panic in moveFiles: %v", r)
			err = fmt.Errorf("moveFiles failed unexpectedly: %v", r)
		}
	}()

	ui.UpdateDetail("Moving files to backup directory...")

	backupDir := filepath.Join(os.TempDir(), "ACRABACK")
	ui.UpdateDetail(fmt.Sprintf("back dir is %s", backupDir))

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

		if file.IsDir() {
			// 디렉토리인 경우
			if err := os.MkdirAll(newPath, os.ModePerm); err != nil {
				return fmt.Errorf("error creating directory %s: %v", file.Name(), err)
			}
			// 디렉토리 내용을 재귀적으로 복사
			if err := copyDir(oldPath, newPath); err != nil {
				return fmt.Errorf("error copying directory %s: %v", file.Name(), err)
			}
			// 원본 디렉토리 삭제
			if err := os.RemoveAll(oldPath); err != nil {
				return fmt.Errorf("error removing original directory %s: %v", file.Name(), err)
			}
		} else {
			// 파일인 경우
			if err := moveFile(oldPath, newPath); err != nil {
				return fmt.Errorf("error moving file %s: %v", file.Name(), err)
			}
		}
	}
	return nil
}

// 압축해제 함수
func unzipFile(zipFile, destDir string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			LogError("Panic in unzipFile: %v", r)
			err = fmt.Errorf("error unzip: unzipFile failed unexpectedly: %v", r)
		}
	}()

	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("error unzip: %v", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		filePath := filepath.Join(destDir, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		// 파일을 위한 디렉토리 생성
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return fmt.Errorf("error unzip: failed to create directory for %s: %v", filePath, err)
		}

		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("error unzip: failed to create %s: %v", filePath, err)
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("error unzip: %v", err)
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return fmt.Errorf("error unzip: %v", err)
		}
	}
	return nil
}

// 파일 이동 함수
func moveFile(sourcePath, destPath string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			LogError("Panic in moveFile: %v", r)
			err = fmt.Errorf("moveFile failed unexpectedly: %v", r)
		}
	}()

	// Try to move the file
	err = os.Rename(sourcePath, destPath)
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
func copyFile(sourcePath, destPath string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			LogError("Panic in copyFile: %v", r)
			err = fmt.Errorf("copyFile failed unexpectedly: %v", r)
		}
	}()

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

// 디렉토리를 복사하는 함수
func copyDir(src string, dst string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			LogError("Panic in copyDir: %v", r)
			err = fmt.Errorf("copyDir failed unexpectedly: %v", r)
		}
	}()

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, os.ModePerm); err != nil {
				return err
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func removeFile(downloadFilePath string, ui *UpdaterUI) (err error) {
	defer func() {
		if r := recover(); r != nil {
			LogError("Panic in removeFile: %v", r)
			err = fmt.Errorf("removeFile failed unexpectedly: %v", r)
		}
	}()

	if err := os.Remove(downloadFilePath); err != nil {
		return err
	}
	return nil
}

func restoreFiles(ui *UpdaterUI) (err error) {
	defer func() {
		if r := recover(); r != nil {
			LogError("Panic in restoreFiles: %v", r)
			err = fmt.Errorf("restoreFiles failed unexpectedly: %v", r)
		}
	}()

	ui.UpdateDetail("Restoring files from backup directory...")

	backupDir := filepath.Join(os.TempDir(), "ACRABACK")
	ui.UpdateDetail(fmt.Sprintf("Backup directory is %s", backupDir))

	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return fmt.Errorf("backup directory does not exist: %v", err)
	}

	// 현재 디렉토리의 내용을 모두 제거
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %v", err)
	}

	ui.UpdateDetail("Removing existing files and directories...")
	err = removeDirContents(currentDir)
	if err != nil {
		return fmt.Errorf("error removing current directory contents: %v", err)
	}

	files, err := os.ReadDir(backupDir)
	if err != nil {
		return fmt.Errorf("error reading backup directory: %v", err)
	}

	for _, file := range files {
		oldPath := filepath.Join(backupDir, file.Name())
		newPath := filepath.Join(".", file.Name())

		if file.IsDir() {
			// 디렉토리인 경우
			if err := os.MkdirAll(newPath, os.ModePerm); err != nil {
				return fmt.Errorf("error creating directory %s: %v", file.Name(), err)
			}
			// 디렉토리 내용을 재귀적으로 복사
			if err := copyDir(oldPath, newPath); err != nil {
				return fmt.Errorf("error copying directory %s: %v", file.Name(), err)
			}
			// 백업 디렉토리에서 해당 디렉토리 삭제
			if err := os.RemoveAll(oldPath); err != nil {
				return fmt.Errorf("error removing backup directory %s: %v", file.Name(), err)
			}
		} else {
			// 파일인 경우
			if err := moveFile(oldPath, newPath); err != nil {
				return fmt.Errorf("error moving file %s: %v", file.Name(), err)
			}
		}
	}

	// 백업 디렉토리 삭제
	if err := os.RemoveAll(backupDir); err != nil {
		return fmt.Errorf("error removing backup directory: %v", err)
	}

	return nil
}

// 디렉토리의 내용을 모두 제거하는 함수
func removeDirContents(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			if err := os.RemoveAll(path); err != nil {
				return err
			}
		} else {
			for retries := 0; retries < 3; retries++ {
				err := os.Remove(path)
				if err == nil {
					break
				}
				if os.IsPermission(err) {
					return err
				}
				if retries == 2 {
					newPath := path + ".old"
					if err := os.Rename(path, newPath); err == nil {
						if err := os.Remove(newPath); err != nil {
							return fmt.Errorf("failed to remove renamed file %s: %v", newPath, err)
						}
					} else {
						return fmt.Errorf("failed to rename and remove file %s: %v", path, err)
					}
				}
				time.Sleep(time.Second)
			}

		}
	}
	return nil
}
