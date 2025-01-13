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
	title      string
	detail     string
	status     StepStatus
	container  *fyne.Container
	icon       *canvas.Image
	titleText  *canvas.Text
	detailText *canvas.Text
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
	// 아이콘 설정
	s.icon = canvas.NewImageFromResource(resourcePendingIconSvg)
	s.icon.Resize(fyne.NewSize(24, 24))
	s.icon.FillMode = canvas.ImageFillOriginal

	// 제목 텍스트
	s.titleText = canvas.NewText(s.title, theme.TextColor)
	s.titleText.TextStyle = fyne.TextStyle{Bold: true}
	s.titleText.TextSize = 14

	// 상세 설명 텍스트
	s.detailText = canvas.NewText(s.detail, theme.SubTextColor)
	s.detailText.TextSize = 12

	// 텍스트 컨테이너
	textContainer := container.NewVBox(
		s.titleText,
		s.detailText,
	)

	// 메인 컨테이너
	s.container = container.NewHBox(
		s.icon,
		container.NewPadded(textContainer),
	)
}

// UpdateStatus는 단계의 상태를 업데이트합니다
func (s *StepIndicator) UpdateStatus(status StepStatus) {
	s.status = status
	s.updateUI()
}

// UpdateDetail은 상세 설명을 업데이트합니다
func (s *StepIndicator) UpdateDetail(detail string) {
	s.detail = detail
	if s.detailText != nil {
		s.detailText.Text = detail
		s.detailText.Refresh()
	}
}

// updateUI는 상태에 따라 UI를 업데이트합니다
func (s *StepIndicator) updateUI() {
	var iconResource fyne.Resource
	var textColor color.Color

	switch s.status {
	case StepPending:
		iconResource = resourcePendingIconSvg
		textColor = theme.SubTextColor
	case StepInProgress:
		iconResource = resourceProgressIconSvg
		textColor = theme.InfoColor
	case StepCompleted:
		iconResource = resourceCompletedIconSvg
		textColor = theme.SuccessColor
	case StepFailed:
		iconResource = resourceFailedIconSvg
		textColor = theme.ErrorColor
	}

	s.icon.Resource = iconResource
	s.titleText.Color = textColor

	s.icon.Refresh()
	s.titleText.Refresh()
	s.container.Refresh()
}
