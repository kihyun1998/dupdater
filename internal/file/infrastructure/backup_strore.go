package infrastructure

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/kihyun1998/dupdater/internal/file/domain/entity"
	"github.com/kihyun1998/dupdater/internal/file/domain/repository"
)

// backupStore는 백업 정보의 저장소를 구현하는 구조체입니다
type backupStore struct {
	mu     sync.RWMutex
	logger Logger
	config *repository.Config
}

// NewBackupStore는 새로운 backupStore 인스턴스를 생성합니다
func NewBackupStore(config *repository.Config, logger Logger) repository.BackupRepository {
	return &backupStore{
		config: config,
		logger: logger,
	}
}

// SaveBackup은 백업 정보를 저장합니다
func (s *backupStore) SaveBackup(backup *entity.BackupInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 메타데이터 파일 경로 설정
	metaPath := filepath.Join(s.config.BackupDir, "backup.meta")

	// 메타데이터를 JSON으로 변환
	data, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		return fmt.Errorf("백업 메타데이터 직렬화 실패: %w", err)
	}

	// 파일 저장
	if err := os.WriteFile(metaPath, data, 0644); err != nil {
		return fmt.Errorf("백업 메타데이터 저장 실패: %w", err)
	}

	return nil
}

// GetBackup은 백업 정보를 조회합니다
func (s *backupStore) GetBackup(sourceDir string) (*entity.BackupInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 메타데이터 파일 경로
	metaPath := filepath.Join(s.config.BackupDir, "backup.meta")

	// 파일 읽기
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("백업 메타데이터 읽기 실패: %w", err)
	}

	// JSON 파싱
	var backup entity.BackupInfo
	if err := json.Unmarshal(data, &backup); err != nil {
		return nil, fmt.Errorf("백업 메타데이터 파싱 실패: %w", err)
	}

	// 소스 디렉토리 검증
	if backup.SourceDir != sourceDir {
		return nil, fmt.Errorf("백업 소스 디렉토리 불일치: %s != %s", backup.SourceDir, sourceDir)
	}

	return &backup, nil
}

// DeleteBackup은 백업을 삭제합니다
func (s *backupStore) DeleteBackup(sourceDir string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 메타데이터 파일 삭제
	metaPath := filepath.Join(s.config.BackupDir, "backup.meta")
	if err := os.Remove(metaPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("백업 메타데이터 삭제 실패: %w", err)
	}

	// 백업 디렉토리 삭제
	if err := os.RemoveAll(s.config.BackupDir); err != nil {
		return fmt.Errorf("백업 디렉토리 삭제 실패: %w", err)
	}

	return nil
}
