// Package usecase는 해시 도메인의 비즈니스 로직을 구현합니다
package usecase

import (
	"fmt"

	"github.com/kihyun1998/dupdater/internal/hash/domain/repository"
)

// HashService는 해시 검증 관련 비즈니스 로직을 구현합니다
type HashService struct {
	repo   repository.HashRepository
	logger repository.Logger
}

// NewHashService는 새로운 HashService 인스턴스를 생성합니다
func NewHashService(repo repository.HashRepository, logger repository.Logger) *HashService {
	return &HashService{
		repo:   repo,
		logger: logger,
	}
}

// VerifyFile은 단일 파일의 해시를 검증합니다
func (s *HashService) VerifyFile(filePath string, expectedHash string) error {
	s.logger.Info("파일 해시 검증 시작: %s", filePath)

	if err := s.repo.VerifyFile(filePath, expectedHash); err != nil {
		s.logger.Error("파일 해시 검증 실패: %v", err)
		return fmt.Errorf("파일 해시 검증 실패: %w", err)
	}

	s.logger.Info("파일 해시 검증 완료: %s", filePath)
	return nil
}

// VerifyUpdateFile은 업데이트 파일의 해시를 검증합니다
func (s *HashService) VerifyUpdateFile(filePath string) error {
	s.logger.Info("업데이트 파일 해시 검증 시작: %s", filePath)

	if err := s.repo.VerifyUpdateFile(filePath); err != nil {
		s.logger.Error("업데이트 파일 해시 검증 실패: %v", err)
		return fmt.Errorf("업데이트 파일 해시 검증 실패: %w", err)
	}

	s.logger.Info("업데이트 파일 해시 검증 완료: %s", filePath)
	return nil
}

// VerifyHashSum은 모든 파일의 해시섬을 검증합니다
func (s *HashService) VerifyHashSum() error {
	s.logger.Info("전체 파일 해시섬 검증 시작")

	if err := s.repo.VerifyHashSum(); err != nil {
		s.logger.Error("전체 파일 해시섬 검증 실패: %v", err)
		return fmt.Errorf("전체 파일 해시섬 검증 실패: %w", err)
	}

	s.logger.Info("전체 파일 해시섬 검증 완료")
	return nil
}
