// internal/ui/domain/usecase/ui_service.go

package usecase

import (
	"sync"

	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/ui/domain/repository"
)

// UIService는 UI 관련 비즈니스 로직을 처리합니다
type UIService struct {
	repo   repository.UIRepository
	logger logger.Logger
	mu     sync.RWMutex
}

// NewUIService는 새로운 UIService 인스턴스를 생성합니다
func NewUIService(repo repository.UIRepository, logger logger.Logger) *UIService {
	return &UIService{
		repo:   repo,
		logger: logger,
	}
}

// SetCurrentStep은 현재 진행 단계를 설정합니다
func (s *UIService) SetCurrentStep(step int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state := s.repo.GetState()
	state.UpdateProgress(step, state.DetailMessage)
	s.repo.UpdateState(state)
	s.repo.SetCurrentStep(step)
	s.logger.Info("현재 단계 업데이트: %d", step)
}

// UpdateDetail은 상세 메시지를 업데이트합니다
func (s *UIService) UpdateDetail(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state := s.repo.GetState()
	state.DetailMessage = message
	s.repo.UpdateState(state)
	s.repo.UpdateDetail(message)
	s.logger.Info("상세 메시지 업데이트: %s", message)
}

// ShowError는 에러 상태를 표시합니다
func (s *UIService) ShowError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state := s.repo.GetState()
	state.SetError()
	s.repo.UpdateState(state)
	s.repo.ShowError(err)
	s.logger.Error("에러 발생: %v", err)
}

// Run은 UI를 실행합니다
func (s *UIService) Run() {
	s.logger.Info("UI 실행")
	s.repo.Run()
}

// Close는 UI를 종료합니다
func (s *UIService) Close() {
	s.logger.Info("UI 종료")
	s.repo.Close()
}

// SetRestoreHandler는 복원 핸들러를 설정합니다
func (s *UIService) SetRestoreHandler(handler func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.repo.SetRestoreHandler(handler)
	s.logger.Info("복원 핸들러 설정됨")
}

// ShowRestoring은 복원 중 상태를 표시합니다
func (s *UIService) ShowRestoring() {
	s.mu.Lock()
	defer s.mu.Unlock()

	state := s.repo.GetState()
	state.SetRestoring()
	s.repo.UpdateState(state)
	s.repo.ShowRestoring()
	s.logger.Info("복원 중 상태로 변경")
}

// ShowRestoreComplete는 복원 완료 상태를 표시합니다
func (s *UIService) ShowRestoreComplete() {
	s.mu.Lock()
	defer s.mu.Unlock()

	state := s.repo.GetState()
	state.SetRestoreCompleted()
	s.repo.UpdateState(state)
	s.repo.ShowRestoreComplete()
	s.logger.Info("복원 완료 상태로 변경")
}

// GetTotalSteps는 전체 단계 수를 반환합니다
func (s *UIService) GetTotalSteps() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.repo.GetTotalSteps()
}

// SetCompletionCallback은 완료 콜백을 설정합니다
func (s *UIService) SetCompletionCallback(callback func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.repo.SetCompletionCallback(callback)
	s.logger.Info("완료 콜백 설정됨")
}
