package theme

import "image/color"

// Colors는 앱에서 사용하는 색상 테마를 정의합니다
var (
	// 기본 색상
	BackgroundColor = color.NRGBA{R: 24, G: 24, B: 27, A: 255}    // gray-900
	TextColor       = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // white
	SubTextColor    = color.NRGBA{R: 156, G: 163, B: 175, A: 255} // gray-400

	// 상태 색상
	SuccessColor = color.NRGBA{R: 74, G: 222, B: 128, A: 255}  // green-400
	InfoColor    = color.NRGBA{R: 96, G: 165, B: 250, A: 255}  // blue-400
	WarningColor = color.NRGBA{R: 251, G: 146, B: 60, A: 255}  // orange-400
	ErrorColor   = color.NRGBA{R: 248, G: 113, B: 113, A: 255} // red-400

	// 컴포넌트 색상
	CardBackgroundColor = color.NRGBA{R: 31, G: 41, B: 55, A: 255} // gray-800
	DividerColor        = color.NRGBA{R: 75, G: 85, B: 99, A: 255} // gray-600
)

// Padding은 컴포넌트 간격을 정의합니다
const (
	PaddingSmall  = 8
	PaddingMedium = 16
	PaddingLarge  = 24
)

// FontSize는 텍스트 크기를 정의합니다
const (
	FontSizeSmall  = 12
	FontSizeMedium = 14
	FontSizeLarge  = 18
)
