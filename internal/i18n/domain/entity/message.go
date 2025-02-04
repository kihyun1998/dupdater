package entity

import "fmt"

// Message는 다국어 메시지를 나타내는 엔티티입니다
type Message struct {
	Key     string
	Content string
	Locale  *Locale
}

// NewMessage는 새로운 Message를 생성합니다
func NewMessage(key string, content string, locale *Locale) *Message {
	return &Message{
		Key:     key,
		Content: content,
		Locale:  locale,
	}
}

// Format은 메시지를 주어진 인자로 포맷팅합니다
func (m *Message) Format(args ...interface{}) string {
	if len(args) > 0 {
		return fmt.Sprintf(m.Content, args...)
	}
	return m.Content
}
