package entity

// LogLevel은 로그의 중요도를 나타냅니다
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String은 LogLevel을 문자열로 변환합니다
func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}
