package scenarios

import (
	"time"

	"github.com/kihyun1998/dupdater/internal/scenario/entity"
	"github.com/kihyun1998/dupdater/internal/ui"
)

// StepExecutor는 각 단계의 실행을 담당하는 인터페이스입니다
type StepExecutor interface {
	Execute(uiManager ui.Manager) error
	GetStep() *entity.Step
}

// BaseScenario는 기본 시나리오 구현을 제공합니다
type BaseScenario struct {
	executors []StepExecutor
}

// Execute는 시나리오의 단계들을 순차적으로 실행합니다
func (b *BaseScenario) Execute(uiManager ui.Manager) error {
	for _, executor := range b.executors {
		step := executor.GetStep()
		uiManager.SetCurrentStep(step.Index)
		time.Sleep(2 * time.Second) // 단계 실행 시뮬레이션

		if err := executor.Execute(uiManager); err != nil {
			return err
		}
	}
	return nil
}
