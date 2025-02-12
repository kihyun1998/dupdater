package usecase

import (
	"fmt"

	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/version/domain/repository"
)

// VersionService는 버전 관리의 비즈니스 로직을 구현합니다
type VersionService struct {
	repo   repository.VersionRepository
	logger logger.Logger
}

// NewVersionService는 새로운 VersionService 인스턴스를 생성합니다
func NewVersionService(repo repository.VersionRepository, logger logger.Logger) *VersionService {
	return &VersionService{
		repo:   repo,
		logger: logger,
	}
}

// GetFromVersion은 현재 버전을 반환합니다
func (s *VersionService) GetFromVersion() string {
	version := s.repo.GetFromVersion()
	if version == nil {
		s.logger.Error("현재 버전 정보를 찾을 수 없습니다")
		return ""
	}
	return version.String()
}

// GetToVersion은 대상 버전을 반환합니다
func (s *VersionService) GetToVersion() string {
	version := s.repo.GetToVersion()
	if version == nil {
		s.logger.Error("대상 버전 정보를 찾을 수 없습니다")
		return ""
	}
	return version.String()
}

// ValidateVersions는 버전 정보의 유효성을 검사합니다
func (s *VersionService) ValidateVersions() error {
	if err := s.repo.ValidateVersions(); err != nil {
		s.logger.Error("버전 정보 검증 실패: %v", err)
		return fmt.Errorf("버전 정보 검증 실패: %w", err)
	}
	return nil
}
