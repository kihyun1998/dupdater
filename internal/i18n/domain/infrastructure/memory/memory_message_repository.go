package memory

import (
	"fmt"
	"sync"

	"github.com/kihyun1998/dupdater/internal/i18n/domain/entity"
)

// MemoryMessageRepository는 메모리 기반 메시지 저장소를 구현합니다
type MemoryMessageRepository struct {
	messages map[string]map[entity.Language]*entity.Message
	mu       sync.RWMutex
}

// NewMemoryMessageRepository는 새로운 MemoryMessageRepository를 생성합니다
func NewMemoryMessageRepository() *MemoryMessageRepository {
	return &MemoryMessageRepository{
		messages: make(map[string]map[entity.Language]*entity.Message),
	}
}

// GetMessage는 주어진 키와 로케일에 해당하는 메시지를 반환합니다
func (r *MemoryMessageRepository) GetMessage(key string, locale *entity.Locale) (*entity.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if localeMessages, exists := r.messages[key]; exists {
		if msg, ok := localeMessages[locale.Code]; ok {
			return msg, nil
		}
	}
	return nil, fmt.Errorf("메시지를 찾을 수 없음: %s (%s)", key, locale.String())
}

// AddMessages는 여러 메시지를 저장소에 추가합니다
func (r *MemoryMessageRepository) AddMessages(messages []*entity.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, msg := range messages {
		if _, exists := r.messages[msg.Key]; !exists {
			r.messages[msg.Key] = make(map[entity.Language]*entity.Message)
		}
		r.messages[msg.Key][msg.Locale.Code] = msg
	}
	return nil
}

// GetAvailableLocales는 사용 가능한 모든 로케일을 반환합니다
func (r *MemoryMessageRepository) GetAvailableLocales() []*entity.Locale {
	r.mu.RLock()
	defer r.mu.RUnlock()

	localeMap := make(map[entity.Language]bool)
	for _, localeMessages := range r.messages {
		for lang := range localeMessages {
			localeMap[lang] = true
		}
	}

	locales := make([]*entity.Locale, 0, len(localeMap))
	for lang := range localeMap {
		locales = append(locales, entity.NewLocale(lang))
	}
	return locales
}
