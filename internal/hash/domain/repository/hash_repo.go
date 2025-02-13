package repository

import (
	"github.com/kihyun1998/dupdater/internal/hash/domain/entity"
	"github.com/kihyun1998/dupdater/internal/logger"
)

// HashRepository는 해시 검증을 위한 저장소 인터페이스입니다
type HashRepository interface {
	// VerifyFile은 단일 파일의 해시를 검증합니다
	VerifyFile(filePath string, expectedHash string) error

	// VerifyUpdateFile은 업데이트 파일의 해시를 검증합니다
	VerifyUpdateFile(filePath string) error

	// VerifyHashSum은 모든 파일의 해시섬을 검증합니다
	VerifyHashSum() error

	// GetFileHash는 파일의 해시 정보를 조회합니다
	GetFileHash(filePath string) (*entity.FileHash, error)
}

// Config는 저장소 설정을 정의합니다
type Config struct {
	CurrentDir  string
	HashSumPath string
	Logger      logger.Logger
}
