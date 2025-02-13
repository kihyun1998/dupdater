package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/kihyun1998/dupdater/internal/i18n"
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/ui/components"
	"github.com/kihyun1998/dupdater/pkg/utils/theme"
)

// State는 UI의 현재 상태를 나타내는 구조체입니다
type State struct {
	CurrentStep int     // 현재 진행 단계
	TotalSteps  int     // 전체 단계 수
	Progress    float64 // 진행률 (0-100)
	Detail      string  // 상세 메시지
}

// Manager는 UI를 관리하는 구조체입니다
type Manager struct {
	app        fyne.App
	mainWindow fyne.Window
	logger     logger.Logger
	i18n       i18n.Manager

	totalSteps   int
	currentStep  int
	currentTheme theme.ThemeVariant
	statusCard   *components.StatusCard
	onRestore    func()

	completionCallback func()
	animationComplete  bool
}

// Config는 Manager 생성에 필요한 설정을 담는 구조체입니다
type Config struct {
	AppName     string
	TotalSteps  int
	FromVersion string
	ToVersion   string
	Logger      logger.Logger
	Theme       theme.ThemeVariant
	I18n        i18n.Manager
}

// New는 새로운 Manager 인스턴스를 생성합니다
func New(config Config) *Manager {
	// 앱 생성 및 테마 설정
	fyneApp := app.New()
	customTheme := theme.NewCustomTheme(config.Theme)
	fyneApp.Settings().SetTheme(customTheme)

	manager := &Manager{
		app:               fyneApp,
		logger:            config.Logger,
		i18n:              config.I18n,
		totalSteps:        config.TotalSteps,
		currentStep:       0,
		animationComplete: false,
		currentTheme:      config.Theme,
	}

	// 메인 윈도우 생성
	manager.mainWindow = manager.app.NewWindow(
		manager.i18n.GetMessage("update.title"),
	)
	manager.initializeUI(config)

	return manager
}

// initializeUI는 UI 컴포넌트들을 초기화하고 배치합니다
func (m *Manager) initializeUI(config Config) {
	// StatusCard 생성 시 테마 전달
	m.statusCard = components.NewStatusCard(
		config.FromVersion,
		config.ToVersion,
		m.currentTheme,
		m.i18n,
	)

	// 이미 설정된 completion callback이 있다면 설정
	if m.completionCallback != nil {
		m.statusCard.SetCompletionCallback(func() {
			if !m.animationComplete {
				m.animationComplete = true
				m.completionCallback()
			}
		})
	}

	// 배경색 설정
	content := container.NewPadded(m.statusCard)
	content.Resize(fyne.NewSize(400, 150))

	m.mainWindow.SetContent(content)
	m.mainWindow.Resize(fyne.NewSize(400, 150))
	m.mainWindow.CenterOnScreen()
	m.mainWindow.SetFixedSize(true)
}

// SetCurrentStep은 현재 진행 단계를 업데이트합니다
func (m *Manager) SetCurrentStep(step int) {
	// step에 따른 적절한 메시지 설정
	messageKeys := map[int]string{
		0: "update.status.checking",
		1: "update.status.getting_info",
		2: "update.status.preparing",
		3: "update.status.downloading",
		4: "update.status.verifying",
		5: "update.status.installing",
		6: "update.status.finalizing",
		7: "update.status.completed",
	}
	if msgKey, ok := messageKeys[step]; ok {
		var progress float64
		if step == m.totalSteps-1 {
			progress = 1.0
		} else {
			progress = float64(step) / float64(m.totalSteps-1)
		}

		m.logger.Info("업데이트 진행률: %.2f%%, 단계: %d/%d", progress*100, step, m.totalSteps-1)
		m.statusCard.UpdateStatus(progress, m.i18n.GetMessage(msgKey))
	}
}

// UpdateDetail은 상세 메시지를 업데이트합니다
func (m *Manager) UpdateDetail(message string) {
	m.statusCard.UpdateStatus(-1, message)
}

// ShowError는 에러 메시지를 표시합니다
func (m *Manager) ShowError(err error) {
	m.statusCard.SetError(err.Error())
	if m.onRestore != nil {
		m.ShowRestoring()
		go m.onRestore()
	}
}

// ShowProgress는 다운로드 진행률을 표시합니다
func (m *Manager) ShowProgress(current, total int64) {
	m.statusCard.SetProgress(current, total)
}

// ShowRestoring은 복원 진행 중임을 표시합니다
func (m *Manager) ShowRestoring() {
	m.statusCard.SetRestoring(m.i18n.GetMessage("update.restore.in_progress"))
}

// ShowRestoreComplete는 복원 완료를 표시합니다
func (m *Manager) ShowRestoreComplete() {
	m.statusCard.SetRestoreComplete()
}

// SetRestoreHandler는 복구 핸들러를 설정합니다
func (m *Manager) SetRestoreHandler(handler func()) {
	m.onRestore = handler
}

// GetTotalSteps는 전체 단계 수를 반환합니다
func (m *Manager) GetTotalSteps() int {
	return m.totalSteps
}

// SetCompletionCallback은 완료 콜백을 설정합니다
func (m *Manager) SetCompletionCallback(callback func()) {
	m.completionCallback = callback
	if m.statusCard != nil {
		m.statusCard.SetCompletionCallback(func() {
			if m.completionCallback != nil && !m.animationComplete {
				m.animationComplete = true
				m.completionCallback()
			}
		})
	}
}

// Run은 UI를 실행합니다
func (m *Manager) Run() {
	m.mainWindow.ShowAndRun()
}

// Close는 UI를 종료합니다
func (m *Manager) Close() {
	m.mainWindow.Close()
}
