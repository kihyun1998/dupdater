package entity

// Language는 지원되는 언어 코드를 정의합니다
type Language string

const (
	Korean  Language = "ko"
	English Language = "en"
)

// DefaultLanguage는 기본 언어를 정의합니다
const DefaultLanguage = Korean

// Locale은 지역화 설정을 나타내는 엔티티입니다
type Locale struct {
	Code Language
}

// NewLocale은 새로운 Locale을 생성합니다
func NewLocale(code Language) *Locale {
	return &Locale{
		Code: code,
	}
}

// IsValid는 언어 코드가 유효한지 검사합니다
func (l *Locale) IsValid() bool {
	switch l.Code {
	case Korean, English:
		return true
	default:
		return false
	}
}

// String은 언어 코드를 문자열로 반환합니다
func (l *Locale) String() string {
	return string(l.Code)
}
