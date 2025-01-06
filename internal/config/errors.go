package config

import "errors"

var (
	// 설정 관련 에러
	ErrMissingFromVersion = errors.New("fromVersion flag is required")
	ErrMissingServerName  = errors.New("serverName flag is required")

	// 경로 관련 에러
	ErrCreateDirectory   = errors.New("failed to create directory")
	ErrInvalidConfigPath = errors.New("invalid config path")
)
