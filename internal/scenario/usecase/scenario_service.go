package usecase

import (
	"fmt"
	"time"

	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/scenario/entity"
	"github.com/kihyun1998/dupdater/internal/ui"
)

// ScenarioService는 시나리오 실행을 담당하는 서비스입니다
type ScenarioService struct {
	logger    logger.Logger
	uiManager ui.Manager
	scenario  *entity.Scenario
}

// NewScenarioService는 새로운 ScenarioService를 생성합니다
func NewScenarioService(
	logger logger.Logger,
	uiManager ui.Manager,
	scenario *entity.Scenario,
) *ScenarioService {
	return &ScenarioService{
		logger:    logger,
		uiManager: uiManager,
		scenario:  scenario,
	}
}

// Run은 시나리오를 실행합니다
func (s *ScenarioService) Run() {
	s.logger.Info("시나리오 시작: %s", s.scenario.Description)

	// UI 복원 핸들러 설정
	s.uiManager.SetRestoreHandler(func() {
		s.handleRestore()
	})

	// 시나리오 실행
	go func() {
		if s.scenario.IsErrorScenario() {
			s.runErrorScenario()
		} else {
			s.runSuccessScenario()
		}
	}()

	// UI 실행
	s.uiManager.Run()
}

// runSuccessScenario는 성공 시나리오를 실행합니다
func (s *ScenarioService) runSuccessScenario() {
	for i, _ := range entity.Steps {
		s.uiManager.SetCurrentStep(i)
		time.Sleep(2 * time.Second) // 단계별 지연
	}
}

// runErrorScenario는 에러 시나리오를 실행합니다
func (s *ScenarioService) runErrorScenario() {
	targetStep := s.scenario.TargetStep

	// 목표 단계까지 정상 진행
	for i := 0; i <= targetStep.Index; i++ {
		s.uiManager.SetCurrentStep(i)

		if i == targetStep.Index {
			// 에러 발생 단계
			time.Sleep(1 * time.Second)
			errMsg := fmt.Sprintf("%s 단계에서 오류 발생", targetStep.Name)
			s.uiManager.ShowError(fmt.Errorf(errMsg))
			return
		}

		time.Sleep(2 * time.Second)
	}
}

// handleRestore는 복원 프로세스를 처리합니다
func (s *ScenarioService) handleRestore() {
	if !s.scenario.TargetStep.NeedRestore {
		return
	}

	s.logger.Info("복원 프로세스 시작")
	time.Sleep(3 * time.Second) // 복원 시간 시뮬레이션
	s.uiManager.ShowRestoreComplete()
}
