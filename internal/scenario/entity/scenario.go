// internal/scenario/entity/scenario.go
package entity

import "fmt"

// ScenarioType은 시나리오 유형을 정의합니다
type ScenarioType string

const (
	ScenarioSuccess ScenarioType = "success"
	ScenarioError1  ScenarioType = "error1"
	ScenarioError2  ScenarioType = "error2"
	ScenarioError3  ScenarioType = "error3"
	ScenarioError4  ScenarioType = "error4"
	ScenarioError5  ScenarioType = "error5"
	ScenarioError6  ScenarioType = "error6"
	ScenarioError7  ScenarioType = "error7"
	ScenarioError8  ScenarioType = "error8"
)

// Scenario는 테스트 시나리오의 기본 정보를 정의합니다
type Scenario struct {
	Type        ScenarioType // 시나리오 유형
	TargetStep  *Step        // 목표 단계 (에러 발생 단계)
	Description string       // 시나리오 설명
}

// NewScenario는 새로운 시나리오를 생성합니다
func NewScenario(scenarioType ScenarioType) *Scenario {
	s := &Scenario{
		Type: scenarioType,
	}

	// 시나리오 유형에 따른 설정
	switch scenarioType {
	case ScenarioSuccess:
		s.Description = "정상 업데이트 시나리오"
		s.TargetStep = nil
	default:
		// error1 ~ error8 처리
		stepIndex := int(scenarioType[5] - '1') // "error1"에서 숫자 추출
		if step := GetStep(stepIndex); step != nil {
			s.TargetStep = step
			s.Description = fmt.Sprintf("%s 단계 실패 시나리오", step.Name)
		}
	}

	return s
}

// IsErrorScenario는 현재 시나리오가 에러 시나리오인지 확인합니다
func (s *Scenario) IsErrorScenario() bool {
	return s.Type != ScenarioSuccess
}
