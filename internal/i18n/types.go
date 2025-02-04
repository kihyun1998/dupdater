// Package i18n은 다국어 지원을 위한 패키지입니다
package i18n

// LocaleManager는 다국어 지원을 위한 인터페이스입니다
type LocaleManager interface {
	// SetLanguage는 현재 언어를 설정합니다
	SetLanguage(lang string) error

	// GetMessage는 지정된 키에 해당하는 메시지를 현재 설정된 언어로 반환합니다
	GetMessage(key string) string

	// GetCurrentLanguage는 현재 설정된 언어를 반환합니다
	GetCurrentLanguage() string
}

// MessageProvider는 각 언어별 메시지를 제공하는 인터페이스입니다
type MessageProvider interface {
	// GetMessages는 해당 언어의 모든 메시지를 반환합니다
	GetMessages() map[string]string

	// GetLanguageCode는 해당 언어의 코드를 반환합니다 (예: "ko", "en")
	GetLanguageCode() string
}

// Language는 지원되는 언어 코드를 정의합니다
type Language string

const (
	// Korean 한국어
	Korean Language = "ko"

	// English 영어
	English Language = "en"

	// DefaultLanguage 기본 언어
	DefaultLanguage = Korean
)

// IsValidLanguage는 주어진 언어 코드가 유효한지 검사합니다
func IsValidLanguage(lang string) bool {
	switch Language(lang) {
	case Korean, English:
		return true
	default:
		return false
	}
}
