package entity

import (
	"fmt"
	"time"
)

// LogEntry는 하나의 로그 항목을 나타냅니다.
type LogEntry struct {
	// Level은 로그의 중요도를 나타냅니다.
	Level LogLevel
	// Message는 로그 메시지 내용입니다.
	Message string
	// Timestamp는 로그가 생성된 시간입니다.
	Timestamp time.Time
	// CallerInfo는 로그를 생성한 코드의 위치 정보입니다.
	CallerInfo string
	// Fields는 로그에 추가되는 구조화된 데이터입니다.
	Fields map[string]interface{}
}

// NewLogEntry는 새로운 LogEntry를 생성합니다.
func NewLogEntry(level LogLevel, message string, callerInfo string) *LogEntry {
	return &LogEntry{
		Level:      level,
		Message:    message,
		Timestamp:  time.Now(),
		CallerInfo: callerInfo,
		Fields:     make(map[string]interface{}),
	}
}

// WithField는 로그 엔트리에 필드를 추가합니다.
func (e *LogEntry) WithField(key string, value interface{}) *LogEntry {
	e.Fields[key] = value
	return e
}

// Format은 로그 엔트리를 문자열로 포맷팅합니다.
func (e *LogEntry) Format() string {
	base := fmt.Sprintf("[%s] %s [%s] %s",
		e.Timestamp.Format("2006-01-02 15:04:05"),
		e.Level.String(),
		e.CallerInfo,
		e.Message,
	)

	// 추가 필드가 있는 경우 포함
	if len(e.Fields) > 0 {
		fields := ""
		for k, v := range e.Fields {
			fields += fmt.Sprintf(" %s=%v", k, v)
		}
		base += fields
	}

	return base
}

// IsValid는 로그 엔트리가 유효한지 검사합니다.
func (e *LogEntry) IsValid() bool {
	return e.Level.IsValid() &&
		e.Message != "" &&
		!e.Timestamp.IsZero() &&
		e.CallerInfo != ""
}
