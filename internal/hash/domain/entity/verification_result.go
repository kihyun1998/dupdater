package entity

import "fmt"

// VerificationResult는 해시 검증 결과를 나타냅니다
type VerificationResult struct {
	IsValid      bool   // 검증 결과
	FilePath     string // 검증된 파일 경로
	ExpectedHash string // 기대하는 해시값
	ActualHash   string // 실제 해시값
	Error        error  // 발생한 오류
}

// NewVerificationResult는 새로운 VerificationResult를 생성합니다
func NewVerificationResult(path string, expected, actual string, err error) *VerificationResult {
	return &VerificationResult{
		IsValid:      err == nil && expected == actual,
		FilePath:     path,
		ExpectedHash: expected,
		ActualHash:   actual,
		Error:        err,
	}
}

// GetError는 검증 결과에 대한 오류 메시지를 반환합니다
func (r *VerificationResult) GetError() error {
	if r.Error != nil {
		return r.Error
	}
	if !r.IsValid {
		return fmt.Errorf("해시 불일치 - 파일: %s (예상: %s, 실제: %s)",
			r.FilePath, r.ExpectedHash, r.ActualHash)
	}
	return nil
}
