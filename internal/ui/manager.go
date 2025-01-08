package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/kihyun1998/dupdater/internal/logger"
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
	state      State
	logger     *logger.Logger

	// UI 컴포넌트들
	progressBar   *widget.ProgressBar // 진행 막대
	stepLabel     *widget.Label       // 단계 표시 레이블
	progressLabel *widget.Label       // 진행률 표시 레이블
	detailLabel   *widget.Label       // 상세 메시지 레이블
	statusCard    *widget.Card        // 상태 표시 카드
	contentBox    *fyne.Container     // 메인 컨텐츠 컨테이너
	errorCard     *widget.Card        // 에러 표시 카드
	actionButtons *fyne.Container     // 작업 버튼 컨테이너

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
		app: app.New(),
		state: State{
			CurrentStep: -1,
			TotalSteps:  config.TotalSteps,
			Progress:    0,
			Detail:      "업데이트 준비 중...",
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

// initializeUI는 UI 컴포넌트들을 초기화하고 배치합니다
func (m *Manager) initializeUI() {
	// 헤더 섹션
	headerLabel := widget.NewLabelWithStyle("업데이트 관리자", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	header := container.NewHBox(
		layout.NewSpacer(),
		headerLabel,
		layout.NewSpacer(),
	)

	// 진행 상태 섹션
	m.progressBar = widget.NewProgressBar()
	m.progressLabel = widget.NewLabel("0%")
	m.stepLabel = widget.NewLabel(fmt.Sprintf("단계: 0/%d", m.state.TotalSteps))
	m.detailLabel = widget.NewLabel(m.state.Detail)
	m.detailLabel.Wrapping = fyne.TextWrapWord

	progressBox := container.NewVBox(
		container.NewHBox(
			m.progressBar,
			m.progressLabel,
		),
		m.stepLabel,
		widget.NewSeparator(),
		m.detailLabel,
	)

	// 상태 카드
	m.statusCard = widget.NewCard("", "", progressBox)

	// 작업 버튼 영역
	m.actionButtons = container.NewHBox(
		layout.NewSpacer(),
	)

	// 메인 컨텐츠 구성
	m.contentBox = container.NewVBox(
		header,
		widget.NewSeparator(),
		m.statusCard,
		m.actionButtons,
	)

	// 윈도우 설정
	m.mainWindow.SetContent(container.NewPadded(m.contentBox))
	m.mainWindow.Resize(fyne.NewSize(500, 300))
	m.mainWindow.SetFixedSize(true) // 윈도우 크기 고정
	m.mainWindow.CenterOnScreen()
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
	m.state.Progress = float64(step+1) / float64(m.state.TotalSteps) * 100

	// UI 업데이트
	m.progressBar.SetValue(m.state.Progress / 100)
	m.progressLabel.SetText(fmt.Sprintf("%.1f%%", m.state.Progress))
	m.stepLabel.SetText(fmt.Sprintf("단계: %d/%d", step+1, m.state.TotalSteps))
}

// UpdateDetail은 상세 메시지를 업데이트합니다
func (m *Manager) UpdateDetail(message string) {
	m.state.Detail = message
	m.detailLabel.SetText(message)
}

// ShowError는 에러 메시지를 표시합니다
func (m *Manager) ShowError(err error) {
	m.logger.Error("오류 발생: %v", err)

	// 에러 메시지 생성
	errorText := widget.NewLabel(fmt.Sprintf("오류: %v", err))
	errorText.Wrapping = fyne.TextWrapBreak
	errorText.TextStyle = fyne.TextStyle{Bold: true}
	errorText.Wrapping = fyne.TextWrapBreak // 텍스트 자동 줄바꿈

	// 스크롤 가능한 컨테이너에 에러 메시지 배치
	errorScroll := container.NewScroll(errorText)
	errorScroll.SetMinSize(fyne.NewSize(460, 100)) // 스크롤 영역 크기 고정

	// 버튼 생성
	copyButton := widget.NewButtonWithIcon("복사", theme.ContentCopyIcon(), func() {
		m.mainWindow.Clipboard().SetContent(err.Error())
		dialog.ShowInformation("알림", "오류 메시지가 클립보드에 복사되었습니다.", m.mainWindow)
	})

	restoreButton := widget.NewButtonWithIcon("복구", theme.ViewRefreshIcon(), func() {
		dialog.ShowConfirm("복구 확인", "파일을 복구하시겠습니까?", func(restore bool) {
			if restore {
				m.triggerRestore()
			}
		}, m.mainWindow)
	})

	// 버튼 컨테이너
	buttons := container.NewHBox(
		layout.NewSpacer(),
		copyButton,
		restoreButton,
	)

	// 에러 카드 생성
	m.errorCard = widget.NewCard(
		"업데이트 오류",
		"",
		container.NewVBox(
			errorScroll, // 스크롤 가능한 에러 메시지
			widget.NewSeparator(),
			buttons,
		),
	)

	// 카드 크기 고정
	m.errorCard.Resize(fyne.NewSize(480, 250))

	// 컨텐츠 업데이트
	m.contentBox.Remove(m.statusCard)
	m.contentBox.Add(m.errorCard)
	m.mainWindow.Content().Refresh()
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

// ShowProgress는 진행률을 주기적으로 업데이트합니다 (다운로드 등에서 사용)
func (m *Manager) ShowProgress(current, total int64) {
	progress := float64(current) / float64(total) * 100
	m.state.Progress = progress
	m.progressBar.SetValue(progress / 100)
	m.progressLabel.SetText(fmt.Sprintf("%.1f%%", progress))

	// 진행 상태 메시지 업데이트
	speed := float64(current) / (1024 * 1024) // MB 단위로 변환
	totalSize := float64(total) / (1024 * 1024)
	m.UpdateDetail(fmt.Sprintf("다운로드 중... %.1f MB / %.1f MB", speed, totalSize))
}

// ShowSuccess는 성공 메시지를 표시합니다
func (m *Manager) ShowSuccess(message string) {
	successIcon := widget.NewIcon(theme.ConfirmIcon())
	successLabel := widget.NewLabelWithStyle(message, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	content := container.NewVBox(
		container.NewHBox(layout.NewSpacer(), successIcon, successLabel, layout.NewSpacer()),
	)

	m.statusCard = widget.NewCard("", "", content)
	m.contentBox.Remove(m.errorCard)
	m.contentBox.Add(m.statusCard)
	m.mainWindow.Content().Refresh()

	// 3초 후 자동으로 창 닫기
	go func() {
		time.Sleep(3 * time.Second)
		m.mainWindow.Close()
	}()
}
