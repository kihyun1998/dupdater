package repository

import "github.com/kihyun1998/dupdater/internal/i18n/domain/entity"

// I18nRepository는 다국어 지원을 위한 저장소 인터페이스입니다
type I18nRepository interface {
	// GetMessage는 현재 언어의 메시지를 조회합니다
	GetMessage(key string) string

	// SetLanguage는 현재 언어를 설정합니다
	SetLanguage(lang entity.Language) error

	// GetCurrentLanguage는 현재 설정된 언어를 반환합니다
	GetCurrentLanguage() entity.Language

	// RegisterProvider는 새로운 메시지 제공자를 등록합니다
	RegisterProvider(provider entity.MessageProvider) error
}

// Config는 저장소 설정을 정의합니다
type Config struct {
	DefaultLanguage entity.Language
	Logger          Logger
}

// Logger는 로깅을 위한 인터페이스입니다
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}
