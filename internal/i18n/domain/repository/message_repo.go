package repository

import (
	"fmt"

	"github.com/kihyun1998/dupdater/internal/i18n/domain/entity"
)

// MessageRepository는 메시지 저장소의 인터페이스를 정의합니다
type MessageRepository interface {
	// GetMessage는 주어진 키와 로케일에 해당하는 메시지를 조회합니다
	GetMessage(key string, locale *entity.Locale) (*entity.Message, error)

	// GetAllMessages는 특정 로케일의 모든 메시지를 조회합니다
	GetAllMessages(locale *entity.Locale) (entity.MessageMap, error)

	// RegisterProvider는 새로운 메시지 제공자를 등록합니다
	RegisterProvider(provider entity.MessageProvider) error
}

// Config는 저장소 설정을 정의합니다
type Config struct {
	// DefaultLocale은 기본 로케일 설정입니다
	DefaultLocale *entity.Locale
}

// NewConfig는 새로운 저장소 설정을 생성합니다
func NewConfig(defaultLocale *entity.Locale) *Config {
	return &Config{
		DefaultLocale: defaultLocale,
	}
}

// Validate는 설정이 유효한지 검증합니다
func (c *Config) Validate() error {
	if c.DefaultLocale == nil || !c.DefaultLocale.IsValid() {
		return fmt.Errorf("유효하지 않은 기본 로케일 설정")
	}
	return nil
}
