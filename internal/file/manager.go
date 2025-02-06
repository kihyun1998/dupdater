// Package file은 파일 시스템 작업을 담당하는 패키지입니다
package file

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Manager는 파일 시스템 작업을 관리하는 구조체입니다
type Manager struct {
	logger     Logger // 로깅을 위한 인터페이스
	backupDir  string // 백업 디렉토리 경로
	currentDir string // 현재 작업 디렉토리
}

// Config는 Manager 생성에 필요한 설정을 담는 구조체입니다
type Config struct {
	Logger     Logger
	BackupDir  string // 백업 디렉토리 경로 (기본값: %TEMP%/ACRABACK)
	CurrentDir string // 현재 작업 디렉토리
}

// New는 새로운 Manager 인스턴스를 생성합니다
func New(config Config) (*Manager, error) {
	if config.BackupDir == "" {
		config.BackupDir = filepath.Join(os.TempDir(), "ACRABACK")
	}

	if config.CurrentDir == "" {
		dir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("현재 디렉토리 확인 실패: %w", err)
		}
		config.CurrentDir = dir
	}

	return &Manager{
		logger:     config.Logger,
		backupDir:  config.BackupDir,
		currentDir: config.CurrentDir,
	}, nil
}

// Backup은 현재 디렉토리의 파일들을 백업합니다
func (m *Manager) Backup() error {
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

// Restore는 백업된 파일들을 복원합니다
func (m *Manager) Restore() error {
	m.logger.Info("파일 복원 시작")

	// 백업 디렉토리 존재 확인
	if _, err := os.Stat(m.backupDir); os.IsNotExist(err) {
		m.logger.Error("백업 디렉토리가 존재하지 않음: %s", m.backupDir)
		return fmt.Errorf("백업 디렉토리가 존재하지 않습니다: %s", m.backupDir)
	}

	// 현재 디렉토리 정리
	if err := m.cleanCurrentDirectory(); err != nil {
		m.logger.Error("현재 디렉토리 정리 실패: %v", err)
		return fmt.Errorf("현재 디렉토리 정리 실패: %w", err)
	}

	// 백업 파일 복원
	files, err := os.ReadDir(m.backupDir)
	if err != nil {
		m.logger.Error("백업 디렉토리 읽기 실패: %v", err)
		return fmt.Errorf("백업 디렉토리 읽기 실패: %w", err)
	}

	for _, file := range files {
		srcPath := filepath.Join(m.backupDir, file.Name())
		destPath := filepath.Join(m.currentDir, file.Name())

		if file.IsDir() {
			m.logger.Info("폴더 복원 시도: %s -> %s", srcPath, destPath)
			if err := m.restoreDirectory(srcPath, destPath); err != nil {
				m.logger.Error("폴더 복원 실패: %s, 에러: %v", srcPath, err)
				return fmt.Errorf("디렉토리 복원 실패 (%s): %w", file.Name(), err)
			}
		} else {
			m.logger.Info("파일 복원 시도: %s -> %s", srcPath, destPath)
			if err := m.restoreFile(srcPath, destPath); err != nil {
				m.logger.Error("파일 복원 실패: %s, 에러: %v", srcPath, err)
				return fmt.Errorf("파일 복원 실패 (%s): %w", file.Name(), err)
			}
		}
	}

	m.logger.Info("파일 복원 완료")
	return nil
}

// ExtractZip은 ZIP 파일을 지정된 디렉토리에 압축 해제합니다
func (m *Manager) ExtractZip(zipFile string) error {
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

// DeleteFile은 지정된 파일을 삭제합니다
func (m *Manager) DeleteFile(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("파일 삭제 실패: %w", err)
	}
	return nil
}

func (m *Manager) GetBackupDir() string {
	return m.backupDir
}

// 내부 헬퍼 함수들
// 파일 백업 함수
func (m *Manager) backupFile(src, dest string) error {
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

// 디렉토리 백업 함수
func (m *Manager) backupDirectory(src, dest string) error {
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

// 파일 복구 함수
func (m *Manager) restoreFile(src, dest string) error {
	// 기존 파일이 존재하면 삭제 시도
	if _, err := os.Stat(dest); err == nil {
		if err := os.Remove(dest); err != nil {
			m.logger.Info("기존 파일 삭제 시도: %s", dest)
			if err := os.Remove(dest); err != nil {
				return fmt.Errorf("기존 파일 삭제 실패: %v", err)
			}
		}
	}

	// Rename 시도
	if err := os.Rename(src, dest); err == nil {
		m.logger.Info("파일 정상적으로 Rename으로 복원됨: %s -> %s", src, dest)
		return nil
	} else {
		// Rename 실패 로그 추가
		m.logger.Error("파일 Rename 실패, 복사 방식으로 복원 시도: %s -> %s, 에러: %v", src, dest, err)
	}

	// Rename 실패 시 기존 방식으로 복원 진행
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("파일 열기 실패: %v", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("파일 생성 실패: %v", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("파일 복사 실패: %v", err)
	}

	m.logger.Info("파일 복사 방식으로 복원됨: %s -> %s", src, dest)

	return os.Remove(src)
}

// 디렉토리 복구 함수
func (m *Manager) restoreDirectory(src, dest string) error {
	m.logger.Info("디렉토리 복원 시도: %s -> %s", src, dest)

	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		m.logger.Error("디렉토리 생성 실패: %s, 에러: %v", dest, err)
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		m.logger.Error("디렉토리 읽기 실패: %s, 에러: %v", src, err)
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			m.logger.Info("하위 디렉토리 복원 시도: %s -> %s", srcPath, destPath)
			if err := m.restoreDirectory(srcPath, destPath); err != nil {
				m.logger.Error("하위 디렉토리 복원 실패: %s, 에러: %v", srcPath, err)
				return err
			}
		} else {
			m.logger.Info("파일 복원 시도: %s -> %s", srcPath, destPath)
			if err := m.restoreFile(srcPath, destPath); err != nil {
				m.logger.Error("파일 복원 실패: %s, 에러: %v", srcPath, err)
				return err
			}
		}
	}

	m.logger.Info("디렉토리 복원 완료: %s", dest)
	return os.RemoveAll(src)
}

// 현재 디렉토리 정리 함수
func (m *Manager) cleanCurrentDirectory() error {
	files, err := os.ReadDir(m.currentDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		path := filepath.Join(m.currentDir, file.Name())
		for i := 0; i < 3; i++ {
			var err error
			if file.IsDir() {
				m.logger.Error("폴더 삭제 시도: %s", path)
				err = os.RemoveAll(path)
			} else {
				m.logger.Error("파일 삭제 시도: %s", path)
				err = os.Remove(path)
			}

			if err == nil {
				m.logger.Info("삭제 성공: %s", path)
				break
			}

			// 마지막 시도에는 old 붙이고 삭제 시도도
			if i == 2 {
				newPath := path + ".old"
				m.logger.Error("파일/폴더 삭제 실패, 이름 변경 후 삭제 시도: %s -> %s", path, newPath)
				if os.Rename(path, newPath) == nil {
					m.logger.Error(".old 파일 삭제 시도: %s", newPath)
					os.Remove(newPath)
				}
			}

			time.Sleep(time.Second)
		}
	}
	return nil
}

// 파일 압축 해제 함수
func (m *Manager) extractFile(file *zip.File) error {
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

// Logger는 로깅을 위한 인터페이스입니다
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}
