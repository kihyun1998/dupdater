// factory.go
package i18n

import (
	"fmt"

	"github.com/kihyun1998/dupdater/internal/i18n/domain/entity"
	"github.com/kihyun1998/dupdater/internal/i18n/domain/repository"
	"github.com/kihyun1998/dupdater/internal/i18n/domain/usecase"
	"github.com/kihyun1998/dupdater/internal/i18n/infrastructure"
	"github.com/kihyun1998/dupdater/internal/i18n/infrastructure/providers"
)

// Manager는 다국어 지원을 위한 인터페이스입니다
type Manager interface {
	GetMessage(key string) string
	SetLanguage(lang string) error
	GetCurrentLanguage() string
}

// Config는 Manager 생성에 필요한 설정입니다
type Config struct {
	DefaultLanguage string // 기본 언어 설정
	Logger          repository.Logger
}

// managerImpl은 Manager 인터페이스 구현체입니다
type managerImpl struct {
	service *usecase.I18nService
}

// New는 새로운 Manager 인스턴스를 생성합니다
func New(config Config) (Manager, error) {
	// 설정 검증
	if config.DefaultLanguage == "" {
		config.DefaultLanguage = string(entity.DefaultLanguage)
	}

	// 저장소 설정
	repoConfig := &repository.Config{
		DefaultLanguage: entity.Language(config.DefaultLanguage),
		Logger:          config.Logger,
	}

	// 저장소 생성
	store := infrastructure.NewI18nStore(repoConfig)

	// 서비스 생성
	service := usecase.NewI18nService(store, config.Logger)

	// 기본 메시지 제공자 등록
	defaultProviders := []entity.MessageProvider{
		providers.NewKoreanProvider(),
		providers.NewEnglishProvider(),
	}

	for _, provider := range defaultProviders {
		if err := service.RegisterProvider(provider); err != nil {
			return nil, fmt.Errorf("메시지 제공자 등록 실패: %w", err)
		}
	}

	// 초기 언어 설정
	if err := service.SetLanguage(config.DefaultLanguage); err != nil {
		return nil, fmt.Errorf("초기 언어 설정 실패: %w", err)
	}

	return &managerImpl{
		service: service,
	}, nil
}

// Manager 인터페이스 구현
func (m *managerImpl) GetMessage(key string) string {
	return m.service.GetMessage(key)
}

func (m *managerImpl) SetLanguage(lang string) error {
	return m.service.SetLanguage(lang)
}

func (m *managerImpl) GetCurrentLanguage() string {
	return m.service.GetCurrentLanguage()
}
