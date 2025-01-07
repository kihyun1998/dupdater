package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// State는 UI의 현재 상태를 나타내는 구조체입니다
type State struct {
	CurrentStep int    // 현재 진행 단계
	TotalSteps  int    // 전체 단계 수
	Detail      string // 상세 메시지
}

// Manager는 UI를 관리하는 구조체입니다
type Manager struct {
	app        fyne.App
	mainWindow fyne.Window
	state      State
	logger     Logger

	// UI 컴포넌트들
	spinnerIcon *canvas.Image
	stepLabel   *widget.Label
	detailLabel *widget.Label

	// 복구 핸들러
	onRestore func()
}

// Config는 Manager 생성에 필요한 설정을 담는 구조체입니다
type Config struct {
	AppName    string // 어플리케이션 이름
	TotalSteps int    // 전체 단계 수
	Logger     Logger // 로거
}

// New는 새로운 Manager 인스턴스를 생성합니다
func New(config Config) *Manager {
	manager := &Manager{
		app: app.New(),
		state: State{
			CurrentStep: -1,
			TotalSteps:  config.TotalSteps,
			Detail:      "",
		},
		logger: config.Logger,
	}

	// 메인 윈도우 생성
	manager.mainWindow = manager.app.NewWindow(
		fmt.Sprintf("%s Updater", config.AppName),
	)

	// UI 초기화
	manager.initializeUI()

	return manager
}

// Run은 UI를 실행합니다
func (m *Manager) Run() {
	m.mainWindow.ShowAndRun()
}

// Close는 UI를 종료합니다
func (m *Manager) Close() {
	m.mainWindow.Close()
}

// SetCurrentStep은 현재 진행 단계를 업데이트합니다
func (m *Manager) SetCurrentStep(step int) {
	if step < 0 || step >= m.state.TotalSteps {
		m.logger.Error("잘못된 단계 번호: %d", step)
		return
	}

	m.state.CurrentStep = step
	m.updateStepLabel()
}

// UpdateDetail은 상세 메시지를 업데이트합니다
func (m *Manager) UpdateDetail(message string) {
	m.state.Detail = message
	m.updateDetailLabel()
}

// ShowError는 에러 메시지를 표시합니다
func (m *Manager) ShowError(err error) {
	m.logger.Error("UI 에러 발생: %v", err)

	// 에러 표시용 컨테이너 생성
	errorText := widget.NewLabel(fmt.Sprintf("에러: %v", err))

	// 복구 버튼 생성
	restoreButton := widget.NewButton("복구", func() {
		m.logger.Info("복구 버튼 클릭됨")
		m.triggerRestore()
	})

	// 컨테이너에 위젯 추가
	content := container.NewVBox(
		errorText,
		restoreButton,
	)

	// 메인 윈도우 컨텐츠 교체
	m.mainWindow.SetContent(content)
}

// 내부 메서드들

func (m *Manager) initializeUI() {
	// 아이콘 대신 프로그레스 바 사용
	progress := widget.NewProgressBar()
	progress.Resize(fyne.NewSize(64, 64))

	// 레이블 초기화
	m.stepLabel = widget.NewLabel("")
	m.stepLabel.TextStyle = fyne.TextStyle{Bold: true}

	m.detailLabel = widget.NewLabel("")
	m.detailLabel.TextStyle = fyne.TextStyle{Italic: true}
	m.detailLabel.Alignment = fyne.TextAlignCenter

	// 레이아웃 구성
	content := container.NewVBox(
		progress,
		m.stepLabel,
		m.detailLabel,
	)

	// 여백 설정
	paddedContent := container.NewPadded(content)

	// 윈도우에 컨텐츠 설정
	m.mainWindow.SetContent(paddedContent)
	m.mainWindow.Resize(fyne.NewSize(400, 200))
	m.mainWindow.CenterOnScreen()

	// 초기 상태 업데이트
	m.updateStepLabel()
	m.updateDetailLabel()
}

func (m *Manager) updateStepLabel() {
	if m.state.CurrentStep < 0 {
		m.stepLabel.SetText("준비 중...")
		return
	}

	m.stepLabel.SetText(fmt.Sprintf(
		"진행 단계: %d / %d",
		m.state.CurrentStep+1,
		m.state.TotalSteps,
	))
}

func (m *Manager) updateDetailLabel() {
	m.detailLabel.SetText(m.state.Detail)
}

func (m *Manager) triggerRestore() {
	if m.onRestore != nil {
		go m.onRestore()
	}
}

// SetRestoreHandler는 복구 핸들러를 설정합니다
func (m *Manager) SetRestoreHandler(handler func()) {
	m.onRestore = handler
}

// Logger는 로깅을 위한 인터페이스입니다
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}
