package theme

import "image/color"

// Colors는 앱에서 사용하는 색상 테마를 정의합니다
var (
	// 기본 색상
	BackgroundColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255} // 흰색
	TextColor       = color.NRGBA{R: 17, G: 24, B: 39, A: 255}    // 진한 회색
	SubTextColor    = color.NRGBA{R: 107, G: 114, B: 128, A: 255} // 중간 회색

	// 상태 표시 색상
	ProgressColor = color.NRGBA{R: 59, G: 130, B: 246, A: 255}  // 파란색
	SuccessColor  = color.NRGBA{R: 74, G: 222, B: 128, A: 255}  // 초록색
	WarningColor  = color.NRGBA{R: 251, G: 191, B: 36, A: 255}  // 노란색
	ErrorColor    = color.NRGBA{R: 248, G: 113, B: 113, A: 255} // 빨간색
)

// FontSize는 텍스트 크기를 정의합니다
const (
	FontSizeSmall  = 14 // 부가 설명용
	FontSizeMedium = 16 // 일반 텍스트용
	FontSizeLarge  = 20 // 제목용
)
