// Package infrastructure는 파일 시스템 작업의 실제 구현체를 제공합니다
package infrastructure

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/kihyun1998/dupdater/internal/file/domain/repository"
)

// LocalFileRepository는 로컬 파일 시스템 기반의 저장소 구현체입니다
type LocalFileRepository struct {
	config *repository.Config
	mu     sync.RWMutex
}

// NewLocalFileRepository는 새로운 LocalFileRepository 인스턴스를 생성합니다
func NewLocalFileRepository(config *repository.Config) repository.FileRepository {
	return &LocalFileRepository{
		config: config,
	}
}

// cleanCurrentDirectory는 현재 디렉토리를 정리합니다
func (r *LocalFileRepository) cleanCurrentDirectory() error {
	files, err := os.ReadDir(r.config.DirInfo.CurrentDir)
	if err != nil {
		return fmt.Errorf("디렉토리 읽기 실패: %w", err)
	}

	for _, file := range files {
		path := filepath.Join(r.config.DirInfo.CurrentDir, file.Name())

		// 여러 번 삭제 시도
		for i := 0; i < 3; i++ {
			var err error
			if file.IsDir() {
				r.config.Logger.Info("디렉토리 삭제 시도: %s", path)
				err = os.RemoveAll(path)
			} else {
				r.config.Logger.Info("파일 삭제 시도: %s", path)
				err = os.Remove(path)
			}

			if err == nil {
				break
			}

			// 마지막 시도에서는 .old를 붙여서 시도
			if i == 2 {
				newPath := path + ".old"
				r.config.Logger.Error("삭제 실패, 이름 변경 후 재시도: %s -> %s", path, newPath)
				if os.Rename(path, newPath) == nil {
					os.Remove(newPath)
				}
			}

			time.Sleep(time.Second)
		}
	}
	return nil
}

// DeleteFile은 지정된 파일을 삭제합니다
func (r *LocalFileRepository) DeleteFile(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("파일 삭제 실패: %w", err)
	}
	return nil
}

// Extract ----------------------------------------------------------
// ExtractZip은 ZIP 파일을 압축 해제합니다
func (r *LocalFileRepository) ExtractZip(zipFile string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("ZIP 파일 열기 실패: %w", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		if err := r.extractFile(file); err != nil {
			return fmt.Errorf("파일 압축해제 실패 (%s): %w", file.Name, err)
		}
	}

	return nil
}

// extractFile은 ZIP 파일의 항목을 압축 해제합니다
func (r *LocalFileRepository) extractFile(file *zip.File) error {
	filePath := filepath.Join(r.config.DirInfo.CurrentDir, file.Name)

	if file.FileInfo().IsDir() {
		return os.MkdirAll(filePath, os.ModePerm)
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return fmt.Errorf("디렉토리 생성 실패: %w", err)
	}

	dest, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return fmt.Errorf("대상 파일 생성 실패: %w", err)
	}
	defer dest.Close()

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("ZIP 파일 열기 실패: %w", err)
	}
	defer src.Close()

	if _, err := io.Copy(dest, src); err != nil {
		return fmt.Errorf("파일 복사 실패: %w", err)
	}

	return nil
}

// ------------------------------------------------------------------

// Restore ----------------------------------------------------------
// Restore는 백업된 파일들을 복원합니다
func (r *LocalFileRepository) Restore() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 백업 디렉토리 존재 확인
	if _, err := os.Stat(r.config.DirInfo.BackupDir); os.IsNotExist(err) {
		return fmt.Errorf("백업 디렉토리가 존재하지 않습니다: %s", r.config.DirInfo.BackupDir)
	}

	// 현재 디렉토리 정리
	if err := r.cleanCurrentDirectory(); err != nil {
		return fmt.Errorf("현재 디렉토리 정리 실패: %w", err)
	}

	// 백업 파일 복원
	files, err := os.ReadDir(r.config.DirInfo.BackupDir)
	if err != nil {
		return fmt.Errorf("백업 디렉토리 읽기 실패: %w", err)
	}

	for _, file := range files {
		srcPath := filepath.Join(r.config.DirInfo.BackupDir, file.Name())
		destPath := filepath.Join(r.config.DirInfo.CurrentDir, file.Name())

		if file.IsDir() {
			if err := r.restoreDirectory(srcPath, destPath); err != nil {
				return fmt.Errorf("디렉토리 복원 실패 (%s): %w", file.Name(), err)
			}
		} else {
			if err := r.restoreFile(srcPath, destPath); err != nil {
				return fmt.Errorf("파일 복원 실패 (%s): %w", file.Name(), err)
			}
		}
	}

	return nil
}

