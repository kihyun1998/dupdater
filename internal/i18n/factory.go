// Package i18n은 다국어 지원 시스템의 진입점을 제공합니다
package i18n

import (
	"fmt"

	"github.com/kihyun1998/dupdater/internal/i18n/domain/entity"
	"github.com/kihyun1998/dupdater/internal/i18n/domain/ports"
	"github.com/kihyun1998/dupdater/internal/i18n/domain/repository"
	"github.com/kihyun1998/dupdater/internal/i18n/domain/usecase"
	"github.com/kihyun1998/dupdater/internal/i18n/infrastructure"
	"github.com/kihyun1998/dupdater/internal/i18n/infrastructure/providers"
)

// Factory는 i18n 시스템의 컴포넌트들을 생성하고 관리합니다
type Factory struct {
	config *ports.LocaleConfig
}

// NewFactory는 새로운 Factory 인스턴스를 생성합니다
func NewFactory(config *ports.LocaleConfig) (*Factory, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("설정 검증 실패: %w", err)
	}

	return &Factory{
		config: config,
	}, nil
}

// Create는 i18n 시스템의 인스턴스를 생성하고 초기화합니다
func (f *Factory) Create() (ports.LocalePort, error) {
	// 1. 저장소 설정 생성
	repoConfig := &repository.Config{
		DefaultLocale: entity.NewLocale(entity.Language(f.config.DefaultLang), "Default"),
	}

	// 2. 메시지 저장소 생성
	messageStore, err := infrastructure.NewMessageStore(repoConfig)
	if err != nil {
		return nil, fmt.Errorf("메시지 저장소 생성 실패: %w", err)
	}

	// 3. 기본 메시지 제공자 등록
	providers := []entity.MessageProvider{
		providers.NewKoreanProvider(),
		providers.NewEnglishProvider(),
	}

	for _, provider := range providers {
		if err := messageStore.RegisterProvider(provider); err != nil {
			return nil, fmt.Errorf("메시지 제공자 등록 실패: %w", err)
		}
	}

	// 4. 로케일 서비스 생성
	localeService := usecase.NewLocaleService(messageStore, f.config.DefaultLang)

	// 5. 초기 언어 설정
	if err := localeService.SetLanguage(f.config.DefaultLang); err != nil {
		return nil, fmt.Errorf("초기 언어 설정 실패: %w", err)
	}

	return localeService, nil
}

// New는 i18n 시스템의 새 인스턴스를 생성하는 편의 함수입니다
func New(config ports.LocaleConfig) (ports.LocalePort, error) {
	factory, err := NewFactory(&config)
	if err != nil {
		return nil, err
	}
	return factory.Create()
}
