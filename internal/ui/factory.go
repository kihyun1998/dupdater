// internal/ui/factory.go

package ui

import (
	"github.com/kihyun1998/dupdater/internal/i18n"
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/ui/domain/repository"
	"github.com/kihyun1998/dupdater/internal/ui/domain/usecase"
	"github.com/kihyun1998/dupdater/internal/ui/infrastructure"
	"github.com/kihyun1998/dupdater/pkg/utils/theme"
)

// Config는 UI 생성에 필요한 설정입니다
type Config struct {
	AppName     string
	TotalSteps  int
	FromVersion string
	ToVersion   string
	Logger      logger.Logger
	Theme       theme.ThemeVariant
	I18n        i18n.Manager
}

// Manager는 updater.UIManager와 동일한 인터페이스를 제공합니다
type Manager interface {
	SetCurrentStep(step int)
	UpdateDetail(message string)
	ShowError(err error)
	Run()
	Close()
	SetRestoreHandler(handler func())
	ShowRestoring()
	ShowRestoreComplete()
	GetTotalSteps() int
	SetCompletionCallback(func())
}

// manager는 UI 관리자의 실제 구현체입니다
type manager struct {
	service *usecase.UIService
}

// New는 새로운 UI Manager를 생성합니다
func New(config Config) Manager {
	// repository 계층 초기화
	repo := infrastructure.NewFyneManager(&repository.Config{
		AppName:     config.AppName,
		TotalSteps:  config.TotalSteps,
		FromVersion: config.FromVersion,
		ToVersion:   config.ToVersion,
		Logger:      config.Logger,
		Theme:       config.Theme,
		I18n:        config.I18n,
	})

	// service 계층 초기화
	service := usecase.NewUIService(repo, config.Logger)

	return &manager{
		service: service,
	}
}

// Manager 인터페이스 구현
func (m *manager) SetCurrentStep(step int) {
	m.service.SetCurrentStep(step)
}

func (m *manager) UpdateDetail(message string) {
	m.service.UpdateDetail(message)
}

func (m *manager) ShowError(err error) {
	m.service.ShowError(err)
}

func (m *manager) Run() {
	m.service.Run()
}

func (m *manager) Close() {
	m.service.Close()
}

func (m *manager) SetRestoreHandler(handler func()) {
	m.service.SetRestoreHandler(handler)
}

func (m *manager) ShowRestoring() {
	m.service.ShowRestoring()
}

func (m *manager) ShowRestoreComplete() {
	m.service.ShowRestoreComplete()
}

func (m *manager) GetTotalSteps() int {
	return m.service.GetTotalSteps()
}

func (m *manager) SetCompletionCallback(callback func()) {
	m.service.SetCompletionCallback(callback)
}
