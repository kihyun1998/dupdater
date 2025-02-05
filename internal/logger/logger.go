package logger

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"runtime"
// 	"sync"
// 	"time"
// )

// // LogLevel은 로그의 중요도를 나타냅니다
// type LogLevel int

// const (
// 	DEBUG LogLevel = iota
// 	INFO
// 	WARN
// 	ERROR
// 	FATAL
// )

// // 로그 레벨을 문자열로 변환
// func (l LogLevel) String() string {
// 	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
// }

// // LogEntry는 하나의 로그 항목을 나타냅니다
// type LogEntry struct {
// 	Level      LogLevel  // 로그 레벨
// 	Message    string    // 로그 메시지
// 	Timestamp  time.Time // 로그 발생 시간
// 	CallerInfo string    // 호출자 정보
// }

// // Logger는 로깅을 관리하는 구조체입니다
// type Logger struct {
// 	mu         sync.Mutex  // 동시성 제어를 위한 뮤텍스
// 	logFile    *os.File    // 로그 파일
// 	logger     *log.Logger // 실제 로깅을 수행하는 logger
// 	logLevel   LogLevel    // 현재 로그 레벨
// 	logPath    string      // 로그 파일 경로
// 	maxSize    int64       // 최대 로그 파일 크기 (바이트)
// 	maxBackups int         // 보관할 최대 백업 파일 수
// }

// // Config는 Logger 생성에 필요한 설정을 담는 구조체입니다
// type Config struct {
// 	LogPath    string   // 로그 파일 경로
// 	LogLevel   LogLevel // 로그 레벨
// 	MaxSize    int64    // 최대 파일 크기 (바이트)
// 	MaxBackups int      // 최대 백업 파일 수
// }

// // New는 새로운 Logger 인스턴스를 생성합니다
// func New(config Config) (*Logger, error) {
// 	// 기본값 설정
// 	if config.LogPath == "" {
// 		homeDir, err := os.UserHomeDir()
// 		if err != nil {
// 			return nil, fmt.Errorf("홈 디렉토리 찾기 실패: %w", err)
// 		}
// 		config.LogPath = filepath.Join(homeDir, ".testfolder", "logs", "updater.log")
// 	}

// 	if config.MaxSize == 0 {
// 		config.MaxSize = 10 * 1024 * 1024 // 기본 10MB
// 	}

// 	if config.MaxBackups == 0 {
// 		config.MaxBackups = 5 // 기본 5개 백업
// 	}

// 	// 로그 디렉토리 생성
// 	if err := os.MkdirAll(filepath.Dir(config.LogPath), 0755); err != nil {
// 		return nil, fmt.Errorf("로그 디렉토리 생성 실패: %w", err)
// 	}

// 	// 로그 파일 생성
// 	logFile, err := os.OpenFile(
// 		config.LogPath,
// 		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
// 		0644,
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("로그 파일 생성 실패: %w", err)
// 	}

// 	return &Logger{
// 		logFile:    logFile,
// 		logger:     log.New(logFile, "", 0),
// 		logLevel:   config.LogLevel,
// 		logPath:    config.LogPath,
// 		maxSize:    config.MaxSize,
// 		maxBackups: config.MaxBackups,
// 	}, nil
// }

// // log는 실제 로깅을 수행하는 내부 메서드입니다
// func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
// 	if level < l.logLevel {
// 		return
// 	}

// 	l.mu.Lock()
// 	defer l.mu.Unlock()

// 	// 로그 순환 체크
// 	if err := l.checkRotate(); err != nil {
// 		fmt.Fprintf(os.Stderr, "로그 순환 실패: %v\n", err)
// 	}

// 	// 호출자 정보 가져오기
// 	_, file, line, ok := runtime.Caller(2)
// 	callerInfo := "unknown"
// 	if ok {
// 		callerInfo = fmt.Sprintf("%s:%d", filepath.Base(file), line)
// 	}

// 	// 로그 엔트리 생성
// 	entry := LogEntry{
// 		Level:      level,
// 		Message:    fmt.Sprintf(format, args...),
// 		Timestamp:  time.Now(),
// 		CallerInfo: callerInfo,
// 	}

// 	// 로그 포맷팅 및 작성
// 	logLine := fmt.Sprintf(
// 		"[%s] %s [%s] %s",
// 		entry.Timestamp.Format("2006-01-02 15:04:05"),
// 		entry.Level,
// 		entry.CallerInfo,
// 		entry.Message,
// 	)

// 	if err := l.logger.Output(0, logLine); err != nil {
// 		fmt.Fprintf(os.Stderr, "로그 작성 실패: %v\n", err)
// 	}
// }

// // 공개 로깅 메서드들
// func (l *Logger) Debug(format string, args ...interface{}) {
// 	l.log(DEBUG, format, args...)
// }

// func (l *Logger) Info(format string, args ...interface{}) {
// 	l.log(INFO, format, args...)
// }

// func (l *Logger) Warn(format string, args ...interface{}) {
// 	l.log(WARN, format, args...)
// }

// func (l *Logger) Error(format string, args ...interface{}) {
// 	l.log(ERROR, format, args...)
// }

// func (l *Logger) Fatal(format string, args ...interface{}) {
// 	l.log(FATAL, format, args...)
// 	os.Exit(1)
// }

// // Close는 로거를 정리합니다
// func (l *Logger) Close() error {
// 	l.mu.Lock()
// 	defer l.mu.Unlock()

// 	if l.logFile != nil {
// 		if err := l.logFile.Sync(); err != nil {
// 			return fmt.Errorf("로그 파일 동기화 실패: %w", err)
// 		}
// 		if err := l.logFile.Close(); err != nil {
// 			return fmt.Errorf("로그 파일 닫기 실패: %w", err)
// 		}
// 		l.logFile = nil
// 	}
// 	return nil
// }

// // checkRotate는 로그 파일 크기를 확인하고 필요시 순환합니다
// func (l *Logger) checkRotate() error {
// 	info, err := l.logFile.Stat()
// 	if err != nil {
// 		return fmt.Errorf("파일 정보 가져오기 실패: %w", err)
// 	}

// 	if info.Size() < l.maxSize {
// 		return nil
// 	}

// 	// 현재 파일 닫기
// 	if err := l.logFile.Close(); err != nil {
// 		return fmt.Errorf("현재 로그 파일 닫기 실패: %w", err)
// 	}

// 	// 기존 백업 파일들 순환
// 	for i := l.maxBackups - 1; i >= 0; i-- {
// 		oldPath := fmt.Sprintf("%s.%d", l.logPath, i)
// 		newPath := fmt.Sprintf("%s.%d", l.logPath, i+1)

// 		if i == 0 {
// 			oldPath = l.logPath
// 		}

// 		// 마지막 백업 파일 삭제
// 		if i == l.maxBackups-1 {
// 			os.Remove(newPath)
// 			continue
// 		}

// 		// 나머지 파일들 이름 변경
// 		if _, err := os.Stat(oldPath); err == nil {
// 			if err := os.Rename(oldPath, newPath); err != nil {
// 				return fmt.Errorf("파일 이름 변경 실패: %w", err)
// 			}
// 		}
// 	}

// 	// 새 로그 파일 생성
// 	newFile, err := os.OpenFile(
// 		l.logPath,
// 		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
// 		0644,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("새 로그 파일 생성 실패: %w", err)
// 	}

// 	l.logFile = newFile
// 	l.logger = log.New(newFile, "", 0)

// 	return nil
// }
