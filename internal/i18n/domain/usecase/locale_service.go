// Package usecase는 i18n 도메인의 비즈니스 로직을 구현합니다
package usecase

import (
	"fmt"
	"sync"

	"github.com/kihyun1998/dupdater/internal/i18n/domain/entity"
	"github.com/kihyun1998/dupdater/internal/i18n/domain/repository"
)

// LocaleService는 로케일 관련 비즈니스 로직을 구현합니다
type LocaleService struct {
	currentLang string
	mutex       sync.RWMutex
	repository  repository.MessageRepository
}

// NewLocaleService는 새로운 LocaleService 인스턴스를 생성합니다
func NewLocaleService(repo repository.MessageRepository, defaultLang string) *LocaleService {
	return &LocaleService{
		currentLang: defaultLang,
		repository:  repo,
	}
}

// SetLanguage는 현재 언어를 설정합니다
func (s *LocaleService) SetLanguage(lang string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	locale := entity.NewLocale(entity.Language(lang), "")
	if !locale.IsValid() {
		return fmt.Errorf("지원하지 않는 언어입니다: %s", lang)
	}

	// 해당 언어의 메시지가 존재하는지 확인
	_, err := s.repository.GetAllMessages(locale)
	if err != nil {
		return fmt.Errorf("언어 리소스를 찾을 수 없습니다: %s", lang)
	}

	s.currentLang = lang
	return nil
}

// GetMessage는 지정된 키에 해당하는 메시지를 현재 설정된 언어로 반환합니다
func (s *LocaleService) GetMessage(key string) string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	locale := entity.NewLocale(entity.Language(s.currentLang), "")
	message, err := s.repository.GetMessage(key, locale)
	if err != nil || message == nil {
		return key // 메시지를 찾을 수 없는 경우 키를 반환
	}

	return message.Value
}

// GetCurrentLanguage는 현재 설정된 언어를 반환합니다
func (s *LocaleService) GetCurrentLanguage() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.currentLang
}

// RegisterProvider는 새로운 메시지 제공자를 등록합니다
func (s *LocaleService) RegisterProvider(provider entity.MessageProvider) error {
	return s.repository.RegisterProvider(provider)
}
