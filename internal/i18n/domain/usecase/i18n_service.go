package usecase

import (
	"fmt"
	"sync"

	"github.com/kihyun1998/dupdater/internal/i18n/domain/entity"
	"github.com/kihyun1998/dupdater/internal/i18n/domain/repository"
)

// I18nService는 다국어 지원의 비즈니스 로직을 구현합니다
type I18nService struct {
	repo   repository.I18nRepository
	logger repository.Logger
	mu     sync.RWMutex
}

// NewI18nService는 새로운 I18nService 인스턴스를 생성합니다
func NewI18nService(repo repository.I18nRepository, logger repository.Logger) *I18nService {
	return &I18nService{
		repo:   repo,
		logger: logger,
	}
}

// GetMessage는 지정된 키에 해당하는 메시지를 현재 언어로 반환합니다
func (s *I18nService) GetMessage(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	msg := s.repo.GetMessage(key)
	if msg == key {
		s.logger.Error("메시지를 찾을 수 없음: %s", key)
	}
	return msg
}

// SetLanguage는 현재 언어를 설정합니다
func (s *I18nService) SetLanguage(lang string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.repo.SetLanguage(entity.Language(lang)); err != nil {
		s.logger.Error("언어 설정 실패: %v", err)
		return fmt.Errorf("언어 설정 실패: %w", err)
	}

	s.logger.Info("언어가 변경됨: %s", lang)
	return nil
}

// GetCurrentLanguage는 현재 설정된 언어를 반환합니다
func (s *I18nService) GetCurrentLanguage() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return string(s.repo.GetCurrentLanguage())
}

// RegisterProvider는 새로운 메시지 제공자를 등록합니다
func (s *I18nService) RegisterProvider(provider entity.MessageProvider) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.repo.RegisterProvider(provider); err != nil {
		s.logger.Error("메시지 제공자 등록 실패: %v", err)
		return fmt.Errorf("메시지 제공자 등록 실패: %w", err)
	}

	s.logger.Info("메시지 제공자가 등록됨: %s", provider.GetLanguage())
	return nil
}
