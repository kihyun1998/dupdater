package components

import (
	"fmt"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kihyun1998/dupdater/internal/i18n"
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
	currentTheme theme.ThemeVariant
	i18n         i18n.LocaleManager

	targetProgress  float64
	currentProgress float64
	animating       bool
	ticker          *time.Ticker
	onComplete      func()
}

// NewStatusCard는 새로운 StatusCard를 생성합니다
func NewStatusCard(fromVersion, toVersion string, themeVariant theme.ThemeVariant, i18n i18n.LocaleManager) *StatusCard {
	card := &StatusCard{
		targetProgress:  0,
		currentProgress: 0,
		animating:       false,
		currentTheme:    themeVariant,
		i18n:            i18n,
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
	// 타이틀 텍스트
	s.titleText = canvas.NewText(
		s.i18n.GetMessage("update.header.new_update"),
		s.currentTheme.TextColor(),
	)
	s.titleText.TextSize = s.currentTheme.FontSizeLarge()
	s.titleText.TextStyle = fyne.TextStyle{Bold: true}

	// 버전 정보
	versionFormat := fmt.Sprintf("%s → %s", fromVersion, toVersion)
	s.versionText = canvas.NewText(versionFormat, s.currentTheme.SubTextColor())
	s.versionText.TextSize = s.currentTheme.FontSizeMedium()

	// 진행 바
	s.progressBar = widget.NewProgressBar()
	s.progressBar.Resize(fyne.NewSize(350, 20))

	// 상태 메시지
	s.subtitleText = canvas.NewText(
		s.i18n.GetMessage("update.info.restart"),
		s.currentTheme.SubTextColor(),
	)
	s.subtitleText.TextSize = s.currentTheme.FontSizeSmall()

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
	s.titleText.Text = s.i18n.GetMessage("update.error.generic")
	s.titleText.Color = s.currentTheme.ErrorColor()
	s.subtitleText.Text = errMsg
	s.subtitleText.Color = s.currentTheme.ErrorColor()
	s.progressBar.Hide()
	s.Refresh()
}

// SetRestoring는 복원 진행 상태를 표시합니다
func (s *StatusCard) SetRestoring(msg string) {
	s.titleText.Text = s.i18n.GetMessage("update.restore.in_progress")
	s.titleText.Color = s.currentTheme.WarningColor()
	s.subtitleText.Text = msg
	s.subtitleText.Color = s.currentTheme.WarningColor()
	s.progressBar.Show()
	s.Refresh()
}

// SetRestoreComplete는 복원 완료 상태를 표시합니다
func (s *StatusCard) SetRestoreComplete() {
	s.titleText.Text = s.i18n.GetMessage("update.restore.completed")
	s.titleText.Color = s.currentTheme.SuccessColor()
	s.subtitleText.Text = s.i18n.GetMessage("update.info.restart")
	s.subtitleText.Color = s.currentTheme.SuccessColor()
	s.progressBar.Hide()
	s.Refresh()
}

// UpdateStatus는 상태와 진행률을 업데이트합니다
func (s *StatusCard) UpdateStatus(progress float64, message string) {
	s.titleText.Color = s.currentTheme.TextColor()
	s.subtitleText.Text = message
	s.subtitleText.Color = s.currentTheme.SubTextColor()

	if progress >= 0 {
		s.targetProgress = progress
		s.animateProgress()
	}

	s.progressBar.Show()
	s.Refresh()
}

// SetProgress는 다운로드 진행률을 업데이트합니다
func (s *StatusCard) SetProgress(current, total int64) {
	progress := float64(current) / float64(total)
	s.progressBar.SetValue(progress)
	s.Refresh()
}

// SetCompletionCallback은 진행률 100% 도달 시 실행될 콜백을 설정합니다
func (s *StatusCard) SetCompletionCallback(callback func()) {
	s.onComplete = callback
}

// Refresh는 위젯을 새로고침합니다
func (s *StatusCard) Refresh() {
	s.container.Refresh()
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
