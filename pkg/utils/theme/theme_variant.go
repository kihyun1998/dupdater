package theme

import (
	"fmt"
	"image/color"
	"sync"
)

var (
	// currentVariant는 현재 사용 중인 테마 변형입니다
	currentVariant ThemeVariant

	// currentTheme은 현재 사용 중인 Fyne 테마입니다
	currentTheme *CustomTheme

	// themeMutex는 테마 변경 시 동시성을 제어합니다
	themeMutex sync.RWMutex
)

// InitTheme은 앱의 테마를 초기화합니다
func InitTheme(mode string) error {
	themeMutex.Lock()
	defer themeMutex.Unlock()

	switch mode {
	case "light":
		currentVariant = NewLightVariant()
	case "dark":
		currentVariant = NewDarkVariant()
	default:
		return fmt.Errorf("지원하지 않는 테마 모드: %s", mode)
	}

	currentTheme = NewCustomTheme(currentVariant)
	return nil
}

// GetCurrentTheme은 현재 Fyne 테마를 반환합니다
func GetCurrentTheme() *CustomTheme {
	themeMutex.RLock()
	defer themeMutex.RUnlock()

	if currentTheme == nil {
		// 기본값으로 라이트 테마 사용
		currentVariant = NewLightVariant()
		currentTheme = NewCustomTheme(currentVariant)
	}

	return currentTheme
}

// GetCurrentVariant는 현재 테마 변형을 반환합니다
func GetCurrentVariant() ThemeVariant {
	themeMutex.RLock()
	defer themeMutex.RUnlock()

	if currentVariant == nil {
		currentVariant = NewLightVariant()
	}

	return currentVariant
}

// SetTheme은 새로운 테마를 설정합니다
func SetTheme(variant ThemeVariant) {
	themeMutex.Lock()
	defer themeMutex.Unlock()

	currentVariant = variant
	currentTheme = NewCustomTheme(variant)
}

// 테마 모드 변경을 위한 헬퍼 함수들
func ToggleTheme() {
	themeMutex.Lock()
	defer themeMutex.Unlock()

	if currentVariant.IsLight() {
		currentVariant = NewDarkVariant()
	} else {
		currentVariant = NewLightVariant()
	}
	currentTheme = NewCustomTheme(currentVariant)
}

// IsLightMode는 현재 라이트 모드인지 확인합니다
func IsLightMode() bool {
	return GetCurrentVariant().IsLight()
}

// IsDarkMode는 현재 다크 모드인지 확인합니다
func IsDarkMode() bool {
	return !GetCurrentVariant().IsLight()
}

// 테마 색상을 직접 가져오는 유틸리티 함수들
func GetBackgroundColor() color.Color {
	return GetCurrentVariant().BackgroundColor()
}

func GetTextColor() color.Color {
	return GetCurrentVariant().TextColor()
}

func GetPrimaryColor() color.Color {
	return GetCurrentVariant().PrimaryColor()
}

func GetSubTextColor() color.Color {
	return GetCurrentVariant().SubTextColor()
}

func GetSuccessColor() color.Color {
	return GetCurrentVariant().SuccessColor()
}

func GetWarningColor() color.Color {
	return GetCurrentVariant().WarningColor()
}

func GetErrorColor() color.Color {
	return GetCurrentVariant().ErrorColor()
}

func GetDividerColor() color.Color {
	return GetCurrentVariant().DividerColor()
}
