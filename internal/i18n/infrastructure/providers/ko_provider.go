package providers

import "github.com/kihyun1998/dupdater/internal/i18n/domain/entity"

// KoreanProvider는 한국어 메시지를 제공하는 구현체입니다
type KoreanProvider struct{}

// NewKoreanProvider는 새로운 KoreanProvider 인스턴스를 생성합니다
func NewKoreanProvider() entity.MessageProvider {
	return &KoreanProvider{}
}

// GetLanguage는 제공하는 언어 코드를 반환합니다
func (p *KoreanProvider) GetLanguage() entity.Language {
	return entity.Korean
}

// GetMessages는 한국어 메시지 맵을 반환합니다
func (p *KoreanProvider) GetMessages() map[string]string {
	return map[string]string{
		// 업데이트 상태 메시지
		"update.status.checking":     "앱 상태를 확인하고 있습니다...",
		"update.status.getting_info": "업데이트 정보를 확인하고 있습니다...",
		"update.status.preparing":    "업데이트 파일을 준비하고 있습니다...",
		"update.status.downloading":  "업데이트 파일을 다운로드하고 있습니다...",
		"update.status.verifying":    "업데이트 파일을 검증하고 있습니다...",
		"update.status.installing":   "업데이트를 설치하고 있습니다...",
		"update.status.finalizing":   "설치를 확인하고 있습니다...",
		"update.status.completed":    "업데이트가 완료되었습니다.",

		// 타이틀 및 헤더
		"update.title":             "업데이트",
		"update.header.new_update": "새로운 업데이트가 있습니다",
		"update.header.version":    "버전",

		// 진행 상태
		"update.progress.downloading": "다운로드 중...",
		"update.progress.percentage":  "%d%%",

		// 에러 메시지
		"update.error.generic":      "업데이트 중 오류가 발생했습니다",
		"update.error.connection":   "서버 연결에 실패했습니다",
		"update.error.download":     "파일 다운로드에 실패했습니다",
		"update.error.verification": "파일 검증에 실패했습니다",

		// 복원 관련
		"update.restore.in_progress": "이전 버전으로 복원 중...",
		"update.restore.completed":   "복원이 완료되었습니다",

		// 기타
		"update.info.restart":    "업데이트가 완료되면 자동으로 앱이 다시 시작됩니다.",
		"update.button.minimize": "최소화",
		"update.button.close":    "닫기",
	}
}
