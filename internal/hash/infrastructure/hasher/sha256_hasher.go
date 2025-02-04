package hasher

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"io"
)

// SHA256Hasher는 SHA256 해시 알고리즘을 구현합니다
type SHA256Hasher struct{}

// NewSHA256Hasher는 새로운 SHA256Hasher를 생성합니다
func NewSHA256Hasher() *SHA256Hasher {
	return &SHA256Hasher{}
}

// CalculateHash implements repository.HashRepository
func (h *SHA256Hasher) CalculateHash(reader io.Reader) ([]byte, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return nil, fmt.Errorf("해시 계산 실패: %w", err)
	}
	return hash.Sum(nil), nil
}

// VerifyHash implements repository.HashRepository
func (h *SHA256Hasher) VerifyHash(hash1, hash2 []byte) bool {
	if len(hash1) != len(hash2) {
		return false
	}
	return subtle.ConstantTimeCompare(hash1, hash2) == 1
}
