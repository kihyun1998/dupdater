package hash

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
)

// HashManager는 해시 관리를 담당하는 구조체
type HashManager struct {
	logger Logger
}

// Logger는 로깅을 위한 인터페이스
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}

// NewHashManager는 새로운 HashManager 인스턴스를 생성
func NewHashManager(logger Logger) *HashManager {
	return &HashManager{
		logger: logger,
	}
}

// VerifyFile은 단일 파일의 해시를 검증
func (hm *HashManager) VerifyFile(filePath string) error {
	hm.logger.Info("파일 해시 검증 시작: %s", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("파일 열기 실패: %w", err)
	}
	defer file.Close()

	// 파일 정보 획득
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("파일 정보 획득 실패: %w", err)
	}

	// 파일 끝에서 해시 데이터 읽기
	hashData := make([]byte, 36) // 4바이트(길이) + 32바이트(SHA-256 해시)
	if _, err := file.ReadAt(hashData, fileInfo.Size()-36); err != nil {
		return fmt.Errorf("해시 데이터 읽기 실패: %w", err)
	}

	// 해시 길이와 해시값 분리
	hashLength := binary.LittleEndian.Uint32(hashData[:4])
	storedHash := hashData[4 : 4+hashLength]

	// 파일 내용의 해시 계산
	hash := sha256.New()
	if _, err := io.CopyN(hash, file, fileInfo.Size()-36); err != nil {
		return fmt.Errorf("파일 해시 계산 실패: %w", err)
	}

	calculatedHash := hash.Sum(nil)

	// 해시 비교
	if !compareHashes(calculatedHash, storedHash) {
		return fmt.Errorf("해시 불일치")
	}

	hm.logger.Info("파일 해시 검증 완료: %s", filePath)
	return nil
}

// VerifyHashSum은 hash_sum.txt 파일을 기반으로 모든 파일의 해시를 검증
func (hm *HashManager) VerifyHashSum(baseDir string) error {
	hm.logger.Info("해시섬 검증 시작")

	sumFilePath := filepath.Join(baseDir, "hash_sum.txt")
	file, err := os.Open(sumFilePath)
	if err != nil {
		return fmt.Errorf("해시섬 파일 열기 실패: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if err := hm.verifyHashSumLine(baseDir, line); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("해시섬 파일 읽기 실패: %w", err)
	}

	hm.logger.Info("해시섬 검증 완료")
	return nil
}

// verifyHashSumLine은 hash_sum.txt의 각 라인을 검증
func (hm *HashManager) verifyHashSumLine(baseDir, line string) error {
	parts := strings.Split(line, ";")
	if len(parts) != 3 {
		return fmt.Errorf("잘못된 해시섬 라인 형식: %s", line)
	}

	// fileType := parts[0]
	pathHash := parts[1]
	dataHash := parts[2]

	// 파일 경로 찾기
	filePath, err := hm.findFileByPathHash(baseDir, pathHash)
	if err != nil {
		return fmt.Errorf("파일 찾기 실패: %w", err)
	}

	// 파일 해시 검증
	if err := hm.verifyFileHash(filePath, dataHash); err != nil {
		return fmt.Errorf("파일 해시 검증 실패 (%s): %w", filePath, err)
	}

	hm.logger.Info("파일 해시 검증 완료: %s", filePath)
	return nil
}

// findFileByPathHash는 경로 해시로 실제 파일을 찾음
func (hm *HashManager) findFileByPathHash(baseDir, targetHash string) (string, error) {
	var matchedPath string

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return fmt.Errorf("상대 경로 계산 실패: %w", err)
		}

		relPath = filepath.ToSlash(relPath)
		hash, err := calculatePathHash(relPath)
		if err != nil {
			return err
		}

		if hash == targetHash {
			matchedPath = path
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("파일 검색 실패: %w", err)
	}

	if matchedPath == "" {
		return "", fmt.Errorf("해시에 해당하는 파일을 찾을 수 없음: %s", targetHash)
	}

	return matchedPath, nil
}

// verifyFileHash는 파일의 해시를 검증
func (hm *HashManager) verifyFileHash(filePath, expectedHash string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	actualHash := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	if actualHash != expectedHash {
		return fmt.Errorf("해시 불일치: %s", filePath)
	}

	return nil
}

// calculatePathHash는 파일 경로의 해시를 계산
func calculatePathHash(relPath string) (string, error) {
	hash := sha256.Sum256([]byte(relPath))
	return base64.StdEncoding.EncodeToString(hash[:]), nil
}

// compareHashes는 두 해시값을 비교
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
