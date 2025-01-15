package components

import (
	"fmt"
	"math"
	"time"

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

	targetProgress  float64
	currentProgress float64
	animating       bool
	ticker          *time.Ticker
	onComplete      func() // 추가
}

// NewStatusCard는 새로운 StatusCard를 생성합니다
func NewStatusCard(fromVersion, toVersion string) *StatusCard {
	card := &StatusCard{
		targetProgress:  0,
		currentProgress: 0,
		animating:       false,
	}
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

	if progress >= 0 {
		s.targetProgress = progress
		s.animateProgress()
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

// SetCompletionCallback은 진행률 100% 도달 시 실행될 콜백을 설정합니다
func (s *StatusCard) SetCompletionCallback(callback func()) {
	s.onComplete = callback
}

// 부드러운 진행률 업데이트를 위한 메서드
func (s *StatusCard) animateProgress() {
	if s.ticker != nil {
		s.ticker.Stop()
	}

	s.animating = true
	s.ticker = time.NewTicker(16 * time.Millisecond) // 약 60FPS

	go func() {
		defer s.ticker.Stop()

		for range s.ticker.C {
			if !s.animating {
				return
			}

			diff := s.targetProgress - s.currentProgress
			step := diff * 0.1

			if math.Abs(diff) < 0.001 {
				s.currentProgress = s.targetProgress
				s.progressBar.SetValue(s.currentProgress)
				s.progressBar.Refresh()
				s.animating = false

				// 애니메이션이 100%에 도달했을 때 콜백 실행
				if s.currentProgress >= 0.999 {
					if s.onComplete != nil {
						s.onComplete()
					}
				}
				return
			}

			s.currentProgress += step
			s.progressBar.SetValue(s.currentProgress)
			s.progressBar.Refresh()
		}
	}()
}
