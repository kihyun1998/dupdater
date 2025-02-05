// Package infrastructure는 파일 도메인의 실제 구현체들을 제공합니다
package infrastructure

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/kihyun1998/dupdater/internal/file/domain/entity"
	"github.com/kihyun1998/dupdater/internal/file/domain/repository"
)

// fileManager는 파일 시스템 작업을 구현하는 구조체입니다
type fileManager struct {
	mu     sync.RWMutex
	logger Logger
	config *repository.Config
}

// NewFileManager는 새로운 fileManager 인스턴스를 생성합니다
func NewFileManager(config *repository.Config, logger Logger) repository.FileRepository {
	return &fileManager{
		config: config,
		logger: logger,
	}
}

// CreateDirectory는 디렉토리를 생성합니다
func (m *fileManager) CreateDirectory(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// DeleteFile은 파일을 삭제합니다
func (m *fileManager) DeleteFile(path string) error {
	return os.Remove(path)
}

// CopyFile은 파일을 복사합니다
func (m *fileManager) CopyFile(src, dest string) error {
	// Rename 시도
	if err := os.Rename(src, dest); err == nil {
		return nil
	}

	// Rename 실패시 복사 후 삭제
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

// MoveFile은 파일을 이동합니다
func (m *fileManager) MoveFile(src, dest string) error {
	if err := os.Rename(src, dest); err != nil {
		return m.CopyFile(src, dest)
	}
	return nil
}

// GetFileInfo는 파일 정보를 가져옵니다
func (m *fileManager) GetFileInfo(path string) (*entity.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return entity.NewFileInfo(path, info), nil
}

// ListFiles는 디렉토리 내의 모든 파일 목록을 반환합니다
func (m *fileManager) ListFiles(dir string) ([]*entity.FileInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	files := make([]*entity.FileInfo, 0, len(entries))
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			m.logger.Error("파일 정보 가져오기 실패 (%s): %v", path, err)
			continue
		}
		files = append(files, entity.NewFileInfo(path, info))
	}

	return files, nil
}

// ExtractZip은 ZIP 파일을 압축 해제합니다
func (m *fileManager) ExtractZip(zipFile, destDir string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("ZIP 파일 열기 실패: %w", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		if err := m.extractFile(file, destDir); err != nil {
			return fmt.Errorf("파일 압축해제 실패 (%s): %w", file.Name, err)
		}
	}

	return nil
}

// extractFile은 단일 파일을 압축 해제합니다
func (m *fileManager) extractFile(file *zip.File, destDir string) error {
	filePath := filepath.Join(destDir, file.Name)

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

// Logger는 로깅을 위한 인터페이스입니다
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}
