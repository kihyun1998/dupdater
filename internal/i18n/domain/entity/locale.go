package entity

// Language는 지원되는 언어 코드를 정의합니다
type Language string

const (
	// Korean 한국어
	Korean Language = "ko"

	// English 영어
	English Language = "en"

	// DefaultLanguage 기본 언어를 한국어로 설정
	DefaultLanguage = Korean
)

// Locale은 로케일 정보를 담는 도메인 엔티티입니다
type Locale struct {
	Language Language
	Messages map[string]string
}

// NewLocale은 새로운 Locale 인스턴스를 생성합니다
func NewLocale(lang Language, messages map[string]string) *Locale {
	if messages == nil {
		messages = make(map[string]string)
	}
	return &Locale{
		Language: lang,
		Messages: messages,
	}
}

// GetMessage는 지정된 키에 해당하는 메시지를 반환합니다
func (l *Locale) GetMessage(key string) string {
	if msg, ok := l.Messages[key]; ok {
		return msg
	}
	return key
}

// IsValid는 로케일이 유효한지 검증합니다
func (l *Locale) IsValid() bool {
	return l.Language == Korean || l.Language == English
}
