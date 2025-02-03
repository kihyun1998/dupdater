// internal/ui/theme/theme.go

package theme

import (
	"fmt"
	"image/color"
	"sync"
)

var (
	// currentTheme은 현재 사용 중인 테마입니다
	currentTheme Theme

	// themeMutex는 테마 변경 시 동시성을 제어합니다
	themeMutex sync.RWMutex
)

// InitTheme은 앱의 테마를 초기화합니다
func InitTheme(mode string) error {
	themeMutex.Lock()
	defer themeMutex.Unlock()

	switch mode {
	case "light":
		currentTheme = NewLightTheme()
	case "dark":
		currentTheme = NewDarkTheme()
	default:
		return fmt.Errorf("지원하지 않는 테마 모드: %s", mode)
	}

	return nil
}

// GetCurrentTheme은 현재 테마를 반환합니다
func GetCurrentTheme() Theme {
	themeMutex.RLock()
	defer themeMutex.RUnlock()

	if currentTheme == nil {
		// 기본값으로 라이트 테마 사용
		currentTheme = NewLightTheme()
	}

	return currentTheme
}

// SetTheme은 새로운 테마를 설정합니다
func SetTheme(theme Theme) {
	themeMutex.Lock()
	defer themeMutex.Unlock()

	currentTheme = theme
}

// 테마 모드 변경을 위한 헬퍼 함수들
func ToggleTheme() {
	themeMutex.Lock()
	defer themeMutex.Unlock()

	if currentTheme.IsLight() {
		currentTheme = NewDarkTheme()
	} else {
		currentTheme = NewLightTheme()
	}
}

// IsLightMode는 현재 라이트 모드인지 확인합니다
func IsLightMode() bool {
	return GetCurrentTheme().IsLight()
}

// IsDarkMode는 현재 다크 모드인지 확인합니다
func IsDarkMode() bool {
	return !GetCurrentTheme().IsLight()
}

// 테마 색상을 직접 가져오는 유틸리티 함수들
func GetBackgroundColor() color.Color {
	return GetCurrentTheme().BackgroundColor()
}

func GetTextColor() color.Color {
	return GetCurrentTheme().TextColor()
}

func GetPrimaryColor() color.Color {
	return GetCurrentTheme().PrimaryColor()
}

func GetSubTextColor() color.Color {
	return GetCurrentTheme().SubTextColor()
}

func GetSuccessColor() color.Color {
	return GetCurrentTheme().SuccessColor()
}

func GetWarningColor() color.Color {
	return GetCurrentTheme().WarningColor()
}

func GetErrorColor() color.Color {
	return GetCurrentTheme().ErrorColor()
}

func GetDividerColor() color.Color {
	return GetCurrentTheme().DividerColor()
}
