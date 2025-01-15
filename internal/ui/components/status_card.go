package components

import (
	"fmt"

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
	titleText    *canvas.Text
	versionText  *canvas.Text
	progressBar  *widget.ProgressBar
	subtitleText *canvas.Text
}

// NewStatusCard는 새로운 StatusCard를 생성합니다
func NewStatusCard(fromVersion, toVersion string) *StatusCard {
	card := &StatusCard{}
	card.ExtendBaseWidget(card)
	card.setupUI(fromVersion, toVersion)
	return card
}

// CreateRenderer는 Fyne 위젯 인터페이스를 구현합니다
func (s *StatusCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(s.container)
}

// setupUI는 UI 컴포넌트를 초기화합니다
func (s *StatusCard) setupUI(fromVersion, toVersion string) {
	// 타이틀
	s.titleText = canvas.NewText("새로운 업데이트가 있습니다", theme.TextColor)
	s.titleText.TextSize = 16
	s.titleText.TextStyle = fyne.TextStyle{Bold: true}

	// 버전 정보
	s.versionText = canvas.NewText(fmt.Sprintf("%s → %s", fromVersion, toVersion), theme.SubTextColor)
	s.versionText.TextSize = 14

	// 진행 바
	s.progressBar = widget.NewProgressBar()

	// 상태 메시지
	s.subtitleText = canvas.NewText("업데이트가 완료되면 자동으로 앱이 다시 시작됩니다.", theme.SubTextColor)
	s.subtitleText.TextSize = 12

	// 레이아웃 구성
	s.container = container.NewVBox(
		container.NewVBox(
			s.titleText,
			s.versionText,
		),
		container.NewPadded(s.progressBar),
		container.NewPadded(s.subtitleText),
	)
}

// SetError는 에러 상태를 표시합니다
func (s *StatusCard) SetError(errMsg string) {
	s.titleText.Text = "업데이트 중 오류가 발생했습니다"
	s.titleText.Color = theme.ErrorColor
	s.subtitleText.Text = errMsg
	s.subtitleText.Color = theme.ErrorColor
	s.progressBar.Hide()
	s.container.Refresh()
}

// SetRestoring는 복원 진행 상태를 표시합니다
func (s *StatusCard) SetRestoring(msg string) {
	s.titleText.Text = "이전 버전으로 복원 중"
	s.titleText.Color = theme.WarningColor
	s.subtitleText.Text = msg
	s.subtitleText.Color = theme.WarningColor
	s.progressBar.Show()
	s.container.Refresh()
}

// SetRestoreComplete는 복원 완료 상태를 표시합니다
func (s *StatusCard) SetRestoreComplete() {
	s.titleText.Text = "복원이 완료되었습니다"
	s.titleText.Color = theme.SuccessColor
	s.subtitleText.Text = "앱이 곧 다시 시작됩니다"
	s.subtitleText.Color = theme.SuccessColor
	s.progressBar.Hide()
	s.container.Refresh()
}

// UpdateStatus는 상태와 진행률을 업데이트합니다
func (s *StatusCard) UpdateStatus(progress float64, message string) {
	s.titleText.Text = "업데이트 진행 중"
	s.titleText.Color = theme.TextColor
	s.subtitleText.Text = message
	s.subtitleText.Color = theme.SubTextColor

	// 진행률 업데이트를 명시적으로 수행
	if progress >= 0 { // -1이 아닐 때만 진행률 업데이트
		s.progressBar.SetValue(progress)
		s.progressBar.Refresh() // 명시적 리프레시 추가
	}
	s.progressBar.Show()
	s.container.Refresh()
}

// SetProgress는 다운로드 진행률을 업데이트합니다
func (s *StatusCard) SetProgress(current, total int64) {
	progress := float64(current) / float64(total)
	s.progressBar.SetValue(progress)
	s.container.Refresh()
}
