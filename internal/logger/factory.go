package logger

import (
	"github.com/kihyun1998/dupdater/internal/logger/domain/entity"
	"github.com/kihyun1998/dupdater/internal/logger/domain/repository"
	"github.com/kihyun1998/dupdater/internal/logger/domain/usecase"
	"github.com/kihyun1998/dupdater/internal/logger/infrastructure"
)

// Config는 로거 생성에 필요한 설정을 정의합니다.
type Config struct {
	LogPath    string          // 로그 파일 경로
	LogLevel   entity.LogLevel // 로그 레벨
	MaxSize    int64           // 최대 파일 크기 (바이트)
	MaxBackups int             // 최대 백업 파일 수
}

// Logger는 로깅을 위한 인터페이스입니다.
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Fatal(format string, v ...interface{})
	Close() error
}

// New는 새로운 로거 인스턴스를 생성합니다.
func New(config Config) (Logger, error) {
	// 저장소 설정 생성
	repoConfig := repository.NewLogConfig(
		config.LogPath,
		config.MaxSize,
		config.MaxBackups,
	)

	// 파일 로거 생성
	fileLogger, err := infrastructure.NewFileLogger(repoConfig)
	if err != nil {
		return nil, err
	}

	// 로깅 서비스 생성 및 반환
	return usecase.NewLoggerService(fileLogger, config.LogLevel), nil
}
