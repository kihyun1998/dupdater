package fileutils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// FileManager는 파일 작업을 관리하는 구조체
type FileManager struct {
	logger Logger
}

// Logger는 로깅을 위한 인터페이스
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}

// NewFileManager는 새로운 FileManager 인스턴스를 생성
func NewFileManager(logger Logger) *FileManager {
	return &FileManager{
		logger: logger,
	}
}

// BackupFiles는 지정된 디렉토리의 파일들을 백업
func (fm *FileManager) BackupFiles(srcDir, backupDir string) error {
	fm.logger.Info("백업 시작: %s -> %s", srcDir, backupDir)

	// 백업 디렉토리 생성
	if err := os.MkdirAll(backupDir, os.ModePerm); err != nil {
		return fmt.Errorf("백업 디렉토리 생성 실패: %w", err)
	}

	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 상대 경로 계산
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return fmt.Errorf("상대 경로 계산 실패: %w", err)
		}

		// 대상 경로 생성
		dstPath := filepath.Join(backupDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// 파일 복사
		if err := fm.CopyFile(path, dstPath); err != nil {
			return fmt.Errorf("파일 복사 실패 (%s): %w", relPath, err)
		}

		fm.logger.Info("파일 백업 완료: %s", relPath)
		return nil
	})
}

// RestoreFiles는 백업된 파일들을 복원
func (fm *FileManager) RestoreFiles(backupDir, targetDir string) error {
	fm.logger.Info("복원 시작: %s -> %s", backupDir, targetDir)

	// 대상 디렉토리의 기존 내용 삭제
	if err := fm.CleanDirectory(targetDir); err != nil {
		return fmt.Errorf("대상 디렉토리 정리 실패: %w", err)
	}

	// 백업에서 복원
	return filepath.Walk(backupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 상대 경로 계산
		relPath, err := filepath.Rel(backupDir, path)
		if err != nil {
			return fmt.Errorf("상대 경로 계산 실패: %w", err)
		}

		// 대상 경로 생성
		dstPath := filepath.Join(targetDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// 파일 복사
		if err := fm.CopyFile(path, dstPath); err != nil {
			return fmt.Errorf("파일 복사 실패 (%s): %w", relPath, err)
		}

		fm.logger.Info("파일 복원 완료: %s", relPath)
		return nil
	})
}

// CopyFile은 단일 파일을 복사
func (fm *FileManager) CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("소스 파일 열기 실패: %w", err)
	}
	defer srcFile.Close()

	// 파일 정보 획득
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("파일 정보 획득 실패: %w", err)
	}

	// 대상 디렉토리 생성
	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return fmt.Errorf("대상 디렉토리 생성 실패: %w", err)
	}

	// 대상 파일 생성
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("대상 파일 생성 실패: %w", err)
	}
	defer dstFile.Close()

	// 파일 복사
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("파일 복사 실패: %w", err)
	}

	return nil
}

// UnzipFile은 ZIP 파일의 압축을 해제
func (fm *FileManager) UnzipFile(zipPath, destDir string) error {
	fm.logger.Info("ZIP 파일 압축 해제 시작: %s -> %s", zipPath, destDir)

	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("ZIP 파일 열기 실패: %w", err)
	}
	defer reader.Close()

	// 대상 디렉토리 생성
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return fmt.Errorf("대상 디렉토리 생성 실패: %w", err)
	}

	// 각 파일 압축 해제
	for _, file := range reader.File {
		if err := fm.extractFile(file, destDir); err != nil {
			return err
		}
	}

	return nil
}

// extractFile은 ZIP 파일의 단일 항목을 추출
func (fm *FileManager) extractFile(file *zip.File, destDir string) error {
	filePath := filepath.Join(destDir, file.Name)

	if file.FileInfo().IsDir() {
		return os.MkdirAll(filePath, os.ModePerm)
	}

	// 대상 디렉토리 생성
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return fmt.Errorf("디렉토리 생성 실패: %w", err)
	}

	// 소스 파일 열기
	rc, err := file.Open()
	if err != nil {
		return fmt.Errorf("ZIP 엔트리 열기 실패: %w", err)
	}
	defer rc.Close()

	// 대상 파일 생성
	outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return fmt.Errorf("파일 생성 실패: %w", err)
	}
	defer outFile.Close()

	// 파일 복사
	if _, err := io.Copy(outFile, rc); err != nil {
		return fmt.Errorf("파일 복사 실패: %w", err)
	}

	return nil
}

// CleanDirectory는 디렉토리의 모든 내용을 삭제
func (fm *FileManager) CleanDirectory(dir string) error {
	fm.logger.Info("디렉토리 정리 시작: %s", dir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("디렉토리 읽기 실패: %w", err)
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if err := fm.SafeRemove(path); err != nil {
			return fmt.Errorf("항목 삭제 실패 (%s): %w", entry.Name(), err)
		}
	}

	return nil
}

// SafeRemove는 파일이나 디렉토리를 안전하게 제거
func (fm *FileManager) SafeRemove(path string) error {
	const maxRetries = 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		var err error
		if info, err := os.Lstat(path); err == nil {
			if info.IsDir() {
				err = os.RemoveAll(path)
			} else {
				err = os.Remove(path)
			}
		}

		if err == nil || os.IsNotExist(err) {
			return nil
		}

		lastErr = err
		time.Sleep(time.Second)
	}

	// 마지막 시도: 이름 변경 후 삭제
	tempPath := path + ".old"
	if err := os.Rename(path, tempPath); err == nil {
		if err := os.RemoveAll(tempPath); err != nil {
			return fmt.Errorf("임시 파일 삭제 실패: %w", err)
		}
		return nil
	}

	return fmt.Errorf("삭제 실패 (최대 재시도 횟수 초과): %w", lastErr)
}
