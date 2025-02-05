// Package entity는 파일 도메인의 핵심 엔티티를 정의합니다
package entity

import (
	"os"
	"path/filepath"
	"time"
)

// FileInfo는 파일 정보를 나타내는 구조체입니다
type FileInfo struct {
	// 파일 경로
	Path string

	// 파일 이름
	Name string

	// 파일 크기 (바이트)
	Size int64

	// 파일 모드 (권한)
	Mode os.FileMode

	// 수정 시간
	ModTime time.Time

	// 파일 여부 (디렉토리가 아닌 경우 true)
	IsFile bool
}

// NewFileInfo는 새로운 FileInfo 인스턴스를 생성합니다
func NewFileInfo(path string, info os.FileInfo) *FileInfo {
	return &FileInfo{
		Path:    path,
		Name:    info.Name(),
		Size:    info.Size(),
		Mode:    info.Mode(),
		ModTime: info.ModTime(),
		IsFile:  !info.IsDir(),
	}
}

// GetRelativePath는 기준 디렉토리에 대한 상대 경로를 반환합니다
func (f *FileInfo) GetRelativePath(baseDir string) (string, error) {
	return filepath.Rel(baseDir, f.Path)
}

// Exists는 파일이 실제로 존재하는지 확인합니다
func (f *FileInfo) Exists() bool {
	_, err := os.Stat(f.Path)
	return err == nil
}
