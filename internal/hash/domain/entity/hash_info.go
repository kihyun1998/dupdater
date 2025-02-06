// Package entity는 해시 도메인의 핵심 개념을 정의합니다
package entity

import (
	"crypto/sha256"
	"encoding/base64"
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

// HashResult는 해시 계산 결과를 담는 값 객체입니다
type HashResult struct {
	Hash []byte
}

// NewHashResult는 바이트 슬라이스로부터 새로운 HashResult를 생성합니다
func NewHashResult(hash []byte) *HashResult {
	return &HashResult{Hash: hash}
}

// CalculateHash는 데이터의 SHA-256 해시를 계산합니다
func CalculateHash(data []byte) *HashResult {
	hash := sha256.Sum256(data)
	return &HashResult{Hash: hash[:]}
}

// ToBase64는 해시값을 base64 인코딩된 문자열로 변환합니다
func (r *HashResult) ToBase64() string {
	return base64.StdEncoding.EncodeToString(r.Hash)
}

// Compare는 두 해시값을 비교합니다
func (r *HashResult) Compare(other *HashResult) bool {
	if len(r.Hash) != len(other.Hash) {
		return false
	}
	for i := range r.Hash {
		if r.Hash[i] != other.Hash[i] {
			return false
		}
	}
	return true
}