// restoreDirectory는 디렉토리를 복구합니다.
func (r *LocalFileRepository) restoreDirectory(src, dest string) error {
	// 디렉토리 존재하면 생성
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return fmt.Errorf("디렉토리 생성 실패: %w", err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("디렉토리 읽기 실패: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			if err := r.restoreDirectory(srcPath, destPath); err != nil {
				return err
			}
		} else {
			if err := r.restoreFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}
	r.config.Logger.Info("디렉토리 복원 완료: %s", dest)
	return os.RemoveAll(src)
}

// restoreFile은 파일을 복원합니다
func (r *LocalFileRepository) restoreFile(src, dest string) error {
	// 기존 파일이 존재하면 삭제
	if _, err := os.Stat(dest); err == nil {
		if err := os.Remove(dest); err != nil {
			r.config.Logger.Info("기존 파일 삭제 시도: %s", dest)
			if err := os.Remove(dest); err != nil {
				return fmt.Errorf("기존 파일 삭제 실패: %w", err)
			}
		}
	}

	// Rename 시도
	if err := os.Rename(src, dest); err == nil {
		r.config.Logger.Info("파일 복원 완료(Rename): %s -> %s", src, dest)
		return nil
	}

	// Rename 실패시 복사
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("소스 파일 열기 실패: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("대상 파일 생성 실패: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("파일 복사 실패: %w", err)
	}

	r.config.Logger.Info("파일 복원 완료(Copy): %s -> %s", src, dest)
	return os.Remove(src)
}

// ------------------------------------------------------------------

// Backup------------------------------------------------------------
// Backup은 현재 디렉토리의 파일들을 백업합니다
func (r *LocalFileRepository) Backup() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 백업 디렉토리 생성
	if err := os.MkdirAll(r.config.DirInfo.BackupDir, os.ModePerm); err != nil {
		return fmt.Errorf("백업 디렉토리 생성 실패: %w", err)
	}

	// 현재 디렉토리 파일 목록 조회
	files, err := os.ReadDir(r.config.DirInfo.CurrentDir)
	if err != nil {
		return fmt.Errorf("디렉토리 읽기 실패: %w", err)
	}

	// 각 파일 백업
	for _, file := range files {
		oldPath := filepath.Join(r.config.DirInfo.CurrentDir, file.Name())
		newPath := filepath.Join(r.config.DirInfo.BackupDir, file.Name())

		if file.IsDir() {
			if err := r.backupDirectory(oldPath, newPath); err != nil {
				return fmt.Errorf("디렉토리 백업 실패 (%s): %w", file.Name(), err)
			}
		} else {
			if err := r.backupFile(oldPath, newPath); err != nil {
				return fmt.Errorf("파일 백업 실패 (%s): %w", file.Name(), err)
			}
		}
	}

	return nil
}

// GetBackupDir은 백업 디렉토리 경로를 반환합니다
func (r *LocalFileRepository) GetBackupDir() string {
	return r.config.DirInfo.BackupDir
}

// backupDirectory는 디렉토리를 백업합니다
func (r *LocalFileRepository) backupDirectory(src, dest string) error {
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return fmt.Errorf("디렉토리 생성 실패: %w", err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("디렉토리 읽기 실패: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			if err := r.backupDirectory(srcPath, destPath); err != nil {
				return err
			}
		} else {
			if err := r.backupFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return os.RemoveAll(src)
}

// backupFile은 단일 파일을 백업합니다
func (r *LocalFileRepository) backupFile(src, dest string) error {
	// 먼저 Rename 시도
	if err := os.Rename(src, dest); err == nil {
		return nil
	}

	// Rename 실패시 복사 후 삭제
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("소스 파일 열기 실패: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("대상 파일 생성 실패: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("파일 복사 실패: %w", err)
	}

	// 원본 파일 삭제 시도
	for i := 0; i < 3; i++ {
		if err := os.Remove(src); err == nil {
			return nil
		}
		time.Sleep(time.Second)
	}

	return os.Remove(src)
}

// ------------------------------------------------------------------
