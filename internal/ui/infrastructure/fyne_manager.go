// internal/ui/infrastructure/fyne_manager.go

package infrastructure

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/kihyun1998/dupdater/internal/i18n"
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/ui/components"
	"github.com/kihyun1998/dupdater/internal/ui/domain/entity"
	"github.com/kihyun1998/dupdater/internal/ui/domain/repository"
	"github.com/kihyun1998/dupdater/pkg/utils/logo"
	"github.com/kihyun1998/dupdater/pkg/utils/theme"
)

// FyneManager는 Fyne 기반 UI 구현체입니다
type FyneManager struct {
	app        fyne.App
	mainWindow fyne.Window
	statusCard *components.StatusCard
	state      *entity.UIState
	logger     logger.Logger
	i18n       i18n.Manager
	theme      theme.ThemeVariant
	mu         sync.RWMutex
}

// NewFyneManager는 새로운 FyneManager를 생성합니다
func NewFyneManager(config *repository.Config) *FyneManager {
	// Fyne 앱 생성
	fyneApp := app.New()
	customTheme := theme.NewCustomTheme(config.Theme)
	fyneApp.Settings().SetTheme(customTheme)

	// 상태 초기화
	state := entity.NewUIState(
		config.TotalSteps,
		config.FromVersion,
		config.ToVersion,
	)

	manager := &FyneManager{
		app:    fyneApp,
		state:  state,
		logger: config.Logger,
		i18n:   config.I18n,
		theme:  config.Theme,
	}

	// 메인 윈도우 초기화
	manager.initWindow()
	manager.initUI()

	return manager
}

// initWindow은 메인 윈도우를 초기화합니다
func (m *FyneManager) initWindow() {
	m.mainWindow = m.app.NewWindow(m.i18n.GetMessage("update.title"))

	m.mainWindow.Resize(fyne.NewSize(400, 150))
	m.mainWindow.CenterOnScreen()
	m.mainWindow.SetIcon(logo.ACRALogo)
	m.mainWindow.SetFixedSize(false)

}

// initUI는 UI 컴포넌트를 초기화합니다
func (m *FyneManager) initUI() {
	m.statusCard = components.NewStatusCard(
		m.state.FromVersion,
		m.state.ToVersion,
		m.theme,
		m.i18n,
	)

	content := container.NewPadded(m.statusCard)
	m.mainWindow.SetContent(content)
}

// GetState는 현재 UI 상태를 반환합니다
func (m *FyneManager) GetState() *entity.UIState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}

// UpdateState는 UI 상태를 업데이트합니다
func (m *FyneManager) UpdateState(state *entity.UIState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state = state
}

// SetCurrentStep은 현재 진행 단계를 설정합니다
func (m *FyneManager) SetCurrentStep(step int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	progress := float64(step) / float64(m.state.TotalSteps-1)
	m.statusCard.UpdateStatus(progress, m.state.DetailMessage)
}

// UpdateDetail은 상세 메시지를 업데이트합니다
func (m *FyneManager) UpdateDetail(message string) {
	m.statusCard.UpdateStatus(-1, message)
}

// ShowError는 에러 상태를 표시합니다
func (m *FyneManager) ShowError(err error) {
	m.statusCard.SetError(err.Error())
}

// Run은 UI를 실행합니다
func (m *FyneManager) Run() {
	m.mainWindow.ShowAndRun()
}

// Close는 UI를 종료합니다
func (m *FyneManager) Close() {
	m.mainWindow.Close()
}

// SetRestoreHandler는 복원 핸들러를 설정합니다
func (m *FyneManager) SetRestoreHandler(handler func()) {
	m.statusCard.SetRestoreHandler(handler)
}

// ShowRestoring은 복원 중 상태를 표시합니다
func (m *FyneManager) ShowRestoring() {
	m.statusCard.SetRestoring(m.i18n.GetMessage("update.restore.in_progress"))
}

// ShowRestoreComplete는 복원 완료 상태를 표시합니다
func (m *FyneManager) ShowRestoreComplete() {
	m.statusCard.SetRestoreComplete()
}

// GetTotalSteps는 전체 단계 수를 반환합니다
func (m *FyneManager) GetTotalSteps() int {
	return m.state.TotalSteps
}

// SetCompletionCallback은 완료 콜백을 설정합니다
func (m *FyneManager) SetCompletionCallback(callback func()) {
	m.statusCard.SetCompletionCallback(callback)
}
