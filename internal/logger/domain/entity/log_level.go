package entity

// LogLevel은 로그의 중요도를 나타냅니다.
type LogLevel int

const (
	// DEBUG는 디버깅 목적의 상세 정보를 나타냅니다.
	DEBUG LogLevel = iota
	// INFO는 일반적인 정보성 메시지를 나타냅니다.
	INFO
	// WARN은 잠재적인 문제를 나타냅니다.
	WARN
	// ERROR는 오류 상황을 나타냅니다.
	ERROR
	// FATAL은 애플리케이션을 중단시킬 수 있는 심각한 오류를 나타냅니다.
	FATAL
)

// String은 로그 레벨을 문자열로 변환합니다.
func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}

// IsValid는 로그 레벨이 유효한지 검사합니다.
func (l LogLevel) IsValid() bool {
	return l >= DEBUG && l <= FATAL
}
