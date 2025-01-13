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
	progressCard   *components.ProgressCard
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

	// 진행 카드 생성
	m.progressCard = components.NewProgressCard()

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
		m.progressCard,
		widget.NewSeparator(),
		stepsContainer,
	)

	// 패딩 추가
	paddedContent := container.NewPadded(m.contentBox)

	// 메인 윈도우 설정
	m.mainWindow.SetContent(paddedContent)
	m.mainWindow.Resize(fyne.NewSize(500, 400))
	m.mainWindow.CenterOnScreen()
	m.mainWindow.SetFixedSize(true)
}

// UI Manager 인터페이스 구현
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
	m.progressCard.UpdateMessage(message)
}

func (m *Manager) ShowProgress(current, total int64) {
	m.progressCard.UpdateProgress(current, total)
}

// ShowError는 에러 메시지를 표시합니다
func (m *Manager) ShowError(err error) {
	// 모든 단계를 실패 상태로 표시
	for _, step := range m.stepIndicators {
		step.UpdateStatus(components.StepFailed)
	}
	m.progressCard.SetError(err.Error())
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

// // ShowProgress는 진행률을 주기적으로 업데이트합니다 (다운로드 등에서 사용)
// func (m *Manager) ShowProgress(current, total int64) {
// 	progress := float64(current) / float64(total) * 100
// 	m.state.Progress = progress
// 	m.progressBar.SetValue(progress / 100)
// 	m.progressLabel.SetText(fmt.Sprintf("%.1f%%", progress))

// 	// 진행 상태 메시지 업데이트
// 	speed := float64(current) / (1024 * 1024) // MB 단위로 변환
// 	totalSize := float64(total) / (1024 * 1024)
// 	m.UpdateDetail(fmt.Sprintf("다운로드 중... %.1f MB / %.1f MB", speed, totalSize))
// }

// // ShowSuccess는 성공 메시지를 표시합니다
// func (m *Manager) ShowSuccess(message string) {
// 	successIcon := widget.NewIcon(theme.ConfirmIcon())
// 	successLabel := widget.NewLabelWithStyle(message, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

// 	content := container.NewVBox(
// 		container.NewHBox(layout.NewSpacer(), successIcon, successLabel, layout.NewSpacer()),
// 	)

// 	m.statusCard = widget.NewCard("", "", content)
// 	m.contentBox.Remove(m.errorCard)
// 	m.contentBox.Add(m.statusCard)
// 	m.mainWindow.Content().Refresh()

// 	// 3초 후 자동으로 창 닫기
// 	go func() {
// 		time.Sleep(3 * time.Second)
// 		m.mainWindow.Close()
// 	}()
// }
