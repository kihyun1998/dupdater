// infrastructure/version_store.go
package infrastructure

import (
	"fmt"
	"sync"

	"github.com/kihyun1998/dupdater/internal/version/domain/entity"
	"github.com/kihyun1998/dupdater/internal/version/domain/repository"
)

// VersionStore는 버전 정보를 관리하는 저장소 구현체입니다
type VersionStore struct {
	fromVersion *entity.Version
	toVersion   *entity.Version
	logger      repository.Logger
	mu          sync.RWMutex
}

// NewVersionStore는 새로운 VersionStore 인스턴스를 생성합니다
func NewVersionStore(config *repository.Config) (*VersionStore, error) {
	fromVersion, err := entity.ParseVersion(config.FromVersion)
	if err != nil {
		return nil, fmt.Errorf("현재 버전 파싱 실패: %w", err)
	}

	toVersion, err := entity.ParseVersion(config.ToVersion)
	if err != nil {
		return nil, fmt.Errorf("대상 버전 파싱 실패: %w", err)
	}

	return &VersionStore{
		fromVersion: fromVersion,
		toVersion:   toVersion,
		logger:      config.Logger,
	}, nil
}

// GetFromVersion은 현재 버전을 반환합니다
func (s *VersionStore) GetFromVersion() *entity.Version {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.fromVersion
}

// GetToVersion은 대상 버전을 반환합니다
func (s *VersionStore) GetToVersion() *entity.Version {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.toVersion
}

// SaveFromVersion은 현재 버전을 저장합니다
func (s *VersionStore) SaveFromVersion(version *entity.Version) error {
	if err := version.Validate(); err != nil {
		return fmt.Errorf("현재 버전 검증 실패: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.fromVersion = version
	s.logger.Info("현재 버전이 업데이트됨: %s", version.String())
	return nil
}

// SaveToVersion은 대상 버전을 저장합니다
func (s *VersionStore) SaveToVersion(version *entity.Version) error {
	if err := version.Validate(); err != nil {
		return fmt.Errorf("대상 버전 검증 실패: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.toVersion = version
	s.logger.Info("대상 버전이 업데이트됨: %s", version.String())
	return nil
}

// ValidateVersions는 버전 정보의 유효성을 검증합니다
func (s *VersionStore) ValidateVersions() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.fromVersion == nil || s.toVersion == nil {
		return fmt.Errorf("버전 정보가 초기화되지 않았습니다")
	}

	if err := s.fromVersion.Validate(); err != nil {
		return fmt.Errorf("현재 버전 검증 실패: %w", err)
	}

	if err := s.toVersion.Validate(); err != nil {
		return fmt.Errorf("대상 버전 검증 실패: %w", err)
	}

	if !s.toVersion.IsNewer(s.fromVersion) {
		return fmt.Errorf("대상 버전(%s)이 현재 버전(%s)보다 낮거나 같습니다",
			s.toVersion.String(), s.fromVersion.String())
	}

	return nil
}
