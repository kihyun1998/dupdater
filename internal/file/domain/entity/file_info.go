// Package entity는 파일 도메인의 핵심 개념을 정의합니다
package entity

import (
	"os"
	"time"
)

// FileInfo는 파일의 메타데이터 정보를 담는 도메인 엔티티입니다
type FileInfo struct {
	Path        string      // 파일 경로
	Name        string      // 파일 이름
	Size        int64       // 파일 크기
	ModTime     time.Time   // 수정 시간
	IsDirectory bool        // 디렉토리 여부
	Mode        os.FileMode // 파일 권한
}

// NewFileInfo는 새로운 FileInfo 인스턴스를 생성합니다
func NewFileInfo(info os.FileInfo, path string) *FileInfo {
	return &FileInfo{
		Path:        path,
		Name:        info.Name(),
		Size:        info.Size(),
		ModTime:     info.ModTime(),
		IsDirectory: info.IsDir(),
		Mode:        info.Mode(),
	}
}

// BackupInfo는 백업 작업에 관련된 정보를 담는 도메인 엔티티입니다
type BackupInfo struct {
	SourcePath      string    // 원본 경로
	BackupPath      string    // 백업 경로
	BackupTime      time.Time // 백업 시간
	IsBackupSuccess bool      // 백업 성공 여부
}

// NewBackupInfo는 새로운 BackupInfo 인스턴스를 생성합니다
func NewBackupInfo(sourcePath, backupPath string) *BackupInfo {
	return &BackupInfo{
		SourcePath: sourcePath,
		BackupPath: backupPath,
		BackupTime: time.Now(),
	}
}

// SetBackupResult는 백업 결과를 설정합니다
func (b *BackupInfo) SetBackupResult(success bool) {
	b.IsBackupSuccess = success
}
