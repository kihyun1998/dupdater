package entity

// Message는 다국어 메시지를 나타내는 구조체입니다
type Message struct {
	// Key는 메시지의 고유 식별자입니다
	Key string

	// Value는 해당 언어로 번역된 메시지 내용입니다
	Value string

	// Locale은 메시지가 속한 로케일 정보입니다
	Locale *Locale
}

// NewMessage는 새로운 Message 인스턴스를 생성합니다
func NewMessage(key string, value string, locale *Locale) *Message {
	return &Message{
		Key:    key,
		Value:  value,
		Locale: locale,
	}
}

// IsValid는 메시지가 유효한지 검증합니다
func (m *Message) IsValid() bool {
	return m.Key != "" &&
		m.Value != "" &&
		m.Locale != nil &&
		m.Locale.IsValid()
}

// MessageMap은 메시지 키-값 쌍의 맵을 나타냅니다
type MessageMap map[string]string

// MessageProvider는 특정 로케일의 메시지를 제공하는 인터페이스입니다
type MessageProvider interface {
	// GetLanguageCode는 제공자가 지원하는 언어 코드를 반환합니다
	GetLanguageCode() string

	// GetMessages는 해당 언어의 전체 메시지 맵을 반환합니다
	GetMessages() MessageMap
}
