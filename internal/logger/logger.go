package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String은 LogLevel의 문자열 표현을 반환
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger는 로깅을 담당하는 구조체
type Logger struct {
	logFile    *os.File
	logger     *log.Logger
	logDir     string
	maxSize    int64 // 바이트 단위
	maxBackups int
}

// LoggerConfig는 Logger 설정을 위한 구조체
type LoggerConfig struct {
	LogDir     string // 로그 저장 디렉토리
	MaxSize    int    // MB 단위
	MaxBackups int    // 보관할 최대 로그 파일 수
}

// NewLogger는 새로운 Logger 인스턴스를 생성
func NewLogger(config LoggerConfig) (*Logger, error) {
	// 로그 디렉토리 생성
	if err := os.MkdirAll(config.LogDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("로그 디렉토리 생성 실패: %w", err)
	}

	// 로그 파일명 생성 (날짜_시간 기준)
	logFileName := fmt.Sprintf("update_%s.log", time.Now().Format("2006-01-02_15-04-05"))
	logFilePath := filepath.Join(config.LogDir, logFileName)

	// 로그 파일 생성
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("로그 파일 생성 실패: %w", err)
	}

	return &Logger{
		logFile:    file,
		logger:     log.New(file, "", log.Ldate|log.Ltime),
		logDir:     config.LogDir,
		maxSize:    int64(config.MaxSize) * 1024 * 1024, // MB를 바이트로 변환
		maxBackups: config.MaxBackups,
	}, nil
}

// log는 실제 로깅을 수행하는 내부 메서드
func (l *Logger) log(level LogLevel, format string, v ...interface{}) {
	// 호출자 정보 획득
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}
	// 파일명만 추출
	file = filepath.Base(file)

	// 로그 메시지 포맷팅
	msg := fmt.Sprintf(format, v...)
	logMsg := fmt.Sprintf("[%s] [%s:%d] %s", level, file, line, msg)

	// 로그 작성
	l.logger.Println(logMsg)

	// 파일 크기 체크 및 로테이션
	l.checkRotate()
}

// Debug 레벨 로그
func (l *Logger) Debug(format string, v ...interface{}) {
	l.log(DEBUG, format, v...)
}

// Info 레벨 로그
func (l *Logger) Info(format string, v ...interface{}) {
	l.log(INFO, format, v...)
}

// Warn 레벨 로그
func (l *Logger) Warn(format string, v ...interface{}) {
	l.log(WARN, format, v...)
}

// Error 레벨 로그
func (l *Logger) Error(format string, v ...interface{}) {
	l.log(ERROR, format, v...)
}

// Close는 로거 리소스를 정리
func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// checkRotate는 로그 파일 크기를 체크하고 필요시 로테이션 수행
func (l *Logger) checkRotate() {
	if l.logFile == nil {
		return
	}

	// 현재 파일 크기 확인
	info, err := l.logFile.Stat()
	if err != nil {
		return
	}

	// 최대 크기를 초과하면 로테이션
	if info.Size() > l.maxSize {
		l.rotate()
	}
}

// rotate는 로그 파일 로테이션을 수행
func (l *Logger) rotate() {
	// 현재 로그 파일 닫기
	l.logFile.Close()

	// 기존 로그 파일들의 이름 변경
	pattern := filepath.Join(l.logDir, "update_*.log")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	// 오래된 로그 파일 삭제
	if len(matches) >= l.maxBackups {
		// 날짜순으로 정렬
		// (구현 생략 - 실제로는 파일명에서 날짜를 추출하여 정렬 필요)
		// 가장 오래된 파일 삭제
		os.Remove(matches[0])
	}

	// 새 로그 파일 생성
	logFileName := fmt.Sprintf("update_%s.log", time.Now().Format("2006-01-02_15-04-05"))
	logFilePath := filepath.Join(l.logDir, logFileName)
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}

	// 로거 업데이트
	l.logFile = file
	l.logger = log.New(file, "", log.Ldate|log.Ltime)
}
