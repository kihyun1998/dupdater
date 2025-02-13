package scenarios

import (
	"github.com/kihyun1998/dupdater/internal/scenario/entity"
	"github.com/kihyun1998/dupdater/internal/ui"
)

// SuccessStep은 성공 시나리오의 단계를 구현합니다
type SuccessStep struct {
	step *entity.Step
}

func NewSuccessStep(step *entity.Step) *SuccessStep {
	return &SuccessStep{step: step}
}

func (s *SuccessStep) Execute(uiManager ui.Manager) error {
	return nil // 성공 시나리오는 항상 성공
}

func (s *SuccessStep) GetStep() *entity.Step {
	return s.step
}

// SuccessScenario는 성공 시나리오를 구현합니다
type SuccessScenario struct {
	BaseScenario
}

// NewSuccessScenario는 새로운 성공 시나리오를 생성합니다
func NewSuccessScenario() *SuccessScenario {
	scenario := &SuccessScenario{}

	// 모든 단계를 성공 단계로 추가
	for _, step := range entity.Steps {
		scenario.executors = append(scenario.executors, NewSuccessStep(step))
	}

	return scenario
}
