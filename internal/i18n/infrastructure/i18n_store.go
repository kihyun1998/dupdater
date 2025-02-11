package infrastructure

import (
	"fmt"
	"sync"

	"github.com/kihyun1998/dupdater/internal/i18n/domain/entity"
	"github.com/kihyun1998/dupdater/internal/i18n/domain/repository"
	"github.com/kihyun1998/dupdater/internal/logger"
)

// I18nStore는 메모리 기반 다국어 저장소입니다
type I18nStore struct {
	currentLang entity.Language
	messages    map[entity.Language]map[string]string
	logger      logger.Logger
	mu          sync.RWMutex
}

// NewI18nStore는 새로운 I18nStore 인스턴스를 생성합니다
func NewI18nStore(config *repository.Config) *I18nStore {
	return &I18nStore{
		currentLang: config.DefaultLanguage,
		messages:    make(map[entity.Language]map[string]string),
		logger:      config.Logger,
	}
}

// GetMessage는 현재 언어의 메시지를 조회합니다
func (s *I18nStore) GetMessage(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if messages, ok := s.messages[s.currentLang]; ok {
		if msg, ok := messages[key]; ok {
			return msg
		}
	}
	return key
}

// SetLanguage는 현재 언어를 설정합니다
func (s *I18nStore) SetLanguage(lang entity.Language) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.messages[lang]; !ok {
		return fmt.Errorf("지원하지 않는 언어: %s", lang)
	}

	s.currentLang = lang
	return nil
}

// GetCurrentLanguage는 현재 설정된 언어를 반환합니다
func (s *I18nStore) GetCurrentLanguage() entity.Language {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentLang
}

// RegisterProvider는 새로운 메시지 제공자를 등록합니다
func (s *I18nStore) RegisterProvider(provider entity.MessageProvider) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	lang := provider.GetLanguage()
	if !entity.NewLocale(lang, nil).IsValid() {
		return fmt.Errorf("유효하지 않은 언어: %s", lang)
	}

	s.messages[lang] = provider.GetMessages()
	return nil
}
