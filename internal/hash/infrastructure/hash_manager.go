// Package infrastructure는 해시 도메인의 실제 구현체를 제공합니다
package infrastructure

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kihyun1998/dupdater/internal/hash/domain/entity"
	"github.com/kihyun1998/dupdater/internal/hash/domain/repository"
)

// FileHashManager는 파일 시스템 기반의 해시 저장소 구현체입니다
type FileHashManager struct {
	config *repository.Config
}

// NewFileHashManager는 새로운 FileHashManager 인스턴스를 생성합니다
func NewFileHashManager(config *repository.Config) repository.HashRepository {
	return &FileHashManager{
		config: config,
	}
}

// VerifyFile은 단일 파일의 해시를 검증합니다
func (m *FileHashManager) VerifyFile(filePath string, expectedHash string) error {
	// 패닉 복구
	defer func() {
		if r := recover(); r != nil {
			m.config.Logger.Error("VerifyFile 함수에서 패닉 발생: %v", r)
		}
	}()

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("파일 열기 실패: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("해시 계산 실패: %w", err)
	}

	calculatedHash := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	if calculatedHash != expectedHash {
		return fmt.Errorf("해시 불일치 - 파일: %s", filePath)
	}

	return nil
}

// VerifyUpdateFile은 업데이트 ZIP 파일의 해시를 검증합니다
func (m *FileHashManager) VerifyUpdateFile(filePath string) error {
	// 패닉 복구
	defer func() {
		if r := recover(); r != nil {
			m.config.Logger.Error("VerifyUpdateFile 함수에서 패닉 발생: %v", r)
		}
	}()

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

	// 해시 데이터는 파일 끝에 위치 (크기(4바이트) + 해시값)
	hashData := make([]byte, 36) // 4바이트 크기 + 32바이트 해시
	if _, err := file.ReadAt(hashData, fileInfo.Size()-36); err != nil {
		return fmt.Errorf("해시 데이터 읽기 실패: %w", err)
	}

	hashLength := binary.LittleEndian.Uint32(hashData[:4])
	storedHash := hashData[4 : 4+hashLength]

	// 파일 내용 해시 계산
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("파일 포인터 이동 실패: %w", err)
	}

	hash := sha256.New()
	if _, err := io.CopyN(hash, file, fileInfo.Size()-36); err != nil {
		return fmt.Errorf("해시 계산 실패: %w", err)
	}

	calculatedHash := hash.Sum(nil)

	// 해시 비교
	if !compareHashes(calculatedHash, storedHash) {
		return fmt.Errorf("업데이트 파일 해시 검증 실패")
	}

	return nil
}

// VerifyHashSum은 hash_sum.txt 파일의 내용을 검증합니다
func (m *FileHashManager) VerifyHashSum() error {
	// 패닉 복구
	defer func() {
		if r := recover(); r != nil {
			m.config.Logger.Error("VerifyHashSum 함수에서 패닉 발생: %v", r)
		}
	}()

	sumFilePath := filepath.Join(m.config.CurrentDir, "hash_sum.txt")
	file, err := os.Open(sumFilePath)
	if err != nil {
		return fmt.Errorf("hash_sum.txt 파일 열기 실패: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fileHash, err := m.parseHashLine(line)
		if err != nil {
			return fmt.Errorf("해시 라인 파싱 실패: %w", err)
		}

		if err := m.VerifyFile(fileHash.Path, fileHash.HashSum); err != nil {
			return fmt.Errorf("파일 검증 실패 (%s): %w", fileHash.Path, err)
		}
	}

	return scanner.Err()
}

// GetFileHash는 파일의 해시 정보를 조회합니다
func (m *FileHashManager) GetFileHash(filePath string) (*entity.FileHash, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("파일 열기 실패: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("파일 정보 가져오기 실패: %w", err)
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, fmt.Errorf("해시 계산 실패: %w", err)
	}

	fileType := "f"
	if fileInfo.IsDir() {
		fileType = "d"
	}

	return entity.NewFileHash(
		filePath,
		base64.StdEncoding.EncodeToString(hash.Sum(nil)),
		fileType,
	), nil
}

// 내부 헬퍼 함수
func (m *FileHashManager) parseHashLine(line string) (*entity.FileHash, error) {
	parts := strings.Split(line, ";")
	if len(parts) != 3 {
		return nil, fmt.Errorf("잘못된 해시 라인 형식: %s", line)
	}

	// parts[0]: fileType, parts[1]: pathHash, parts[2]: dataHash
	filePath, err := m.getFilePathFromHash(parts[1])
	if err != nil {
		return nil, fmt.Errorf("파일 경로 찾기 실패: %w", err)
	}

	return entity.NewFileHash(filePath, parts[2], parts[0]), nil
}

// getFilePathFromHash는 해시값에 해당하는 파일의 경로를 찾습니다
func (m *FileHashManager) getFilePathFromHash(pathHash string) (string, error) {
	var matchedPath string

	err := filepath.Walk(m.config.CurrentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 디렉토리인 경우 계속 진행
		if info.IsDir() {
			return nil
		}

		// 상대 경로로 변환
		relPath, err := filepath.Rel(m.config.CurrentDir, path)
		if err != nil {
			return err
		}

		// Windows 경로를 Unix 스타일로 변환
		relPath = filepath.ToSlash(relPath)

		// 경로의 해시값 계산
		hash := sha256.Sum256([]byte(relPath))
		currentHash := base64.StdEncoding.EncodeToString(hash[:])

		// 해시값이 일치하면 경로 저장 및 검색 중단
		if currentHash == pathHash {
			matchedPath = path
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("파일 검색 실패: %w", err)
	}

	if matchedPath == "" {
		return "", fmt.Errorf("해시에 매칭되는 파일을 찾을 수 없음: %s", pathHash)
	}

	return matchedPath, nil
}

func compareHashes(hash1, hash2 []byte) bool {
	if len(hash1) != len(hash2) {
		return false
	}
	for i := range hash1 {
		if hash1[i] != hash2[i] {
			return false
		}
	}
	return true
}
