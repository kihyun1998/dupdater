package entity

// UpdateStatus는 업데이트 진행 상태를 추적하는 엔티티입니다
type UpdateStatus struct {
	BackupCompleted  bool   // 백업 완료 여부
	RestoreCompleted bool   // 복원 완료 여부
	ServerIP         string // 서버 IP 주소
	UpdateFile       string // 업데이트 파일 경로
}

// NewUpdateStatus는 새로운 UpdateStatus 인스턴스를 생성합니다
func NewUpdateStatus() *UpdateStatus {
	return &UpdateStatus{
		BackupCompleted:  false,
		RestoreCompleted: false,
	}
}

// SetServerIP는 서버 IP를 설정합니다
func (s *UpdateStatus) SetServerIP(ip string) {
	s.ServerIP = ip
}

// SetUpdateFile은 업데이트 파일 경로를 설정합니다
func (s *UpdateStatus) SetUpdateFile(path string) {
	s.UpdateFile = path
}

// MarkBackupCompleted는 백업 완료 상태를 설정합니다
func (s *UpdateStatus) MarkBackupCompleted() {
	s.BackupCompleted = true
}

// MarkRestoreCompleted는 복원 완료 상태를 설정합니다
func (s *UpdateStatus) MarkRestoreCompleted() {
	s.RestoreCompleted = true
}

// NeedsRestore는 복원이 필요한 상태인지 확인합니다
func (s *UpdateStatus) NeedsRestore() bool {
	return s.BackupCompleted && !s.RestoreCompleted
}
