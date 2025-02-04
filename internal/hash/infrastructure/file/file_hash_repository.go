package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/kihyun1998/dupdater/internal/hash/domain/entity"
	"github.com/kihyun1998/dupdater/internal/hash/domain/repository"
	"github.com/kihyun1998/dupdater/internal/hash/infrastructure/hasher"
)

// FileHashRepository는 파일 기반 해시 저장소를 구현합니다
type FileHashRepository struct {
	hasher    *hasher.SHA256Hasher
	hashFile  string
	hashInfos map[string]*entity.HashInfo
	mu        sync.RWMutex
}

// NewFileHashRepository는 새로운 FileHashRepository를 생성합니다
func NewFileHashRepository(config repository.HashConfig) (*FileHashRepository, error) {
	repo := &FileHashRepository{
		hasher:    hasher.NewSHA256Hasher(),
		hashFile:  filepath.Join(config.BasePath, config.HashSumFile),
		hashInfos: make(map[string]*entity.HashInfo),
	}

	if err := repo.loadHashFile(); err != nil {
		return nil, err
	}

	return repo, nil
}

// CalculateHash implements repository.HashRepository
func (r *FileHashRepository) CalculateHash(reader io.Reader) ([]byte, error) {
	return r.hasher.CalculateHash(reader)
}

// VerifyHash implements repository.HashRepository
func (r *FileHashRepository) VerifyHash(hash1, hash2 []byte) bool {
	return r.hasher.VerifyHash(hash1, hash2)
}

// SaveHashInfo implements repository.HashRepository
func (r *FileHashRepository) SaveHashInfo(info *entity.HashInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.hashInfos[info.PathHash] = info

	// 파일에 저장
	file, err := os.OpenFile(r.hashFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("해시 파일 열기 실패: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, hashInfo := range r.hashInfos {
		if _, err := writer.WriteString(hashInfo.String() + "\n"); err != nil {
			return fmt.Errorf("해시 정보 쓰기 실패: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("버퍼 플러시 실패: %w", err)
	}

	return nil
}

// LoadHashInfo implements repository.HashRepository
func (r *FileHashRepository) LoadHashInfo(pathHash string) (*entity.HashInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if info, exists := r.hashInfos[pathHash]; exists {
		return info, nil
	}
	return nil, fmt.Errorf("해시 정보를 찾을 수 없음: %s", pathHash)
}

// ListHashInfos implements repository.HashRepository
func (r *FileHashRepository) ListHashInfos() ([]*entity.HashInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	hashInfos := make([]*entity.HashInfo, 0, len(r.hashInfos))
	for _, info := range r.hashInfos {
		hashInfos = append(hashInfos, info)
	}
	return hashInfos, nil
}

// loadHashFile은 해시 파일을 로드합니다
func (r *FileHashRepository) loadHashFile() error {
	file, err := os.Open(r.hashFile)
	if os.IsNotExist(err) {
		return nil // 파일이 없으면 무시
	}
	if err != nil {
		return fmt.Errorf("해시 파일 열기 실패: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")
		if len(parts) != 3 {
			continue // 잘못된 형식은 무시
		}

		info := &entity.HashInfo{
			FileType: parts[0],
			PathHash: parts[1],
			DataHash: parts[2],
		}
		r.hashInfos[info.PathHash] = info
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("해시 파일 읽기 실패: %w", err)
	}

	return nil
}
