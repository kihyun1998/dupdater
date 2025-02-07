// Package file은 파일 시스템 작업의 진입점을 제공합니다
package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kihyun1998/dupdater/internal/file/domain/repository"
	"github.com/kihyun1998/dupdater/internal/file/domain/usecase"
	"github.com/kihyun1998/dupdater/internal/file/infrastructure"
)

// Manager는 파일 관리를 위한 인터페이스입니다
type Manager interface {
	Backup() error
	Restore() error
	ExtractZip(zipFile string) error
	DeleteFile(path string) error
	GetBackupDir() string
}

// Config는 Manager 생성에 필요한 설정입니다
type Config struct {
	Logger     Logger
	BackupDir  string
	CurrentDir string
}

// Logger는 로깅을 위한 인터페이스입니다
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}

// manager는 Manager 인터페이스의 구현체입니다
type manager struct {
	service *usecase.FileService
}

// New는 새로운 Manager 인스턴스를 생성합니다
func New(config Config) (Manager, error) {
	// 기본값 설정
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

	// Repository 설정
	repoConfig := &repository.Config{
		BackupDir:  config.BackupDir,
		CurrentDir: config.CurrentDir,
		Logger:     config.Logger,
	}

	// Repository 생성
	repo := infrastructure.NewLocalFileRepository(repoConfig)

	// Service 생성
	service := usecase.NewFileService(repo, config.Logger)

	return &manager{
		service: service,
	}, nil
}

// Manager 인터페이스 구현
func (m *manager) Backup() error {
	return m.service.Backup()
}

func (m *manager) Restore() error {
	return m.service.Restore()
}

func (m *manager) ExtractZip(zipFile string) error {
	return m.service.ExtractZip(zipFile)
}

func (m *manager) DeleteFile(path string) error {
	return m.service.DeleteFile(path)
}

func (m *manager) GetBackupDir() string {
	return m.service.GetBackupDir()
}
