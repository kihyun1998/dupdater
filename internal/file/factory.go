package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kihyun1998/dupdater/internal/file/infrastructure"
	loggerPorts "github.com/kihyun1998/dupdater/internal/logger/domain/ports"
)

// Config는 Manager 생성에 필요한 설정을 담는 구조체입니다
type Config struct {
	Logger     loggerPorts.LoggerPort
	BackupDir  string
	CurrentDir string
}

// Manager는 파일 작업을 관리하는 구조체입니다
type Manager struct {
	fileManager *infrastructure.FileManager
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

	fileManager := infrastructure.NewFileManager(
		config.Logger,
		config.BackupDir,
		config.CurrentDir,
	)

	return &Manager{
		fileManager: fileManager,
	}, nil
}

// 기존 인터페이스와 동일한 메서드들
func (m *Manager) Backup() error {
	return m.fileManager.Backup()
}

func (m *Manager) Restore() error {
	return m.fileManager.Restore()
}

func (m *Manager) ExtractZip(zipFile string) error {
	return m.fileManager.ExtractZip(zipFile)
}

func (m *Manager) DeleteFile(path string) error {
	return m.fileManager.DeleteFile(path)
}

// Logger 인터페이스
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}
