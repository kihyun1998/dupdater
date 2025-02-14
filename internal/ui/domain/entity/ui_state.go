package entity

import "fmt"

// UIState는 UI의 상태를 나타내는 도메인 엔티티입니다
type UIState struct {
	CurrentStep      int     // 현재 진행 단계
	TotalSteps       int     // 전체 단계 수
	Progress         float64 // 진행률 (0-1 사이 값)
	DetailMessage    string  // 상세 메시지
	FromVersion      string  // 시작 버전
	ToVersion        string  // 목표 버전
	IsError          bool    // 에러 상태 여부
	IsRestoring      bool    // 복원 중 상태 여부
	RestoreCompleted bool    // 복원 완료 상태
}

// NewUIState는 새로운 UIState 인스턴스를 생성합니다
func NewUIState(totalSteps int, fromVersion, toVersion string) *UIState {
	return &UIState{
		CurrentStep:      0,
		TotalSteps:       totalSteps,
		Progress:         0.0,
		FromVersion:      fromVersion,
		ToVersion:        toVersion,
		IsError:          false,
		IsRestoring:      false,
		RestoreCompleted: false,
	}
}

// UpdateProgress는 현재 진행 상태를 업데이트합니다
func (s *UIState) UpdateProgress(step int, message string) {
	s.CurrentStep = step
	s.DetailMessage = message
	if s.TotalSteps > 0 {
		s.Progress = float64(step) / float64(s.TotalSteps-1)
	}
}

// SetError는 에러 상태로 변경합니다
func (s *UIState) SetError() {
	s.IsError = true
}

// SetRestoring는 복원 중 상태로 변경합니다
func (s *UIState) SetRestoring() {
	s.IsRestoring = true
	s.IsError = false
}

// SetRestoreCompleted는 복원 완료 상태로 변경합니다
func (s *UIState) SetRestoreCompleted() {
	s.RestoreCompleted = true
	s.IsRestoring = false
}

// IsCompleted는 업데이트가 완료되었는지 확인합니다
func (s *UIState) IsCompleted() bool {
	return s.Progress >= 1.0
}

// Validate는 상태가 유효한지 검증합니다
func (s *UIState) Validate() error {
	if s.TotalSteps <= 0 {
		return fmt.Errorf("전체 단계 수는 0보다 커야 합니다")
	}
	if s.FromVersion == "" || s.ToVersion == "" {
		return fmt.Errorf("버전 정보가 누락되었습니다")
	}
	return nil
}
