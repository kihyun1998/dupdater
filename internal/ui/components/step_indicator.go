package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kihyun1998/dupdater/internal/ui/theme"
)

// StepStatus는 단계의 상태를 정의합니다
type StepStatus int

const (
	StepPending StepStatus = iota
	StepInProgress
	StepCompleted
	StepFailed
)

// StepIndicator는 업데이트 과정의 각 단계를 표시하는 컴포넌트입니다
type StepIndicator struct {
	widget.BaseWidget
	title  string     // 단계 제목
	detail string     // 상세 설명
	status StepStatus // 현재 상태

	container  *fyne.Container
	icon       *canvas.Image // 상태 아이콘
	titleText  *canvas.Text  // 제목 텍스트
	detailText *canvas.Text  // 상세 설명 텍스트
}

// NewStepIndicator는 새로운 StepIndicator를 생성합니다
func NewStepIndicator(title, detail string) *StepIndicator {
	step := &StepIndicator{
		title:  title,
		detail: detail,
		status: StepPending,
	}
	step.ExtendBaseWidget(step)
	step.setupUI()
	return step
}

// CreateRenderer는 Fyne 위젯 인터페이스를 구현합니다
func (s *StepIndicator) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(s.container)
}

// setupUI는 UI 컴포넌트를 초기화합니다
func (s *StepIndicator) setupUI() {
	// 아이콘 설정 (기본값은 대기 상태)
	s.icon = canvas.NewImageFromResource(resourceProgressIconSvg)
	s.icon.Resize(fyne.NewSize(20, 20))
	s.icon.FillMode = canvas.ImageFillOriginal

	// 제목 텍스트
	s.titleText = canvas.NewText(s.title, theme.TextColor)
	s.titleText.TextSize = theme.FontSizeMedium

	// 상세 설명 텍스트
	s.detailText = canvas.NewText(s.detail, theme.SubTextColor)
	s.detailText.TextSize = theme.FontSizeSmall

	// 텍스트 컨테이너
	textContainer := container.NewVBox(
		s.titleText,
		s.detailText,
	)

	// 메인 컨테이너 (아이콘 + 텍스트)
	s.container = container.NewHBox(
		container.NewPadded(s.icon),
		textContainer,
	)
}

// UpdateStatus는 단계의 상태를 업데이트합니다
func (s *StepIndicator) UpdateStatus(status StepStatus) {
	s.status = status

	// 상태에 따른 아이콘과 색상 업데이트
	var iconResource fyne.Resource
	var titleColor color.Color

	switch s.status {
	case StepCompleted:
		iconResource = resourceCompletedIconSvg
		titleColor = theme.SuccessColor
	case StepInProgress:
		iconResource = resourceProgressIconSvg
		// titleColor = theme.InfoColor
	case StepFailed:
		iconResource = resourceFailedIconSvg
		titleColor = theme.ErrorColor
	default:
		iconResource = resourceProgressIconSvg
		titleColor = theme.SubTextColor
	}

	s.icon.Resource = iconResource
	s.titleText.Color = titleColor
	s.icon.Refresh()
	s.titleText.Refresh()
	s.container.Refresh()
}

// UpdateDetail은 상세 설명을 업데이트합니다
func (s *StepIndicator) UpdateDetail(detail string) {
	s.detail = detail
	s.detailText.Text = detail
	s.detailText.Refresh()
}
