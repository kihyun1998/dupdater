package entity

import (
	"time"
)

// BackupInfo는 백업 정보를 나타내는 구조체입니다
type BackupInfo struct {
	// 원본 디렉토리 경로
	SourceDir string

	// 백업 디렉토리 경로
	BackupDir string

	// 백업 시작 시간
	StartTime time.Time

	// 백업 완료 시간
	EndTime time.Time

	// 백업된 파일 목록
	Files []*FileInfo

	// 백업 상태
	Status BackupStatus
}

// BackupStatus는 백업 상태를 정의합니다
type BackupStatus int

const (
	// StatusPending는 대기 상태를 나타냅니다
	StatusPending BackupStatus = iota

	// StatusInProgress는 진행 중 상태를 나타냅니다
	StatusInProgress

	// StatusCompleted는 완료 상태를 나타냅니다
	StatusCompleted

	// StatusFailed는 실패 상태를 나타냅니다
	StatusFailed
)

// NewBackupInfo는 새로운 BackupInfo 인스턴스를 생성합니다
func NewBackupInfo(sourceDir, backupDir string) *BackupInfo {
	return &BackupInfo{
		SourceDir: sourceDir,
		BackupDir: backupDir,
		StartTime: time.Now(),
		Status:    StatusPending,
		Files:     make([]*FileInfo, 0),
	}
}

// AddFile은 백업 파일 목록에 파일을 추가합니다
func (b *BackupInfo) AddFile(file *FileInfo) {
	b.Files = append(b.Files, file)
}

// Complete는 백업 완료를 표시합니다
func (b *BackupInfo) Complete() {
	b.Status = StatusCompleted
	b.EndTime = time.Now()
}

// Fail은 백업 실패를 표시합니다
func (b *BackupInfo) Fail() {
	b.Status = StatusFailed
	b.EndTime = time.Now()
}

// Duration은 백업에 소요된 시간을 반환합니다
func (b *BackupInfo) Duration() time.Duration {
	if b.EndTime.IsZero() {
		return time.Since(b.StartTime)
	}
	return b.EndTime.Sub(b.StartTime)
}
