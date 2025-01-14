package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/ui/components"
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
	logger     *logger.Logger

	// UI 컴포넌트들
	statusCard     *components.StatusCard      // 상태 표시 컴포넌트
	stepIndicators []*components.StepIndicator // 단계 표시 컴포넌트들
	contentBox     *fyne.Container

	// 복구 핸들러
	onRestore func()
}

// Config는 Manager 생성에 필요한 설정을 담는 구조체입니다
type Config struct {
	AppName    string
	TotalSteps int
	Logger     *logger.Logger
}

// New는 새로운 Manager 인스턴스를 생성합니다
func New(config Config) *Manager {
	manager := &Manager{
		app:    app.New(),
		logger: config.Logger,
	}

	// 메인 윈도우 생성
	manager.mainWindow = manager.app.NewWindow(fmt.Sprintf("%s Updater", config.AppName))

	// UI 초기화
	manager.initializeUI(config)

	return manager
}

// initializeUI는 UI 컴포넌트들을 초기화하고 배치합니다
func (m *Manager) initializeUI(config Config) {
	// 상태 카드 생성
	m.statusCard = components.NewStatusCard()

	// 단계 표시기들 생성
	m.stepIndicators = make([]*components.StepIndicator, 3)
	m.stepIndicators[0] = components.NewStepIndicator("Backup completed", "Files backed up successfully")
	m.stepIndicators[1] = components.NewStepIndicator("Download completed", "Update package verified")
	m.stepIndicators[2] = components.NewStepIndicator("Installing update", "This may take a few minutes")

	// 단계 표시기 컨테이너
	stepsContainer := container.NewVBox()
	for _, step := range m.stepIndicators {
		stepsContainer.Add(step)
	}

	// 전체 레이아웃 구성
	m.contentBox = container.NewVBox(
		m.statusCard,   // 상태 카드를 최상단에 배치
		stepsContainer, // 단계 표시기들을 그 아래에 배치
	)

	// 메인 윈도우 설정
	m.mainWindow.SetContent(container.NewPadded(m.contentBox))
	m.mainWindow.Resize(fyne.NewSize(400, 500))
	m.mainWindow.CenterOnScreen()
	m.mainWindow.SetFixedSize(true)
}

// SetCurrentStep은 현재 진행 단계를 업데이트합니다
func (m *Manager) SetCurrentStep(step int) {
	if step >= 0 && step < len(m.stepIndicators) {
		// 이전 단계들을 완료 상태로 설정
		for i := 0; i < step; i++ {
			m.stepIndicators[i].UpdateStatus(components.StepCompleted)
		}
		// 현재 단계를 진행 중 상태로 설정
		m.stepIndicators[step].UpdateStatus(components.StepInProgress)
		// 다음 단계들을 대기 상태로 설정
		for i := step + 1; i < len(m.stepIndicators); i++ {
			m.stepIndicators[i].UpdateStatus(components.StepPending)
		}
	}
}

// UpdateDetail은 상세 메시지를 업데이트합니다
func (m *Manager) UpdateDetail(message string) {
	m.statusCard.UpdateStatus("System Update", message)
}

// ShowError는 에러 메시지를 표시합니다
func (m *Manager) ShowError(err error) {
	// 모든 단계를 실패 상태로 표시
	for _, step := range m.stepIndicators {
		step.UpdateStatus(components.StepFailed)
	}
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
	m.statusCard.SetRestoring("Restoring files...")
}

// ShowRestoreComplete는 복원 완료를 표시합니다
func (m *Manager) ShowRestoreComplete() {
	m.statusCard.SetRestoreComplete()
}

// SetRestoreHandler는 복구 핸들러를 설정합니다
func (m *Manager) SetRestoreHandler(handler func()) {
	m.onRestore = handler
}

// Run은 UI를 실행합니다
func (m *Manager) Run() {
	m.mainWindow.ShowAndRun()
}

// Close는 UI를 종료합니다
func (m *Manager) Close() {
	m.mainWindow.Close()
}
