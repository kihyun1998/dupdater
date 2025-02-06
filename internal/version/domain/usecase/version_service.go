// domain/usecase/version_service.go
package usecase

import (
	"fmt"

	"github.com/kihyun1998/dupdater/internal/version/domain/repository"
)

// VersionService는 버전 관리의 비즈니스 로직을 구현합니다
type VersionService struct {
	repo   repository.VersionRepository
	logger repository.Logger
}

// NewVersionService는 새로운 VersionService 인스턴스를 생성합니다
func NewVersionService(repo repository.VersionRepository, logger repository.Logger) *VersionService {
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

// IsUpdateRequired는 업데이트가 필요한지 확인합니다
func (s *VersionService) IsUpdateRequired() bool {
	fromVersion := s.repo.GetFromVersion()
	toVersion := s.repo.GetToVersion()

	if fromVersion == nil || toVersion == nil {
		s.logger.Error("버전 정보가 없어 업데이트 필요 여부를 확인할 수 없습니다")
		return false
	}

	return toVersion.IsNewer(fromVersion)
}
