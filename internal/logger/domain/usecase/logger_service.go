package usecase

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/kihyun1998/dupdater/internal/logger/domain/entity"
	"github.com/kihyun1998/dupdater/internal/logger/domain/repository"
)

// LoggerService는 로깅 서비스를 정의합니다
type LoggerService struct {
	repo     repository.LoggerRepository
	minLevel entity.LogLevel
}

// NewLoggerService는 새로운 LoggerService를 생성합니다
func NewLoggerService(repo repository.LoggerRepository, config repository.LogConfig) *LoggerService {
	return &LoggerService{
		repo:     repo,
		minLevel: config.MinLevel,
	}
}

// log는 실제 로깅을 수행하는 내부 메서드입니다
func (s *LoggerService) log(level entity.LogLevel, format string, args ...interface{}) error {
	if level < s.minLevel {
		return nil
	}

	// 호출자 정보 가져오기
	_, file, line, ok := runtime.Caller(2)
	callerInfo := "unknown"
	if ok {
		callerInfo = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	// 로그 엔트리 생성
	entry := entity.NewLogEntry(
		level,
		fmt.Sprintf(format, args...),
		callerInfo,
		filepath.Base(filepath.Dir(file)),
	)

	return s.repo.Log(entry)
}

// 공개 메서드들
func (s *LoggerService) Debug(format string, args ...interface{}) error {
	return s.log(entity.DEBUG, format, args...)
}

func (s *LoggerService) Info(format string, args ...interface{}) error {
	return s.log(entity.INFO, format, args...)
}

func (s *LoggerService) Warn(format string, args ...interface{}) error {
	return s.log(entity.WARN, format, args...)
}

func (s *LoggerService) Error(format string, args ...interface{}) error {
	return s.log(entity.ERROR, format, args...)
}

func (s *LoggerService) Fatal(format string, args ...interface{}) error {
	return s.log(entity.FATAL, format, args...)
}

// Close는 로거를 정리합니다
func (s *LoggerService) Close() error {
	return s.repo.Close()
}
