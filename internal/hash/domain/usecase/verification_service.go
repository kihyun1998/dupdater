package usecase

import (
	"fmt"
	"io"
	"os"

	"github.com/kihyun1998/dupdater/internal/hash/domain/repository"
)

// VerificationService는 해시 검증 관련 서비스를 제공합니다
type VerificationService struct {
	hashService *HashService
	repo        repository.HashRepository
}

// NewVerificationService는 새로운 VerificationService를 생성합니다
func NewVerificationService(hashService *HashService, repo repository.HashRepository) *VerificationService {
	return &VerificationService{
		hashService: hashService,
		repo:        repo,
	}
}

// VerifyUpdateFile은 업데이트 파일의 해시를 검증합니다
func (s *VerificationService) VerifyUpdateFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("파일 열기 실패: %w", err)
	}
	defer file.Close()

	// 파일 크기 확인
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("파일 정보 가져오기 실패: %w", err)
	}

	// 해시 데이터는 파일 끝에 위치
	hashData := make([]byte, 36) // 4바이트 크기 + 32바이트 해시
	if _, err := file.ReadAt(hashData, fileInfo.Size()-36); err != nil {
		return fmt.Errorf("해시 데이터 읽기 실패: %w", err)
	}

	// 파일 내용 해시 계산
	contentHash, err := s.repo.CalculateHash(io.NewSectionReader(file, 0, fileInfo.Size()-36))
	if err != nil {
		return fmt.Errorf("해시 계산 실패: %w", err)
	}

	if !s.repo.VerifyHash(contentHash, hashData[4:]) {
		return fmt.Errorf("업데이트 파일 해시 검증 실패")
	}

	return nil
}

// VerifyHashSum은 hash_sum.txt 파일의 내용을 검증합니다
func (s *VerificationService) VerifyHashSum() error {
	hashInfos, err := s.repo.ListHashInfos()
	if err != nil {
		return fmt.Errorf("해시 정보 목록 가져오기 실패: %w", err)
	}

	for _, info := range hashInfos {
		result := s.hashService.VerifyFile(info.FilePath, info.DataHash)
		if err := result.GetError(); err != nil {
			return err
		}
	}

	return nil
}
