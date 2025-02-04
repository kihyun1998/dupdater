package repository

import (
	"io"

	"github.com/kihyun1998/dupdater/internal/hash/domain/entity"
)

// HashRepository는 해시 저장소 인터페이스를 정의합니다
type HashRepository interface {
	// CalculateHash는 Reader로부터 해시값을 계산합니다
	CalculateHash(reader io.Reader) ([]byte, error)

	// VerifyHash는 두 해시값을 비교합니다
	VerifyHash(hash1, hash2 []byte) bool

	// SaveHashInfo는 해시 정보를 저장합니다
	SaveHashInfo(info *entity.HashInfo) error

	// LoadHashInfo는 파일 경로에 해당하는 해시 정보를 로드합니다
	LoadHashInfo(pathHash string) (*entity.HashInfo, error)

	// ListHashInfos는 저장된 모든 해시 정보를 반환합니다
	ListHashInfos() ([]*entity.HashInfo, error)
}

// HashConfig는 해시 저장소 설정을 정의합니다
type HashConfig struct {
	BasePath    string // 기본 디렉토리 경로
	HashSumFile string // 해시 정보 저장 파일
}
