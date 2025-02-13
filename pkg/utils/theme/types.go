package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	baseTheme "fyne.io/fyne/v2/theme"
	"github.com/kihyun1998/dupdater/pkg/utils/fonts"
)

// ThemeVariant는 테마의 색상과 스타일을 정의하는 인터페이스입니다
type ThemeVariant interface {
	// 기본 색상
	BackgroundColor() color.Color
	TextColor() color.Color
	SubTextColor() color.Color

	// 강조 색상
	PrimaryColor() color.Color

	// 상태 표시 색상
	SuccessColor() color.Color
	WarningColor() color.Color
	ErrorColor() color.Color

	// 구분선 색상
	DividerColor() color.Color

	// 텍스트 크기
	FontSizeSmall() float32
	FontSizeMedium() float32
	FontSizeLarge() float32

	// 테마 모드
	IsLight() bool
}

// CustomTheme은 Fyne의 테마 인터페이스를 구현하는 커스텀 테마입니다
type CustomTheme struct {
	variant ThemeVariant
}

// CommonFontSizes는 공통으로 사용되는 폰트 크기를 정의합니다
const (
	DefaultFontSizeSmall  float32 = 14
	DefaultFontSizeMedium float32 = 16
	DefaultFontSizeLarge  float32 = 20
)

// NewCustomTheme은 새로운 CustomTheme 인스턴스를 생성합니다
func NewCustomTheme(variant ThemeVariant) *CustomTheme {
	return &CustomTheme{
		variant: variant,
	}
}

// Fyne Theme 인터페이스 구현
func (t *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case baseTheme.ColorNameBackground:
		return t.variant.BackgroundColor()
	case baseTheme.ColorNameForeground:
		return t.variant.TextColor()
	case baseTheme.ColorNamePrimary:
		return t.variant.PrimaryColor()
	case baseTheme.ColorNameError:
		return t.variant.ErrorColor()
	default:
		if t.variant.IsLight() {
			return baseTheme.LightTheme().Color(name, variant)
		}
		return baseTheme.DarkTheme().Color(name, variant)
	}
}

func (t *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	if t.variant.IsLight() {
		return baseTheme.LightTheme().Icon(name)
	}
	return baseTheme.DarkTheme().Icon(name)
}

// Font는 textStyle에 따른 폰트를 반환합니다
func (t *CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Bold {
		return fonts.PretendardBold
	}
	if style.Italic {
		return fonts.PretendardMedium
	}
	return fonts.PretendardRegular
}

func (t *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case baseTheme.SizeNameText:
		return t.variant.FontSizeMedium()
	case baseTheme.SizeNameHeadingText:
		return t.variant.FontSizeLarge()
	case baseTheme.SizeNameCaptionText:
		return t.variant.FontSizeSmall()
	default:
		if t.variant.IsLight() {
			return baseTheme.LightTheme().Size(name)
		}
		return baseTheme.DarkTheme().Size(name)
	}
}
