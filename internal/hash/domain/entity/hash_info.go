package entity

import (
	"encoding/base64"
	"fmt"
	"path/filepath"
)

// HashInfo는 파일의 해시 정보를 나타냅니다
type HashInfo struct {
	FilePath string // 파일 경로
	FileType string // 파일 타입 (f: 파일, d: 디렉토리)
	PathHash string // 경로의 해시값
	DataHash string // 파일 내용의 해시값
}

// NewHashInfo는 새로운 HashInfo를 생성합니다
func NewHashInfo(path string, fileType string, dataHash []byte) *HashInfo {
	// 상대 경로로 변환 및 정규화
	relPath := filepath.ToSlash(path)

	return &HashInfo{
		FilePath: relPath,
		FileType: fileType,
		PathHash: base64.StdEncoding.EncodeToString([]byte(relPath)),
		DataHash: base64.StdEncoding.EncodeToString(dataHash),
	}
}

// String은 HashInfo를 문자열로 변환합니다
func (h *HashInfo) String() string {
	return fmt.Sprintf("%s;%s;%s", h.FileType, h.PathHash, h.DataHash)
}
