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

// HeaderCard는 업데이트 헤더를 표시하는 컴포넌트입니다
type HeaderCard struct {
	widget.BaseWidget
	title       string
	fromVersion string
	toVersion   string
	container   *fyne.Container
	titleText   *canvas.Text
	versionText *canvas.Text
}

// NewHeaderCard는 새로운 HeaderCard를 생성합니다
func NewHeaderCard(title, fromVersion, toVersion string) *HeaderCard {
	card := &HeaderCard{
		title:       title,
		fromVersion: fromVersion,
		toVersion:   toVersion,
	}
	card.ExtendBaseWidget(card)
	card.setupUI()
	return card
}

// CreateRenderer는 Fyne 위젯 인터페이스를 구현합니다
func (h *HeaderCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(h.container)
}

// setupUI는 UI 컴포넌트를 초기화합니다
func (h *HeaderCard) setupUI() {
	// 타이틀 텍스트
	h.titleText = canvas.NewText(h.title, theme.TextColor)
	h.titleText.TextStyle = fyne.TextStyle{Bold: true}
	h.titleText.TextSize = theme.FontSizeLarge

	// 버전 텍스트
	versionStr := fmt.Sprintf("%s → %s", h.fromVersion, h.toVersion)
	h.versionText = canvas.NewText(versionStr, theme.SubTextColor)
	h.versionText.TextSize = theme.FontSizeSmall
	h.versionText.Alignment = fyne.TextAlignTrailing

	// 구분선
	// divider := canvas.NewLine(theme.DividerColor)
	// divider.StrokeWidth = 1

	// 헤더 컨테이너 (타이틀과 버전 정보)
	headerContainer := container.NewHBox(
		h.titleText,
		layout.NewSpacer(),
		h.versionText,
	)

	// 메인 컨테이너
	h.container = container.NewVBox(
		headerContainer,
		// divider,
	)

	// 패딩 추가
	h.container = container.NewPadded(h.container)
	h.container.Resize(fyne.NewSize(400, 60))
}

// UpdateVersions는 버전 정보를 업데이트합니다
func (h *HeaderCard) UpdateVersions(fromVersion, toVersion string) {
	h.fromVersion = fromVersion
	h.toVersion = toVersion
	if h.versionText != nil {
		h.versionText.Text = fmt.Sprintf("%s → %s", fromVersion, toVersion)
		h.versionText.Refresh()
		h.container.Refresh()
	}
}

// MinSize는 최소 크기를 반환합니다
func (h *HeaderCard) MinSize() fyne.Size {
	return fyne.NewSize(400, 60)
}
