package usecase

import (
	"fmt"
	"sync"

	"github.com/kihyun1998/dupdater/internal/i18n/domain/entity"
	"github.com/kihyun1998/dupdater/internal/i18n/domain/repository"
)

// LocaleService는 국제화 서비스를 정의합니다
type LocaleService struct {
	messageRepo   repository.MessageRepository
	currentLocale *entity.Locale
	mu            sync.RWMutex
}

// LocaleConfig는 LocaleService 생성에 필요한 설정입니다
type LocaleConfig struct {
	DefaultLocale *entity.Locale
	Providers     []repository.MessageProvider
}

// NewLocaleService는 새로운 LocaleService를 생성합니다
func NewLocaleService(repo repository.MessageRepository, config LocaleConfig) (*LocaleService, error) {
	service := &LocaleService{
		messageRepo:   repo,
		currentLocale: config.DefaultLocale,
	}

	// 제공자들의 메시지 등록
	for _, provider := range config.Providers {
		messages := make([]*entity.Message, 0)
		for key, content := range provider.GetMessages() {
			messages = append(messages, entity.NewMessage(key, content, provider.GetLocale()))
		}
		if err := repo.AddMessages(messages); err != nil {
			return nil, fmt.Errorf("메시지 등록 실패: %w", err)
		}
	}

	return service, nil
}

// SetLocale은 현재 로케일을 설정합니다
func (s *LocaleService) SetLocale(locale *entity.Locale) error {
	if !locale.IsValid() {
		return fmt.Errorf("지원하지 않는 언어입니다: %s", locale.String())
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentLocale = locale
	return nil
}

// GetMessage는 키에 해당하는 메시지를 현재 로케일로 반환합니다
func (s *LocaleService) GetMessage(key string, args ...interface{}) string {
	s.mu.RLock()
	locale := s.currentLocale
	s.mu.RUnlock()

	msg, err := s.messageRepo.GetMessage(key, locale)
	if err != nil {
		return key
	}

	return msg.Format(args...)
}

// GetCurrentLocale은 현재 로케일을 반환합니다
func (s *LocaleService) GetCurrentLocale() *entity.Locale {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentLocale
}

// GetAvailableLocales는 사용 가능한 모든 로케일을 반환합니다
func (s *LocaleService) GetAvailableLocales() []*entity.Locale {
	return s.messageRepo.GetAvailableLocales()
}
