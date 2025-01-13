package components

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kihyun1998/dupdater/internal/ui/theme"
)

// ProgressCard는 진행 상황을 표시하는 컴포넌트입니다
type ProgressCard struct {
	widget.BaseWidget
	container    *fyne.Container
	progressBar  *widget.ProgressBar
	progressText *canvas.Text
	messageText  *canvas.Text
	currentValue float64
	message      string
}

// NewProgressCard는 새로운 ProgressCard를 생성합니다
func NewProgressCard() *ProgressCard {
	card := &ProgressCard{
		currentValue: 0,
		message:      "준비 중...",
	}
	card.ExtendBaseWidget(card)
	card.setupUI()
	return card
}

// CreateRenderer는 Fyne 위젯 인터페이스를 구현합니다
func (p *ProgressCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(p.container)
}

// setupUI는 UI 컴포넌트를 초기화합니다
func (p *ProgressCard) setupUI() {
	// 진행률 텍스트 (우측 상단)
	p.progressText = canvas.NewText("0%", theme.TextColor)
	p.progressText.TextSize = theme.FontSizeMedium
	p.progressText.Alignment = fyne.TextAlignTrailing

	// 진행 바
	p.progressBar = widget.NewProgressBar()
	p.progressBar.Min = 0
	p.progressBar.Max = 100
	p.progressBar.Value = p.currentValue

	// 메시지 텍스트
	p.messageText = canvas.NewText(p.message, theme.SubTextColor)
	p.messageText.TextSize = theme.FontSizeSmall

	// 상단 컨테이너 (진행률 텍스트)
	topContainer := container.NewHBox(
		layout.NewSpacer(),
		p.progressText,
	)

	// 메인 컨테이너
	p.container = container.NewVBox(
		topContainer,
		p.progressBar,
		p.messageText,
	)

	// 패딩 추가
	p.container = container.NewPadded(p.container)
}

// UpdateProgress는 진행률을 업데이트합니다
func (p *ProgressCard) UpdateProgress(current, total int64) {
	if total > 0 {
		percentage := float64(current) / float64(total) * 100
		p.currentValue = percentage
		p.progressBar.Value = percentage
		p.progressText.Text = fmt.Sprintf("%.1f%%", percentage)

		// 진행률에 따른 메시지 업데이트
		p.messageText.Text = fmt.Sprintf("다운로드 중... %.1f MB / %.1f MB",
			float64(current)/(1024*1024),
			float64(total)/(1024*1024))

		// UI 새로고침
		p.progressText.Refresh()
		p.progressBar.Refresh()
		p.messageText.Refresh()
	}
}

// UpdateMessage는 메시지를 업데이트합니다
func (p *ProgressCard) UpdateMessage(message string) {
	p.message = message
	p.messageText.Text = message
	p.messageText.Refresh()
}

// MinSize는 최소 크기를 반환합니다
func (p *ProgressCard) MinSize() fyne.Size {
	return fyne.NewSize(400, 100)
}

// SetComplete는 진행이 완료되었음을 표시합니다
func (p *ProgressCard) SetComplete() {
	p.currentValue = 100
	p.progressBar.Value = 100
	p.progressText.Text = "100%"
	p.messageText.Text = "완료됨"
	p.messageText.Color = theme.SuccessColor

	// UI 새로고침
	p.progressText.Refresh()
	p.progressBar.Refresh()
	p.messageText.Refresh()
}

// SetError는 오류 상태를 표시합니다
func (p *ProgressCard) SetError(errorMsg string) {
	p.messageText.Text = fmt.Sprintf("오류: %s", errorMsg)
	p.messageText.Color = theme.ErrorColor
	p.messageText.Refresh()
}
