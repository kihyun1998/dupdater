package theme

import "image/color"

// Theme은 앱의 테마를 정의하는 인터페이스입니다
type Theme interface {
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

// CommonFontSizes는 공통으로 사용되는 폰트 크기를 정의합니다
const (
	DefaultFontSizeSmall  float32 = 14
	DefaultFontSizeMedium float32 = 16
	DefaultFontSizeLarge  float32 = 20
)
