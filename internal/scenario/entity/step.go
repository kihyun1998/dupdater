package entity

// Step은 업데이트 프로세스의 각 단계를 정의합니다
type Step struct {
	Index       int    // 단계 인덱스 (0-7)
	Name        string // 단계 이름
	MessageKey  string // i18n 메시지 키
	NeedRestore bool   // 실패시 복원 필요 여부
}

// NewStep은 새로운 단계 정보를 생성합니다
func NewStep(index int, name string, messageKey string, needRestore bool) *Step {
	return &Step{
		Index:       index,
		Name:        name,
		MessageKey:  messageKey,
		NeedRestore: needRestore,
	}
}

// Steps는 전체 업데이트 단계 정보를 제공합니다
var Steps = []*Step{
	NewStep(0, "AppCheck", "update.status.checking", false),
	NewStep(1, "ServerInfo", "update.status.getting_info", false),
	NewStep(2, "Backup", "update.status.preparing", true),
	NewStep(3, "Download", "update.status.downloading", true),
	NewStep(4, "Verify", "update.status.verifying", true),
	NewStep(5, "Install", "update.status.installing", true),
	NewStep(6, "FinalVerify", "update.status.finalizing", true),
	NewStep(7, "Restart", "update.status.completed", true),
}

// GetStep은 인덱스에 해당하는 단계 정보를 반환합니다
func GetStep(index int) *Step {
	if index < 0 || index >= len(Steps) {
		return nil
	}
	return Steps[index]
}
