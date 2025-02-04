package usecase

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kihyun1998/dupdater/internal/hash/domain/entity"
	"github.com/kihyun1998/dupdater/internal/hash/domain/repository"
)

// HashService는 해시 계산과 검증을 담당하는 서비스입니다
type HashService struct {
	repo     repository.HashRepository
	basePath string
}

// NewHashService는 새로운 HashService를 생성합니다
func NewHashService(repo repository.HashRepository, config repository.HashConfig) *HashService {
	return &HashService{
		repo:     repo,
		basePath: config.BasePath,
	}
}

// CalculateFileHash는 파일의 해시값을 계산합니다
func (s *HashService) CalculateFileHash(filePath string) (*entity.HashInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("파일 열기 실패: %w", err)
	}
	defer file.Close()

	hash, err := s.repo.CalculateHash(file)
	if err != nil {
		return nil, fmt.Errorf("해시 계산 실패: %w", err)
	}

	relPath, err := filepath.Rel(s.basePath, filePath)
	if err != nil {
		return nil, fmt.Errorf("상대 경로 변환 실패: %w", err)
	}

	fileType := "f"
	if info, err := file.Stat(); err == nil && info.IsDir() {
		fileType = "d"
	}

	return entity.NewHashInfo(relPath, fileType, hash), nil
}

// VerifyFile은 파일의 해시값을 검증합니다
func (s *HashService) VerifyFile(filePath string, expectedHash string) *entity.VerificationResult {
	info, err := s.CalculateFileHash(filePath)
	if err != nil {
		return entity.NewVerificationResult(filePath, expectedHash, "", err)
	}

	return entity.NewVerificationResult(filePath, expectedHash, info.DataHash, nil)
}
