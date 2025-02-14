package repository

import (
	"github.com/kihyun1998/dupdater/internal/i18n"
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/ui/domain/entity"
	"github.com/kihyun1998/dupdater/pkg/utils/theme"
)

// UIRepository는 UI 상태 관리를 위한 저장소 인터페이스입니다
type UIRepository interface {
	// UIManager 인터페이스와 호환되는 메서드들
	SetCurrentStep(step int)
	UpdateDetail(message string)
	ShowError(err error)
	Run()
	Close()
	SetRestoreHandler(handler func())
	ShowRestoring()
	ShowRestoreComplete()
	GetTotalSteps() int
	SetCompletionCallback(func())

	// 내부 상태 관리를 위한 추가 메서드
	GetState() *entity.UIState
	UpdateState(state *entity.UIState)
}

// Config는 UI Repository 생성에 필요한 설정입니다
type Config struct {
	AppName     string
	TotalSteps  int
	FromVersion string
	ToVersion   string
	Logger      logger.Logger
	Theme       theme.ThemeVariant
	I18n        i18n.Manager
}
