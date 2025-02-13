// Package hash는 해시 도메인의 진입점을 제공합니다
package hash

import (
	"github.com/kihyun1998/dupdater/internal/hash/domain/repository"
	"github.com/kihyun1998/dupdater/internal/hash/domain/usecase"
	"github.com/kihyun1998/dupdater/internal/hash/infrastructure"
	"github.com/kihyun1998/dupdater/internal/logger"
)

// Config는 해시 매니저 생성에 필요한 설정입니다
type Config struct {
	CurrentDir  string
	HashSumPath string
	Logger      logger.Logger
}

// Manager는 해시 검증을 위한 인터페이스입니다
type Manager interface {
	VerifyFile(filePath string, expectedHash string) error
	VerifyUpdateFile(filePath string) error
	VerifyHashSum() error
}

// 실제 구현체
type manager struct {
	service *usecase.HashService
}

// New는 새로운 해시 매니저를 생성합니다
func New(config Config) (Manager, error) {
	// 리포지토리 설정
	repoConfig := &repository.Config{
		CurrentDir:  config.CurrentDir,
		HashSumPath: config.HashSumPath,
		Logger:      config.Logger,
	}

	// 해시 매니저 생성
	repo := infrastructure.NewFileHashManager(repoConfig)

	// 서비스 생성
	service := usecase.NewHashService(repo, config.Logger)

	return &manager{
		service: service,
	}, nil
}

// 인터페이스 구현
func (m *manager) VerifyFile(filePath string, expectedHash string) error {
	return m.service.VerifyFile(filePath, expectedHash)
}

func (m *manager) VerifyUpdateFile(filePath string) error {
	return m.service.VerifyUpdateFile(filePath)
}

func (m *manager) VerifyHashSum() error {
	return m.service.VerifyHashSum()
}
