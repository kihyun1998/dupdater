package repository

import "github.com/kihyun1998/dupdater/internal/logger/domain/entity"

// LoggerRepository는 로그 저장소 인터페이스를 정의합니다
type LoggerRepository interface {
	// Log는 로그 엔트리를 저장합니다
	Log(entry *entity.LogEntry) error

	// Close는 로그 저장소 연결을 종료합니다
	Close() error

	// Rotate는 로그 파일을 순환합니다
	Rotate() error
}

// LogConfig는 로거 설정을 정의합니다
type LogConfig struct {
	LogPath    string
	MaxSize    int64
	MaxBackups int
	MinLevel   entity.LogLevel
}
