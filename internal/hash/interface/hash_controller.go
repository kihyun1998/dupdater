package hashinterface

import (
	"github.com/kihyun1998/dupdater/internal/hash/domain/usecase"
)

// HashController는 해시 관련 작업을 제어합니다
type HashController struct {
	hashService         *usecase.HashService
	verificationService *usecase.VerificationService
}

// NewHashController는 새로운 HashController를 생성합니다
func NewHashController(
	hashService *usecase.HashService,
	verificationService *usecase.VerificationService,
) *HashController {
	return &HashController{
		hashService:         hashService,
		verificationService: verificationService,
	}
}

// VerifyFile은 단일 파일의 해시를 검증합니다
func (c *HashController) VerifyFile(filePath string, expectedHash string) error {
	result := c.hashService.VerifyFile(filePath, expectedHash)
	return result.GetError()
}

// VerifyUpdateFile은 업데이트 파일의 해시를 검증합니다
func (c *HashController) VerifyUpdateFile(filePath string) error {
	return c.verificationService.VerifyUpdateFile(filePath)
}

// VerifyHashSum은 hash_sum.txt 파일의 내용을 검증합니다
func (c *HashController) VerifyHashSum() error {
	return c.verificationService.VerifyHashSum()
}
