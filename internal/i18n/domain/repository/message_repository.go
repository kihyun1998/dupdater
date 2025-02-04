package repository

import "github.com/kihyun1998/dupdater/internal/i18n/domain/entity"

// MessageRepository는 메시지 저장소 인터페이스를 정의합니다
type MessageRepository interface {
	// GetMessage는 주어진 키와 로케일에 해당하는 메시지를 반환합니다
	GetMessage(key string, locale *entity.Locale) (*entity.Message, error)

	// AddMessages는 여러 메시지를 저장소에 추가합니다
	AddMessages(messages []*entity.Message) error

	// GetAvailableLocales는 사용 가능한 모든 로케일을 반환합니다
	GetAvailableLocales() []*entity.Locale
}

// MessageProvider는 특정 언어의 메시지를 제공하는 인터페이스입니다
type MessageProvider interface {
	// GetLocale은 제공자의 로케일을 반환합니다
	GetLocale() *entity.Locale

	// GetMessages는 해당 언어의 모든 메시지를 반환합니다
	GetMessages() map[string]string
}
