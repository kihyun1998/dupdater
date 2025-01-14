package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
	headerCard     *components.HeaderCard
	statusCard     *components.StatusCard
	stepIndicators []*components.StepIndicator
	contentBox     *fyne.Container

	// 복구 핸들러
	onRestore func()
}

// Config는 Manager 생성에 필요한 설정을 담는 구조체입니다
type Config struct {
	AppName     string
	FromVersion string
	ToVersion   string
	TotalSteps  int
	Logger      *logger.Logger
}

// New는 새로운 Manager 인스턴스를 생성합니다
func New(config Config) *Manager {
	manager := &Manager{
		app:    app.New(),
		logger: config.Logger,
	}

	// 메인 윈도우 생성
	manager.mainWindow = manager.app.NewWindow(
		fmt.Sprintf("%s Updater", config.AppName),
	)

	// UI 초기화
	manager.initializeUI(config)

	return manager
}

// initializeUI는 UI 컴포넌트들을 초기화하고 배치합니다
func (m *Manager) initializeUI(config Config) {
	// 헤더 카드 생성
	m.headerCard = components.NewHeaderCard(
		fmt.Sprintf("%s Update", config.AppName),
		config.FromVersion,
		config.ToVersion,
	)

	// 상태 카드 생성
	m.statusCard = components.NewStatusCard()

	// 단계 표시기 생성
	m.stepIndicators = make([]*components.StepIndicator, 3)
	m.stepIndicators[0] = components.NewStepIndicator("백업", "파일 백업 준비 중...")
	m.stepIndicators[1] = components.NewStepIndicator("다운로드", "업데이트 파일 다운로드 대기 중...")
	m.stepIndicators[2] = components.NewStepIndicator("설치", "설치 준비 중...")

	// 단계 표시기 컨테이너
	stepsContainer := container.NewVBox()
	for _, step := range m.stepIndicators {
		stepsContainer.Add(step)
	}

	// 전체 레이아웃 구성
	m.contentBox = container.NewVBox(
		m.headerCard,
		widget.NewSeparator(),
		stepsContainer, // 순서 변경 - 단계 표시를 먼저 보여줌
		widget.NewSeparator(),
		m.statusCard, // 상태 카드를 마지막에 배치
	)

	// 패딩 추가
	paddedContent := container.NewPadded(m.contentBox)

	// 메인 윈도우 설정
	m.mainWindow.SetContent(paddedContent)
	m.mainWindow.Resize(fyne.NewSize(500, 400))
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
	m.statusCard.UpdateStatus("업데이트 진행 중", message)
}

// ShowProgress는 다운로드 진행상황을 표시합니다
func (m *Manager) ShowProgress(current, total int64) {
	m.statusCard.SetProgress(current, total)
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

// triggerRestore는 복구 프로세스를 시작합니다
func (m *Manager) triggerRestore() {
	if m.onRestore != nil {
		go func() {
			m.UpdateDetail("파일 복구 중...")
			m.onRestore()
		}()
	}
}

// ShowRestoring은 복구 진행 중임을 표시합니다
func (m *Manager) ShowRestoring() {
	for _, step := range m.stepIndicators {
		step.UpdateStatus(components.StepPending)
	}
	m.statusCard.SetRestoring("파일을 원래 상태로 복원하고 있습니다...")
}

// ShowRestoreComplete는 복구 완료를 표시합니다
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
