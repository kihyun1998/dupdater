// internal/scenario/scenarios/error.go
package scenarios

import (
	"fmt"

	"github.com/kihyun1998/dupdater/internal/scenario/entity"
	"github.com/kihyun1998/dupdater/internal/ui"
)

// ErrorStep은 에러 발생 단계를 구현합니다
type ErrorStep struct {
	step      *entity.Step
	shouldErr bool
}

func NewErrorStep(step *entity.Step, shouldErr bool) *ErrorStep {
	return &ErrorStep{
		step:      step,
		shouldErr: shouldErr,
	}
}

func (s *ErrorStep) Execute(uiManager ui.Manager) error {
	if s.shouldErr {
		return fmt.Errorf("%s 단계에서 오류 발생", s.step.Name)
	}
	return nil
}

func (s *ErrorStep) GetStep() *entity.Step {
	return s.step
}

// ErrorScenario는 에러 시나리오를 구현합니다
type ErrorScenario struct {
	BaseScenario
	errorStep *entity.Step
}

// NewErrorScenario는 새로운 에러 시나리오를 생성합니다
func NewErrorScenario(errorStep *entity.Step) *ErrorScenario {
	scenario := &ErrorScenario{
		errorStep: errorStep,
	}

	// 에러 발생 단계까지의 실행자 추가
	for _, step := range entity.Steps {
		shouldErr := step.Index == errorStep.Index
		scenario.executors = append(scenario.executors, NewErrorStep(step, shouldErr))
		if shouldErr {
			break
		}
	}

	return scenario
}
