// Package ports는 i18n 도메인의 외부 인터페이스를 정의합니다
package ports

import "github.com/kihyun1998/dupdater/internal/i18n/domain/entity"

// LocalePort는 i18n 시스템의 외부 인터페이스를 정의합니다
type LocalePort interface {
	// SetLanguage는 현재 언어를 설정합니다
	SetLanguage(lang string) error

	// GetMessage는 지정된 키에 해당하는 메시지를 현재 설정된 언어로 반환합니다
	GetMessage(key string) string

	// GetCurrentLanguage는 현재 설정된 언어를 반환합니다
	GetCurrentLanguage() string

	// RegisterProvider는 새로운 메시지 제공자를 등록합니다
	RegisterProvider(provider entity.MessageProvider) error
}

// LocaleConfig는 로케일 초기화에 필요한 설정을 정의합니다
type LocaleConfig struct {
	// DefaultLang은 기본 언어 설정입니다
	DefaultLang string
}

// Validate는 설정이 유효한지 검증합니다
func (c *LocaleConfig) Validate() error {
	if c.DefaultLang == "" {
		c.DefaultLang = string(entity.DefaultLanguage)
	}
	return nil
}
