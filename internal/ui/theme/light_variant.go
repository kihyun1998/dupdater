package theme

import "image/color"

// LightVariant는 라이트 모드 테마를 구현합니다
type LightVariant struct{}

// NewLightVariant는 새로운 라이트 테마 인스턴스를 생성합니다
func NewLightVariant() *LightVariant {
	return &LightVariant{}
}

// 기본 색상
func (t *LightVariant) BackgroundColor() color.Color {
	return color.NRGBA{R: 255, G: 255, B: 255, A: 255} // #FFFFFF
}

func (t *LightVariant) TextColor() color.Color {
	return color.NRGBA{R: 17, G: 24, B: 39, A: 255} // #111827
}

func (t *LightVariant) SubTextColor() color.Color {
	return color.NRGBA{R: 107, G: 114, B: 128, A: 255} // #6B7280
}

// 강조 색상
func (t *LightVariant) PrimaryColor() color.Color {
	return color.NRGBA{R: 107, G: 105, B: 232, A: 255} // #6B69E8
}

// 상태 표시 색상
func (t *LightVariant) SuccessColor() color.Color {
	return color.NRGBA{R: 74, G: 222, B: 128, A: 255} // #4ADE80
}

func (t *LightVariant) WarningColor() color.Color {
	return color.NRGBA{R: 251, G: 191, B: 36, A: 255} // #FBBF24
}

func (t *LightVariant) ErrorColor() color.Color {
	return color.NRGBA{R: 248, G: 113, B: 113, A: 255} // #F87171
}

// 구분선 색상
func (t *LightVariant) DividerColor() color.Color {
	return color.NRGBA{R: 229, G: 231, B: 235, A: 255} // #E5E7EB
}

// 폰트 크기
func (t *LightVariant) FontSizeSmall() float32 {
	return DefaultFontSizeSmall
}

func (t *LightVariant) FontSizeMedium() float32 {
	return DefaultFontSizeMedium
}

func (t *LightVariant) FontSizeLarge() float32 {
	return DefaultFontSizeLarge
}

// 테마 모드
func (t *LightVariant) IsLight() bool {
	return true
}
