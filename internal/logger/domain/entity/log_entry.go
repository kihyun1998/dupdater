package entity

import (
	"fmt"
	"time"
)

// LogEntry는 하나의 로그 항목을 나타냅니다
type LogEntry struct {
	Level      LogLevel
	Message    string
	Timestamp  time.Time
	CallerInfo string
	Module     string
}

// NewLogEntry는 새로운 LogEntry를 생성합니다
func NewLogEntry(level LogLevel, message string, callerInfo string, module string) *LogEntry {
	return &LogEntry{
		Level:      level,
		Message:    message,
		Timestamp:  time.Now(),
		CallerInfo: callerInfo,
		Module:     module,
	}
}

// FormatMessage는 로그 메시지를 포맷팅합니다
func (e *LogEntry) FormatMessage() string {
	return fmt.Sprintf("[%s] %s [%s] %s: %s",
		e.Timestamp.Format("2006-01-02 15:04:05"),
		e.Level,
		e.Module,
		e.CallerInfo,
		e.Message,
	)
}
