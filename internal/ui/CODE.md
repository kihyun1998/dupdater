# ui
## Project Structure

```
ui/
├── components/
    └── status_card.go
└── manager.go
```

## components/status_card.go
```go
package components

import (
	"fmt"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kihyun1998/dupdater/internal/i18n"
	"github.com/kihyun1998/dupdater/pkg/utils/theme"
)

// StatusCard는 현재 상태를 표시하는 컴포넌트입니다
type StatusCard struct {
	widget.BaseWidget
	container   *fyne.Container
	progressBar *widget.ProgressBar

	titleText    *canvas.Text
	versionText  *canvas.Text
	subtitleText *canvas.Text

	currentTheme theme.ThemeVariant
	i18n         i18n.Manager

	targetProgress  float64
	currentProgress float64
	animating       bool
	ticker          *time.Ticker
	onComplete      func()
}

// NewStatusCard는 새로운 StatusCard를 생성합니다
func NewStatusCard(fromVersion, toVersion string, themeVariant theme.ThemeVariant, i18n i18n.Manager) *StatusCard {
	card := &StatusCard{
		targetProgress:  0,
		currentProgress: 0,
		animating:       false,
		currentTheme:    themeVariant,
		i18n:            i18n,
	}
	card.ExtendBaseWidget(card)
	card.setupUI(fromVersion, toVersion)
	return card
}

// CreateRenderer는 Fyne 위젯 인터페이스를 구현합니다
func (s *StatusCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(s.container)
}

// setupUI는 UI 컴포넌트를 초기화합니다
func (s *StatusCard) setupUI(fromVersion, toVersion string) {
	// 타이틀 텍스트
	s.titleText = canvas.NewText(
		s.i18n.GetMessage("update.header.new_update"),
		s.currentTheme.TextColor(),
	)
	s.titleText.TextSize = s.currentTheme.FontSizeLarge()
	s.titleText.TextStyle = fyne.TextStyle{Bold: true}

	// 버전 정보
	versionFormat := fmt.Sprintf("%s → %s", fromVersion, toVersion)
	s.versionText = canvas.NewText(versionFormat, s.currentTheme.SubTextColor())
	s.versionText.TextSize = s.currentTheme.FontSizeMedium()

	// 진행 바
	s.progressBar = widget.NewProgressBar()
	s.progressBar.Resize(fyne.NewSize(350, 20))

	// 상태 메시지
	s.subtitleText = canvas.NewText(
		s.i18n.GetMessage("update.info.restart"),
		s.currentTheme.SubTextColor(),
	)
	s.subtitleText.TextSize = s.currentTheme.FontSizeSmall()

	// 레이아웃 구성
	s.container = container.NewVBox(
		container.NewVBox(
			s.titleText,
			s.versionText,
		),
		container.NewPadded(s.progressBar),
		container.NewPadded(s.subtitleText),
	)
}

// SetError는 에러 상태를 표시합니다
func (s *StatusCard) SetError(errMsg string) {
	s.titleText.Text = s.i18n.GetMessage("update.error.generic")
	s.titleText.Color = s.currentTheme.ErrorColor()
	s.subtitleText.Text = errMsg
	s.subtitleText.Color = s.currentTheme.ErrorColor()
	s.progressBar.Hide()
	s.Refresh()
}

// SetRestoring는 복원 진행 상태를 표시합니다
func (s *StatusCard) SetRestoring(msg string) {
	s.titleText.Text = s.i18n.GetMessage("update.restore.in_progress")
	s.titleText.Color = s.currentTheme.WarningColor()
	s.subtitleText.Text = msg
	s.subtitleText.Color = s.currentTheme.WarningColor()
	s.progressBar.Show()
	s.Refresh()
}

// SetRestoreComplete는 복원 완료 상태를 표시합니다
func (s *StatusCard) SetRestoreComplete() {
	s.titleText.Text = s.i18n.GetMessage("update.restore.completed")
	s.titleText.Color = s.currentTheme.SuccessColor()
	s.subtitleText.Text = s.i18n.GetMessage("update.info.restart")
	s.subtitleText.Color = s.currentTheme.SuccessColor()
	s.progressBar.Hide()
	s.Refresh()
}

// UpdateStatus는 상태와 진행률을 업데이트합니다
func (s *StatusCard) UpdateStatus(progress float64, message string) {
	s.titleText.Color = s.currentTheme.TextColor()
	s.subtitleText.Text = message
	s.subtitleText.Color = s.currentTheme.SubTextColor()

	if progress >= 0 {
		s.targetProgress = progress
		s.animateProgress()
	}

	s.progressBar.Show()
	s.Refresh()
}

// SetProgress는 다운로드 진행률을 업데이트합니다
func (s *StatusCard) SetProgress(current, total int64) {
	progress := float64(current) / float64(total)
	s.progressBar.SetValue(progress)
	s.Refresh()
}

// SetCompletionCallback은 진행률 100% 도달 시 실행될 콜백을 설정합니다
func (s *StatusCard) SetCompletionCallback(callback func()) {
	s.onComplete = callback
}

// Refresh는 위젯을 새로고침합니다
func (s *StatusCard) Refresh() {
	s.container.Refresh()
}

// 부드러운 진행률 업데이트를 위한 메서드
func (s *StatusCard) animateProgress() {
	if s.ticker != nil {
		s.ticker.Stop()
	}

	s.animating = true
	s.ticker = time.NewTicker(16 * time.Millisecond) // 약 60FPS

	go func() {
		defer s.ticker.Stop()

		for range s.ticker.C {
			if !s.animating {
				return
			}

			diff := s.targetProgress - s.currentProgress
			step := diff * 0.1

			if math.Abs(diff) < 0.001 {
				s.currentProgress = s.targetProgress
				s.progressBar.SetValue(s.currentProgress)
				s.progressBar.Refresh()
				s.animating = false

				// 애니메이션이 100%에 도달했을 때 콜백 실행
				if s.currentProgress >= 0.999 {
					if s.onComplete != nil {
						s.onComplete()
					}
				}
				return
			}

			s.currentProgress += step
			s.progressBar.SetValue(s.currentProgress)
			s.progressBar.Refresh()
		}
	}()
}

```
## manager.go
```go
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

```
