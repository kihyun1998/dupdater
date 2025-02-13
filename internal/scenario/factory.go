// internal/scenario/factory.go
package scenario

import (
	"fmt"

	"github.com/kihyun1998/dupdater/internal/i18n"
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/scenario/entity"
	"github.com/kihyun1998/dupdater/internal/scenario/usecase"
	"github.com/kihyun1998/dupdater/internal/ui"
)

// Manager는 시나리오 실행을 위한 인터페이스입니다
type Manager interface {
	Run()
}

// Config는 시나리오 매니저 생성에 필요한 설정입니다
type Config struct {
	Logger       logger.Logger
	UIManager    ui.Manager
	I18n         i18n.Manager
	ScenarioType string
}

// manager는 시나리오 매니저의 구현체입니다
type manager struct {
	service *usecase.ScenarioService
}

// New는 새로운 시나리오 매니저를 생성합니다
func New(config Config) (Manager, error) {
	// 시나리오 타입 검증 및 생성
	scenarioType := entity.ScenarioType(config.ScenarioType)
	scenario := entity.NewScenario(scenarioType)

	if scenario == nil {
		return nil, fmt.Errorf("잘못된 시나리오 타입: %s", config.ScenarioType)
	}

	// 시나리오 서비스 생성
	service := usecase.NewScenarioService(
		config.Logger,
		config.UIManager,
		scenario,
	)

	return &manager{
		service: service,
	}, nil
}

// Run은 시나리오를 실행합니다
func (m *manager) Run() {
	m.service.Run()
}
