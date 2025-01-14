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
	container   *fyne.Container
	statusText  *canvas.Text
	detailText  *canvas.Text
	currentIcon *canvas.Image
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
	// 상태 아이콘 초기화
	s.currentIcon = canvas.NewImageFromResource(resourcePendingIconSvg)
	s.currentIcon.Resize(fyne.NewSize(24, 24))
	s.currentIcon.FillMode = canvas.ImageFillOriginal

	// 상태 텍스트
	s.statusText = canvas.NewText("업데이트 준비 중...", theme.TextColor)
	s.statusText.TextStyle = fyne.TextStyle{Bold: true}
	s.statusText.TextSize = theme.FontSizeMedium

	// 상세 텍스트
	s.detailText = canvas.NewText("잠시만 기다려주세요", theme.SubTextColor)
	s.detailText.TextSize = theme.FontSizeSmall

	// 텍스트 컨테이너
	textContainer := container.NewVBox(
		s.statusText,
		s.detailText,
	)

	// 메인 컨테이너 (아이콘 + 텍스트)
	s.container = container.NewHBox(
		container.NewPadded(s.currentIcon),
		container.NewPadded(textContainer),
	)
}

// UpdateStatus는 상태와 메시지를 업데이트합니다
func (s *StatusCard) UpdateStatus(status string, detail string) {
	s.statusText.Text = status
	s.detailText.Text = detail
	s.statusText.Refresh()
	s.detailText.Refresh()
}

// SetProgress는 다운로드 진행 상태를 표시합니다
func (s *StatusCard) SetProgress(current, total int64) {
	downloadedMB := float64(current) / (1024 * 1024)
	totalMB := float64(total) / (1024 * 1024)

	s.detailText.Text = fmt.Sprintf("%.1f MB / %.1f MB", downloadedMB, totalMB)
	s.detailText.Refresh()
}

// SetSuccess는 성공 상태를 표시합니다
func (s *StatusCard) SetSuccess(message string) {
	s.currentIcon.Resource = resourceCompletedIconSvg
	s.statusText.Color = theme.SuccessColor
	s.UpdateStatus("완료되었습니다", message)
	s.currentIcon.Refresh()
}

// SetError는 에러 상태를 표시합니다
func (s *StatusCard) SetError(errMsg string) {
	s.currentIcon.Resource = resourceFailedIconSvg
	s.statusText.Color = theme.ErrorColor
	s.UpdateStatus("오류가 발생했습니다", errMsg)
	s.currentIcon.Refresh()
}

// MinSize는 최소 크기를 반환합니다
func (s *StatusCard) MinSize() fyne.Size {
	return fyne.NewSize(400, 80)
}

// SetRestoring은 복구 진행 상태를 표시합니다
func (s *StatusCard) SetRestoring(message string) {
	s.currentIcon.Resource = resourceProgressIconSvg
	s.statusText.Color = theme.WarningColor
	s.UpdateStatus("파일 복원 중", message)
	s.currentIcon.Refresh()
}

// SetRestoreComplete는 복구 완료 상태를 표시합니다
func (s *StatusCard) SetRestoreComplete() {
	s.currentIcon.Resource = resourceCompletedIconSvg
	s.statusText.Color = theme.SuccessColor
	s.UpdateStatus("복원 완료", "파일이 성공적으로 복원되었습니다")
	s.currentIcon.Refresh()
}
