package usecase

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/kihyun1998/dupdater/internal/logger/domain/entity"
	"github.com/kihyun1998/dupdater/internal/logger/domain/ports"
	"github.com/kihyun1998/dupdater/internal/logger/domain/repository"
)

// loggerService는 LoggerPort의 구현체입니다.
type loggerService struct {
	mu         sync.RWMutex
	repository repository.LogRepository
	level      entity.LogLevel
}

// NewLoggerService는 새로운 LoggerService 인스턴스를 생성합니다.
func NewLoggerService(repo repository.LogRepository, level entity.LogLevel) ports.LoggerPort {
	return &loggerService{
		repository: repo,
		level:      level,
	}
}

// log는 실제 로깅을 수행하는 내부 메서드입니다.
func (l *loggerService) log(level entity.LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// 호출자 정보 가져오기
	_, file, line, ok := runtime.Caller(2)
	callerInfo := "unknown"
	if ok {
		callerInfo = fmt.Sprintf("%s:%d", file, line)
	}

	// 로그 엔트리 생성
	entry := entity.NewLogEntry(
		level,
		fmt.Sprintf(format, args...),
		callerInfo,
	)

	// 로그 저장
	if err := l.repository.Write(entry); err != nil {
		fmt.Printf("로그 저장 실패: %v\n", err)
	}

	// FATAL 레벨인 경우 프로그램 종료
	if level == entity.FATAL {
		l.Close()
		os.Exit(1)
	}
}

// LoggerPort 인터페이스 구현
func (l *loggerService) Debug(format string, args ...interface{}) {
	l.log(entity.DEBUG, format, args...)
}

func (l *loggerService) Info(format string, args ...interface{}) {
	l.log(entity.INFO, format, args...)
}

func (l *loggerService) Warn(format string, args ...interface{}) {
	l.log(entity.WARN, format, args...)
}

func (l *loggerService) Error(format string, args ...interface{}) {
	l.log(entity.ERROR, format, args...)
}

func (l *loggerService) Fatal(format string, args ...interface{}) {
	l.log(entity.FATAL, format, args...)
}

func (l *loggerService) Close() error {
	return l.repository.Close()
}
