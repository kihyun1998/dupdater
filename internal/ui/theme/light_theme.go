// internal/ui/theme/light_theme.go

package theme

import "image/color"

// LightTheme는 라이트 모드 테마를 구현합니다
type LightTheme struct{}

// NewLightTheme는 새로운 라이트 테마 인스턴스를 생성합니다
func NewLightTheme() *LightTheme {
	return &LightTheme{}
}

// 기본 색상
func (t *LightTheme) BackgroundColor() color.Color {
	return color.NRGBA{R: 255, G: 255, B: 255, A: 255} // #FFFFFF
}

func (t *LightTheme) TextColor() color.Color {
	return color.NRGBA{R: 17, G: 24, B: 39, A: 255} // #111827
}

func (t *LightTheme) SubTextColor() color.Color {
	return color.NRGBA{R: 107, G: 114, B: 128, A: 255} // #6B7280
}

// 강조 색상
func (t *LightTheme) PrimaryColor() color.Color {
	return color.NRGBA{R: 107, G: 105, B: 232, A: 255} // #6B69E8
}

// 상태 표시 색상
func (t *LightTheme) SuccessColor() color.Color {
	return color.NRGBA{R: 74, G: 222, B: 128, A: 255} // #4ADE80
}

func (t *LightTheme) WarningColor() color.Color {
	return color.NRGBA{R: 251, G: 191, B: 36, A: 255} // #FBBF24
}

func (t *LightTheme) ErrorColor() color.Color {
	return color.NRGBA{R: 248, G: 113, B: 113, A: 255} // #F87171
}

// 구분선 색상
func (t *LightTheme) DividerColor() color.Color {
	return color.NRGBA{R: 229, G: 231, B: 235, A: 255} // #E5E7EB
}

// 폰트 크기
func (t *LightTheme) FontSizeSmall() float32 {
	return DefaultFontSizeSmall
}

func (t *LightTheme) FontSizeMedium() float32 {
	return DefaultFontSizeMedium
}

func (t *LightTheme) FontSizeLarge() float32 {
	return DefaultFontSizeLarge
}

// 테마 모드
func (t *LightTheme) IsLight() bool {
	return true
}
