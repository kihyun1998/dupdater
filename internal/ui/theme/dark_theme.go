// internal/ui/theme/dark_theme.go

package theme

import "image/color"

// DarkTheme는 다크 모드 테마를 구현합니다
type DarkTheme struct{}

// NewDarkTheme는 새로운 다크 테마 인스턴스를 생성합니다
func NewDarkTheme() *DarkTheme {
	return &DarkTheme{}
}

// 기본 색상
func (t *DarkTheme) BackgroundColor() color.Color {
	return color.NRGBA{R: 12, G: 12, B: 19, A: 255} // #0C0C13
}

func (t *DarkTheme) TextColor() color.Color {
	return color.NRGBA{R: 229, G: 231, B: 235, A: 255} // #E5E7EB
}

func (t *DarkTheme) SubTextColor() color.Color {
	return color.NRGBA{R: 156, G: 163, B: 175, A: 255} // #9CA3AF
}

// 강조 색상
func (t *DarkTheme) PrimaryColor() color.Color {
	return color.NRGBA{R: 107, G: 105, B: 232, A: 255} // #6B69E8
}

// 상태 표시 색상
func (t *DarkTheme) SuccessColor() color.Color {
	return color.NRGBA{R: 34, G: 197, B: 94, A: 255} // #22C55E
}

func (t *DarkTheme) WarningColor() color.Color {
	return color.NRGBA{R: 234, G: 179, B: 8, A: 255} // #EAB308
}

func (t *DarkTheme) ErrorColor() color.Color {
	return color.NRGBA{R: 239, G: 68, B: 68, A: 255} // #EF4444
}

// 구분선 색상
func (t *DarkTheme) DividerColor() color.Color {
	return color.NRGBA{R: 31, G: 41, B: 55, A: 255} // #1F2937
}

// 폰트 크기
func (t *DarkTheme) FontSizeSmall() float32 {
	return DefaultFontSizeSmall
}

func (t *DarkTheme) FontSizeMedium() float32 {
	return DefaultFontSizeMedium
}

func (t *DarkTheme) FontSizeLarge() float32 {
	return DefaultFontSizeLarge
}

// 테마 모드
func (t *DarkTheme) IsLight() bool {
	return false
}
