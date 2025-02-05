package repository

import (
	"path/filepath"

	"github.com/kihyun1998/dupdater/internal/logger/domain/entity"
)

// LogRepository는 로그 저장소의 인터페이스를 정의합니다.
type LogRepository interface {
	// Write는 로그 엔트리를 저장합니다.
	Write(entry *entity.LogEntry) error

	// Rotate는 로그 파일을 순환합니다.
	Rotate() error

	// Close는 로그 저장소를 정리합니다.
	Close() error
}

// LogConfig는 로그 저장소의 설정을 정의합니다.
type LogConfig struct {
	// LogPath는 로그 파일의 경로입니다.
	LogPath string
	// MaxSize는 로그 파일의 최대 크기(바이트)입니다.
	MaxSize int64
	// MaxBackups는 보관할 최대 백업 파일 수입니다.
	MaxBackups int
}

// Validate는 로그 설정이 유효한지 검사합니다.
func (c *LogConfig) Validate() error {
	if c.LogPath == "" {
		return filepath.ErrBadPattern
	}
	if c.MaxSize <= 0 {
		c.MaxSize = 10 * 1024 * 1024 // 기본값 10MB
	}
	if c.MaxBackups <= 0 {
		c.MaxBackups = 5 // 기본값 5개
	}
	return nil
}

// NewLogConfig는 새로운 LogConfig를 생성합니다.
func NewLogConfig(path string, maxSize int64, maxBackups int) *LogConfig {
	return &LogConfig{
		LogPath:    path,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
	}
}
