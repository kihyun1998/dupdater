package i18n

// import (
// 	"fmt"
// 	"sync"

// 	"github.com/kihyun1998/dupdater/internal/i18n/locale"
// )

// // Manager는 다국어 지원을 관리하는 구조체입니다
// type Manager struct {
// 	currentLang string
// 	messages    map[string]map[string]string
// 	mutex       sync.RWMutex
// 	providers   map[string]MessageProvider
// }

// // Config는 Manager 생성에 필요한 설정을 담는 구조체입니다
// type Config struct {
// 	DefaultLang string
// }

// // New는 새로운 Manager 인스턴스를 생성합니다
// func New(config Config) (*Manager, error) {
// 	if config.DefaultLang == "" {
// 		config.DefaultLang = string(DefaultLanguage)
// 	}

// 	m := &Manager{
// 		messages:  make(map[string]map[string]string),
// 		providers: make(map[string]MessageProvider),
// 	}

// 	// 기본 제공자 등록
// 	m.RegisterProvider(&locale.KoreanProvider{})
// 	m.RegisterProvider(&locale.EnglishProvider{})

// 	// 기본 언어 설정
// 	if err := m.SetLanguage(config.DefaultLang); err != nil {
// 		return nil, err
// 	}

// 	return m, nil
// }

// // RegisterProvider는 새로운 메시지 제공자를 등록합니다
// func (m *Manager) RegisterProvider(provider MessageProvider) {
// 	m.mutex.Lock()
// 	defer m.mutex.Unlock()

// 	langCode := provider.GetLanguageCode()
// 	m.providers[langCode] = provider
// 	m.messages[langCode] = provider.GetMessages()
// }

// // SetLanguage는 현재 언어를 설정합니다
// func (m *Manager) SetLanguage(lang string) error {
// 	if !IsValidLanguage(lang) {
// 		return fmt.Errorf("지원하지 않는 언어입니다: %s", lang)
// 	}

// 	m.mutex.Lock()
// 	defer m.mutex.Unlock()

// 	if _, exists := m.messages[lang]; !exists {
// 		return fmt.Errorf("언어 리소스를 찾을 수 없습니다: %s", lang)
// 	}

// 	m.currentLang = lang
// 	return nil
// }

// // GetMessage는 지정된 키에 해당하는 메시지를 현재 설정된 언어로 반환합니다
// func (m *Manager) GetMessage(key string) string {
// 	m.mutex.RLock()
// 	defer m.mutex.RUnlock()

// 	if messages, exists := m.messages[m.currentLang]; exists {
// 		if msg, ok := messages[key]; ok {
// 			return msg
// 		}
// 	}

// 	// 메시지를 찾을 수 없는 경우 키를 반환
// 	return key
// }

// // GetCurrentLanguage는 현재 설정된 언어를 반환합니다
// func (m *Manager) GetCurrentLanguage() string {
// 	m.mutex.RLock()
// 	defer m.mutex.RUnlock()

// 	return m.currentLang
// }

// // formatMessage는 메시지에 인자를 적용합니다
// func (m *Manager) formatMessage(message string, args ...interface{}) string {
// 	if len(args) > 0 {
// 		return fmt.Sprintf(message, args...)
// 	}
// 	return message
// }
