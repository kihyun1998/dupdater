package entity

// MessageProvider는 특정 언어의 메시지를 제공하는 인터페이스입니다
type MessageProvider interface {
	// GetLanguageCode는 제공자가 지원하는 언어 코드를 반환합니다
	GetLanguage() Language

	// GetMessages는 해당 언어의 전체 메시지 맵을 반환합니다
	GetMessages() map[string]string
}

// Message는 다국어 메시지를 나타내는 값 객체입니다
type Message struct {
	Key   string
	Value string
	Lang  Language
}

// NewMessage는 새로운 Message 인스턴스를 생성합니다
func NewMessage(key, value string, lang Language) *Message {
	return &Message{
		Key:   key,
		Value: value,
		Lang:  lang,
	}
}

// String은 메시지의 문자열 표현을 반환합니다
func (m *Message) String() string {
	return m.Value
}
