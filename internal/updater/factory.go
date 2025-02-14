package updater

import (
	"fmt"

	"github.com/kihyun1998/dupdater/internal/file"
	"github.com/kihyun1998/dupdater/internal/hash"
	"github.com/kihyun1998/dupdater/internal/i18n"
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/network"
	"github.com/kihyun1998/dupdater/internal/ui"
	"github.com/kihyun1998/dupdater/internal/updater/domain/entity"
	"github.com/kihyun1998/dupdater/internal/updater/domain/repository"
	"github.com/kihyun1998/dupdater/internal/updater/domain/usecase"
	"github.com/kihyun1998/dupdater/internal/updater/infrastructure"
)

// Manager는 업데이트 관리를 위한 인터페이스입니다
type Manager interface {
	Start()
}

// Config는 업데이터 생성에 필요한 설정입니다
type Config struct {
	AppName          string
	FromVersion      string
	ServerName       string
	BackupCompleted  bool
	RestoreCompleted bool
	UIManager        ui.Manager
	Logger           logger.Logger
	NetworkManager   network.Manager
	FileManager      file.Manager
	HashManager      hash.Manager
	I18n             i18n.Manager
}

// manager는 Manager 인터페이스의 구현체입니다
type manager struct {
	service *usecase.UpdateService
}

// New는 새로운 업데이트 Manager를 생성합니다
func New(config Config) (Manager, error) {
	// 1. 엔티티 생성
	updateConfig, err := entity.NewUpdateConfig(
		config.AppName,
		config.FromVersion,
		config.ServerName,
	)
	if err != nil {
		return nil, fmt.Errorf("업데이트 설정 생성 실패: %w", err)
	}

	updateStatus := entity.NewUpdateStatus()
	updateStatus.BackupCompleted = config.BackupCompleted
	updateStatus.RestoreCompleted = config.RestoreCompleted

	// 2. Repository 계층 설정
	repoConfig := &repository.Config{
		Config:      updateConfig,
		Status:      updateStatus,
		UI:          config.UIManager,
		Logger:      config.Logger,
		Network:     config.NetworkManager,
		FileManager: config.FileManager,
		HashManager: config.HashManager,
		I18n:        config.I18n,
	}

	// 3. Repository 구현체 생성
	repo, err := infrastructure.NewUpdateManager(repoConfig)
	if err != nil {
		return nil, fmt.Errorf("업데이트 매니저 생성 실패: %w", err)
	}

	// 4. Service 계층 생성
	service := usecase.NewUpdateService(
		repo,
		config.UIManager,
		config.Logger,
		updateStatus,
		updateConfig,
	)

	return &manager{
		service: service,
	}, nil
}

// Start는 업데이트 프로세스를 시작합니다
func (m *manager) Start() {
	m.service.Start()
}
