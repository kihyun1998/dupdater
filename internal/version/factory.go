package version

import (
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/version/domain/repository"
	"github.com/kihyun1998/dupdater/internal/version/domain/usecase"
	"github.com/kihyun1998/dupdater/internal/version/infrastructure"
)

// Manager는 버전 관리를 위한 외부 인터페이스입니다
type Manager interface {
	// GetFromVersion은 현재 버전을 반환합니다
	GetFromVersion() string

	// GetToVersion은 대상 버전을 반환합니다
	GetToVersion() string
}

// Config는 Manager 생성에 필요한 설정을 담는 구조체입니다
type Config struct {
	Logger      logger.Logger
	FromVersion string
	ToVersion   string
}

// managerImpl은 Manager 인터페이스의 구현체입니다
type managerImpl struct {
	service *usecase.VersionService
}

// New는 새로운 Manager 인스턴스를 생성합니다
func New(config Config) (Manager, error) {
	// 저장소 설정 생성
	repoConfig := &repository.Config{
		FromVersion: config.FromVersion,
		ToVersion:   config.ToVersion,
		Logger:      config.Logger,
	}

	// 버전 저장소 생성
	store, err := infrastructure.NewVersionStore(repoConfig)
	if err != nil {
		return nil, err
	}

	// 버전 서비스 생성
	service := usecase.NewVersionService(store, config.Logger)

	// 버전 정보 유효성 검증
	if err := service.ValidateVersions(); err != nil {
		return nil, err
	}

	return &managerImpl{
		service: service,
	}, nil
}

// GetFromVersion은 현재 버전을 반환합니다
func (m *managerImpl) GetFromVersion() string {
	return m.service.GetFromVersion()
}

// GetToVersion은 대상 버전을 반환합니다
func (m *managerImpl) GetToVersion() string {
	return m.service.GetToVersion()
}
