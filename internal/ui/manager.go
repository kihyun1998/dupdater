package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/ui/components"
	"github.com/kihyun1998/dupdater/internal/ui/theme"
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
	app          fyne.App
	mainWindow   fyne.Window
	logger       *logger.Logger
	totalSteps   int
	currentStep  int
	currentTheme theme.Theme
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
	Logger      *logger.Logger
	Theme       theme.Theme
}

// New는 새로운 Manager 인스턴스를 생성합니다
func New(config Config) *Manager {
	manager := &Manager{
		app:               app.New(),
		logger:            config.Logger,
		totalSteps:        config.TotalSteps,
		currentStep:       0,
		animationComplete: false,
		currentTheme:      config.Theme,
	}

	manager.mainWindow = manager.app.NewWindow(fmt.Sprintf("%s Updater", config.AppName))
	manager.initializeUI(config)

	return manager
}

// initializeUI는 UI 컴포넌트들을 초기화하고 배치합니다
func (m *Manager) initializeUI(config Config) {
	// StatusCard 생성 시 테마 전달
	m.statusCard = components.NewStatusCard(config.FromVersion, config.ToVersion, m.currentTheme)

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
	content.Resize(fyne.NewSize(400, 200))

	// 테마에 따른 배경색 설정
	m.mainWindow.SetContent(content)
	m.mainWindow.Resize(fyne.NewSize(400, 200))
	m.mainWindow.CenterOnScreen()
	m.mainWindow.SetFixedSize(true)
}

// SetCurrentStep은 현재 진행 단계를 업데이트합니다
func (m *Manager) SetCurrentStep(step int) {
	// step에 따른 적절한 메시지 설정
	messages := map[int]string{
		0: "앱 상태를 확인하고 있습니다...",
		1: "업데이트 정보를 확인하고 있습니다...",
		2: "업데이트 파일을 준비하고 있습니다...",
		3: "업데이트 파일을 다운로드하고 있습니다...",
		4: "업데이트 파일을 검증하고 있습니다...",
		5: "업데이트를 설치하고 있습니다...",
		6: "설치를 확인하고 있습니다...",
		7: "업데이트가 완료되었습니다.",
	}
	if msg, ok := messages[step]; ok {
		var progress float64
		if step == m.totalSteps-1 {
			// 마지막 단계에서는 100%로 설정
			progress = 1.0
		} else {
			// 그 외의 경우 진행률 계산
			progress = float64(step) / float64(m.totalSteps-1)
		}

		m.logger.Info("업데이트 진행률: %.2f%%, 단계: %d/%d", progress*100, step, m.totalSteps-1)
		m.statusCard.UpdateStatus(progress, msg)
	}
}

// GetTotalSteps은 전체 단계 수를 반환합니다
func (m *Manager) GetTotalSteps() int {
	return m.totalSteps
}

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

// UpdateDetail은 상세 메시지를 업데이트합니다
func (m *Manager) UpdateDetail(message string) {
	m.statusCard.UpdateStatus(-1, message) // -1은 진행률 변경 없음을 의미
}

// ShowError는 에러 메시지를 표시합니다
func (m *Manager) ShowError(err error) {
	m.statusCard.SetError(err.Error())

	// 복구 가능한 경우 복구 UI 표시
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
	m.statusCard.SetRestoring("이전 버전으로 복원중...")
}

// ShowRestoreComplete는 복원 완료를 표시합니다
func (m *Manager) ShowRestoreComplete() {
	m.statusCard.SetRestoreComplete()
}

// SetRestoreHandler는 복구 핸들러를 설정합니다
func (m *Manager) SetRestoreHandler(handler func()) {
	m.onRestore = handler
}

// 상태 카드의 애니메이션 완료 콜백을 처리하는 메서드 추가
func (m *Manager) onAnimationComplete() {
	if !m.animationComplete && m.currentStep == m.totalSteps-1 {
		m.animationComplete = true
		if m.completionCallback != nil {
			m.completionCallback()
		}
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
