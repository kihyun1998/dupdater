package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kihyun1998/dupdater/internal/ui/theme"
)

// StatusCard는 현재 상태를 표시하는 컴포넌트입니다
type StatusCard struct {
	widget.BaseWidget
	container    *fyne.Container
	titleText    *canvas.Text        // "System Update"
	subtitleText *canvas.Text        // 상태 메시지 (예: "Please keep the app open")
	progressBar  *widget.ProgressBar // 진행 상태 표시용 (복원 시에만 표시)
}

// NewStatusCard는 새로운 StatusCard를 생성합니다
func NewStatusCard() *StatusCard {
	card := &StatusCard{}
	card.ExtendBaseWidget(card)
	card.setupUI()
	return card
}

// CreateRenderer는 Fyne 위젯 인터페이스를 구현합니다
func (s *StatusCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(s.container)
}

// setupUI는 UI 컴포넌트를 초기화합니다
func (s *StatusCard) setupUI() {
	// 메인 타이틀
	s.titleText = canvas.NewText("System Update", theme.TextColor)
	s.titleText.TextSize = theme.FontSizeLarge

	// 서브타이틀 (상태 메시지)
	s.subtitleText = canvas.NewText("", theme.SubTextColor)
	s.subtitleText.TextSize = theme.FontSizeSmall

	// 진행 바 (기본적으로 숨김)
	s.progressBar = widget.NewProgressBar()
	s.progressBar.Hide()

	// 컨테이너 구성
	s.container = container.NewVBox(
		s.titleText,
		s.subtitleText,
		s.progressBar,
	)
}

// UpdateStatus는 상태 메시지를 업데이트합니다
func (s *StatusCard) UpdateStatus(status, detail string) {
	s.titleText.Text = "System Update"
	s.subtitleText.Text = detail
	s.subtitleText.Color = theme.SubTextColor
	s.subtitleText.Refresh()
	s.progressBar.Hide()
	s.container.Refresh()
}

// SetError는 에러 상태를 표시합니다
func (s *StatusCard) SetError(errMsg string) {
	s.subtitleText.Text = "Update failed"
	s.subtitleText.Color = theme.ErrorColor
	s.subtitleText.Refresh()
	s.progressBar.Hide()
	s.container.Refresh()
}

// SetRestoring는 복원 진행 상태를 표시합니다
func (s *StatusCard) SetRestoring(msg string) {
	s.subtitleText.Text = "Restoring previous version"
	s.subtitleText.Color = theme.WarningColor
	s.progressBar.Show()
	s.container.Refresh()
}

// SetProgress는 진행률을 업데이트합니다
func (s *StatusCard) SetProgress(current, total int64) {
	// 진행률 계산 (0.0 ~ 1.0)
	progress := float64(current) / float64(total)
	s.progressBar.SetValue(progress)
}

// SetRestoreComplete는 복원 완료 상태를 표시합니다
func (s *StatusCard) SetRestoreComplete() {
	s.subtitleText.Text = "Restoration completed"
	s.subtitleText.Color = theme.SuccessColor
	s.progressBar.Hide()
	s.container.Refresh()
}
