package entity

// DirInfo는 디렉토리 정보를 담는 도메인 엔티티입니다
type DirInfo struct {
	BackupDir  string // 백업 디렉토리 경로
	CurrentDir string // 현재 작업 디렉토리 경로
}

// NewDirInfo는 새로운 DirInfo 인스턴스를 생성합니다
func NewDirInfo(backupDir, currentDir string) *DirInfo {
	return &DirInfo{
		BackupDir:  backupDir,
		CurrentDir: currentDir,
	}
}

// GetBackupDir은 백업 디렉토리 경로를 반환합니다
func (d *DirInfo) GetBackupDir() string {
	return d.BackupDir
}

// GetCurrentDir은 현재 작업 디렉토리 경로를 반환합니다
func (d *DirInfo) GetCurrentDir() string {
	return d.CurrentDir
}
