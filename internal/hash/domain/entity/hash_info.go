// Package entity는 해시 도메인의 핵심 개념을 정의합니다
package entity

import (
	"fmt"
)

// FileHash는 파일의 해시 정보를 담는 도메인 엔티티입니다
type FileHash struct {
	Path     string // 파일 경로
	HashSum  string // 해시 값
	FileType string // 파일 유형 (파일/디렉토리)
}

// NewFileHash는 새로운 FileHash 인스턴스를 생성합니다
func NewFileHash(path string, hashSum string, fileType string) *FileHash {
	return &FileHash{
		Path:     path,
		HashSum:  hashSum,
		FileType: fileType,
	}
}

// Validate는 해시 정보가 유효한지 검증합니다
func (f *FileHash) Validate() error {
	if f.Path == "" {
		return fmt.Errorf("파일 경로가 비어있습니다")
	}
	if f.HashSum == "" {
		return fmt.Errorf("해시값이 비어있습니다")
	}
	if f.FileType == "" {
		return fmt.Errorf("파일 유형이 비어있습니다")
	}
	return nil
}
