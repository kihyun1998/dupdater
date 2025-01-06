package config

import (
	"os"
	"path/filepath"
)

// Config는 전체 어플리케이션 설정을 담는 구조체
type Config struct {
	// 어플리케이션 기본 설정
	AppName     string `json:"appName"`     // 업데이트할 어플리케이션 이름
	FromVersion string `json:"fromVersion"` // 업데이트 시작 버전
	ServerName  string `json:"serverName"`  // 서버 프로필명

	// 경로 설정
	Paths *Paths `json:"paths"`

	// 로그 설정
	LogConfig *LogConfig `json:"logConfig"`
}

// LogConfig는 로깅 관련 설정을 담는 구조체
type LogConfig struct {
	LogDir     string `json:"logDir"`     // 로그 저장 디렉토리
	MaxSize    int    `json:"maxSize"`    // 최대 로그 파일 크기 (MB)
	MaxBackups int    `json:"maxBackups"` // 보관할 최대 로그 파일 수
}

// NewConfig는 새로운 설정 인스턴스를 생성
func NewConfig() *Config {
	return &Config{
		AppName: "update_test_app_1.exe", // 기본값 설정
		Paths:   NewPaths(),
		LogConfig: &LogConfig{
			LogDir:     filepath.Join(os.Getenv("APPDATA"), "acrapoint", "logs"),
			MaxSize:    10,
			MaxBackups: 5,
		},
	}
}

// LoadFromArgs는 커맨드라인 인자에서 설정을 로드
func (c *Config) LoadFromArgs(fromVersion, serverName string) error {
	if fromVersion != "" {
		c.FromVersion = fromVersion
	}
	if serverName != "" {
		c.ServerName = serverName
	}
	return nil
}

// Validate는 필수 설정값들을 검증
func (c *Config) Validate() error {
	if c.FromVersion == "" {
		return ErrMissingFromVersion
	}
	if c.ServerName == "" {
		return ErrMissingServerName
	}
	return nil
}
