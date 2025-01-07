// Package hash는 파일 해시 검증을 담당하는 패키지입니다
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

// Manager는 해시 검증을 관리하는 구조체입니다
type Manager struct {
	logger     Logger // 로깅을 위한 인터페이스
	currentDir string // 현재 작업 디렉토리
}

// Config는 Manager 생성에 필요한 설정을 담는 구조체입니다
type Config struct {
	Logger     Logger
	CurrentDir string // 현재 작업 디렉토리 (기본값: ".")
}

// FileHash는 파일의 해시 정보를 담는 구조체입니다
type FileHash struct {
	FileType string // 파일 타입 (f: 파일, d: 디렉토리)
	PathHash string // 경로의 해시값
	DataHash string // 파일 내용의 해시값
}

// New는 새로운 Manager 인스턴스를 생성합니다
func New(config Config) (*Manager, error) {
	if config.CurrentDir == "" {
		config.CurrentDir = "."
	}

	return &Manager{
		logger:     config.Logger,
		currentDir: config.CurrentDir,
	}, nil
}

// VerifyFile은 단일 파일의 해시를 검증합니다
func (m *Manager) VerifyFile(filePath string, expectedHash string) error {
	defer func() {
		if r := recover(); r != nil {
			m.logger.Error("VerifyFile 함수에서 패닉 발생: %v", r)
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
func (m *Manager) VerifyUpdateFile(filePath string) error {
	defer func() {
		if r := recover(); r != nil {
			m.logger.Error("VerifyUpdateFile 함수에서 패닉 발생: %v", r)
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
	if !m.compareHashes(calculatedHash, storedHash) {
		return fmt.Errorf("업데이트 파일 해시 검증 실패")
	}

	return nil
}

// VerifyHashSum은 hash_sum.txt 파일의 내용을 검증합니다
func (m *Manager) VerifyHashSum() error {
	defer func() {
		if r := recover(); r != nil {
			m.logger.Error("VerifyHashSum 함수에서 패닉 발생: %v", r)
		}
	}()

	sumFilePath := filepath.Join(m.currentDir, "hash_sum.txt")
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

		filePath, err := m.getFilePathFromHash(fileHash.PathHash)
		if err != nil {
			return fmt.Errorf("파일 경로 찾기 실패: %w", err)
		}

		if err := m.VerifyFile(filePath, fileHash.DataHash); err != nil {
			return fmt.Errorf("파일 검증 실패 (%s): %w", filePath, err)
		}
	}

	return scanner.Err()
}

// 내부 헬퍼 함수들
// parseHashLine는 hash_sum.txt 파일의 한 줄을 파싱합니다
func (m *Manager) parseHashLine(line string) (FileHash, error) {
	parts := strings.Split(line, ";")
	if len(parts) != 3 {
		return FileHash{}, fmt.Errorf("잘못된 해시 라인 형식: %s", line)
	}

	return FileHash{
		FileType: parts[0],
		PathHash: parts[1],
		DataHash: parts[2],
	}, nil
}

// getFilePathFromHash는 해시값에 해당하는 파일의 경로를 찾습니다
func (m *Manager) getFilePathFromHash(pathHash string) (string, error) {
	var matchedPath string

	err := filepath.Walk(m.currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(m.currentDir, path)
		if err != nil {
			return err
		}

		relPath = filepath.ToSlash(relPath)
		currentHash, err := m.calculatePathHash(relPath)
		if err != nil {
			return err
		}

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

// calculatePathHash는 파일 경로의 해시값을 계산합니다
func (m *Manager) calculatePathHash(relPath string) (string, error) {
	hash := sha256.Sum256([]byte(relPath))
	return base64.StdEncoding.EncodeToString(hash[:]), nil
}

// compareHashes는 두 해시값을 비교합니다
func (m *Manager) compareHashes(hash1, hash2 []byte) bool {
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

// Logger는 로깅을 위한 인터페이스입니다
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}
