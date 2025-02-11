// domain/repository/version_repo.go
package repository

import (
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/version/domain/entity"
)

// VersionRepository는 버전 관리를 위한 저장소 인터페이스입니다
type VersionRepository interface {
	// GetFromVersion은 현재 버전 정보를 조회합니다
	GetFromVersion() *entity.Version

	// GetToVersion은 대상 버전 정보를 조회합니다
	GetToVersion() *entity.Version

	// SaveFromVersion은 현재 버전 정보를 저장합니다
	SaveFromVersion(version *entity.Version) error

	// SaveToVersion은 대상 버전 정보를 저장합니다
	SaveToVersion(version *entity.Version) error

	// ValidateVersions는 버전 정보의 유효성을 검증합니다
	ValidateVersions() error
}

// Config는 저장소 설정을 정의합니다
type Config struct {
	FromVersion string        // 현재 버전
	ToVersion   string        // 대상 버전
	Logger      logger.Logger // 로거 인터페이스
}
