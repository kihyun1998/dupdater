package infrastructure

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	loggerPorts "github.com/kihyun1998/dupdater/internal/logger/domain/ports"
)

type FileManager struct {
	logger     loggerPorts.LoggerPort
	backupDir  string
	currentDir string
}

func NewFileManager(logger loggerPorts.LoggerPort, backupDir, currentDir string) *FileManager {
	return &FileManager{
		logger:     logger,
		backupDir:  backupDir,
		currentDir: currentDir,
	}
}

func (m *FileManager) Backup() error {
	// 패닉 복구
	defer func() {
		if r := recover(); r != nil {
			m.logger.Error("Backup 함수에서 패닉 발생: %v", r)
		}
	}()

	// 백업 디렉토리 생성
	if err := os.MkdirAll(m.backupDir, os.ModePerm); err != nil {
		return fmt.Errorf("백업 디렉토리 생성 실패: %w", err)
	}

	// 현재 디렉토리 파일 목록 조회
	files, err := os.ReadDir(m.currentDir)
	if err != nil {
		return fmt.Errorf("디렉토리 읽기 실패: %w", err)
	}

	// 각 파일/디렉토리 백업
	for _, file := range files {
		oldPath := filepath.Join(m.currentDir, file.Name())
		newPath := filepath.Join(m.backupDir, file.Name())

		if file.IsDir() {
			if err := m.backupDirectory(oldPath, newPath); err != nil {
				return fmt.Errorf("디렉토리 백업 실패 (%s): %w", file.Name(), err)
			}
		} else {
			if err := m.backupFile(oldPath, newPath); err != nil {
				return fmt.Errorf("파일 백업 실패 (%s): %w", file.Name(), err)
			}
		}
	}

	m.logger.Info("백업 완료: %s", m.backupDir)
	return nil
}

func (m *FileManager) Restore() error {
	defer func() {
		if r := recover(); r != nil {
			m.logger.Error("Restore 함수에서 패닉 발생: %v", r)
		}
	}()

	// 백업 디렉토리 존재 확인
	if _, err := os.Stat(m.backupDir); os.IsNotExist(err) {
		return fmt.Errorf("백업 디렉토리가 존재하지 않습니다: %s", m.backupDir)
	}

	// 현재 디렉토리 정리
	if err := m.cleanCurrentDirectory(); err != nil {
		return fmt.Errorf("현재 디렉토리 정리 실패: %w", err)
	}

	// 백업 파일 복원
	files, err := os.ReadDir(m.backupDir)
	if err != nil {
		return fmt.Errorf("백업 디렉토리 읽기 실패: %w", err)
	}

	for _, file := range files {
		srcPath := filepath.Join(m.backupDir, file.Name())
		destPath := filepath.Join(m.currentDir, file.Name())

		if file.IsDir() {
			if err := m.restoreDirectory(srcPath, destPath); err != nil {
				return fmt.Errorf("디렉토리 복원 실패 (%s): %w", file.Name(), err)
			}
		} else {
			if err := m.restoreFile(srcPath, destPath); err != nil {
				return fmt.Errorf("파일 복원 실패 (%s): %w", file.Name(), err)
			}
		}
	}

	// 백업 디렉토리 정리
	if err := os.RemoveAll(m.backupDir); err != nil {
		m.logger.Error("백업 디렉토리 삭제 실패: %v", err)
	}

	m.logger.Info("복원 완료")
	return nil
}

func (m *FileManager) ExtractZip(zipFile string) error {
	defer func() {
		if r := recover(); r != nil {
			m.logger.Error("ExtractZip 함수에서 패닉 발생: %v", r)
		}
	}()

	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("ZIP 파일 열기 실패: %w", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		if err := m.extractFile(file); err != nil {
			return fmt.Errorf("파일 압축해제 실패 (%s): %w", file.Name, err)
		}
	}

	return nil
}

func (m *FileManager) DeleteFile(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("파일 삭제 실패: %w", err)
	}
	return nil
}

// 내부 헬퍼 함수들 - 기존 코드 그대로 유지
func (m *FileManager) backupFile(src, dest string) error {
	// 먼저 Rename 시도
	if err := os.Rename(src, dest); err == nil {
		return nil
	}

	// Rename 실패시 복사 후 삭제 시도
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// 여러 번 삭제 시도
	for i := 0; i < 3; i++ {
		if err := os.Remove(src); err == nil {
			return nil
		}
		time.Sleep(time.Second)
	}

	return os.Remove(src)
}

func (m *FileManager) backupDirectory(src, dest string) error {
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			if err := m.backupDirectory(srcPath, destPath); err != nil {
				return err
			}
		} else {
			if err := m.backupFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return os.RemoveAll(src)
}

func (m *FileManager) restoreFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	return os.Remove(src)
}

func (m *FileManager) restoreDirectory(src, dest string) error {
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			if err := m.restoreDirectory(srcPath, destPath); err != nil {
				return err
			}
		} else {
			if err := m.restoreFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return os.RemoveAll(src)
}

func (m *FileManager) cleanCurrentDirectory() error {
	files, err := os.ReadDir(m.currentDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		path := filepath.Join(m.currentDir, file.Name())
		for i := 0; i < 3; i++ {
			var err error
			if file.IsDir() {
				err = os.RemoveAll(path)
			} else {
				err = os.Remove(path)
			}

			if err == nil {
				break
			}

			if i == 2 {
				// 마지막 시도에서는 이름 변경 후 삭제 시도
				newPath := path + ".old"
				if os.Rename(path, newPath) == nil {
					os.Remove(newPath)
				}
			}

			time.Sleep(time.Second)
		}
	}
	return nil
}

func (m *FileManager) extractFile(file *zip.File) error {
	filePath := filepath.Join(m.currentDir, file.Name)

	if file.FileInfo().IsDir() {
		return os.MkdirAll(filePath, os.ModePerm)
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	dest, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer dest.Close()

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = io.Copy(dest, src)
	return err
}
