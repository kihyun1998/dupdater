// Package infrastructure는 i18n 도메인의 실제 구현체들을 제공합니다
package infrastructure

import (
	"fmt"
	"sync"

	"github.com/kihyun1998/dupdater/internal/i18n/domain/entity"
	"github.com/kihyun1998/dupdater/internal/i18n/domain/repository"
)

// MessageStore는 메모리 기반 메시지 저장소입니다
type MessageStore struct {
	mutex    sync.RWMutex
	messages map[string]entity.MessageMap // locale code -> message map
	config   *repository.Config
}

// NewMessageStore는 새로운 MessageStore 인스턴스를 생성합니다
func NewMessageStore(config *repository.Config) (*MessageStore, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("저장소 설정 검증 실패: %w", err)
	}

	return &MessageStore{
		messages: make(map[string]entity.MessageMap),
		config:   config,
	}, nil
}

// GetMessage는 주어진 키와 로케일에 해당하는 메시지를 조회합니다
func (s *MessageStore) GetMessage(key string, locale *entity.Locale) (*entity.Message, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !locale.IsValid() {
		return nil, fmt.Errorf("유효하지 않은 로케일: %s", locale.String())
	}

	messageMap, exists := s.messages[locale.String()]
	if !exists {
		return nil, fmt.Errorf("로케일에 대한 메시지가 없음: %s", locale.String())
	}

	value, exists := messageMap[key]
	if !exists {
		return nil, fmt.Errorf("메시지를 찾을 수 없음: %s", key)
	}

	return entity.NewMessage(key, value, locale), nil
}

// GetAllMessages는 특정 로케일의 모든 메시지를 조회합니다
func (s *MessageStore) GetAllMessages(locale *entity.Locale) (entity.MessageMap, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !locale.IsValid() {
		return nil, fmt.Errorf("유효하지 않은 로케일: %s", locale.String())
	}

	messageMap, exists := s.messages[locale.String()]
	if !exists {
		return nil, fmt.Errorf("로케일에 대한 메시지가 없음: %s", locale.String())
	}

	// 메시지 맵 복사
	result := make(entity.MessageMap)
	for k, v := range messageMap {
		result[k] = v
	}

	return result, nil
}

// RegisterProvider는 새로운 메시지 제공자를 등록합니다
func (s *MessageStore) RegisterProvider(provider entity.MessageProvider) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	langCode := provider.GetLanguageCode()
	locale := entity.NewLocale(entity.Language(langCode), "")

	if !locale.IsValid() {
		return fmt.Errorf("지원하지 않는 언어 코드: %s", langCode)
	}

	// 메시지 맵 등록
	s.messages[langCode] = provider.GetMessages()
	return nil
}
