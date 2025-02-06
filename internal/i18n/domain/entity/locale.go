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

// Locale은 로케일 정보를 나타내는 구조체입니다
type Locale struct {
	// Code는 언어 코드입니다 (예: "ko", "en")
	Code Language

	// Name은 해당 언어의 표시 이름입니다
	Name string
}

// NewLocale은 새로운 Locale 인스턴스를 생성합니다
func NewLocale(code Language, name string) *Locale {
	return &Locale{
		Code: code,
		Name: name,
	}
}

// IsValid는 로케일이 유효한지 검증합니다
func (l *Locale) IsValid() bool {
	switch l.Code {
	case Korean, English:
		return true
	default:
		return false
	}
}

// String은 로케일의 문자열 표현을 반환합니다
func (l *Locale) String() string {
	return string(l.Code)
}
