// Package usecase는 파일 관리의 비즈니스 로직을 구현합니다
package usecase

import (
	"fmt"

	"github.com/kihyun1998/dupdater/internal/file/domain/repository"
	"github.com/kihyun1998/dupdater/internal/logger"
)

// FileService는 파일 관리의 비즈니스 로직을 구현합니다
type FileService struct {
	repo   repository.FileRepository
	logger logger.Logger
}

// NewFileService는 새로운 FileService 인스턴스를 생성합니다
func NewFileService(repo repository.FileRepository, logger logger.Logger) *FileService {
	return &FileService{
		repo:   repo,
		logger: logger,
	}
}

// Backup은 현재 디렉토리의 파일들을 백업합니다
func (s *FileService) Backup() error {
	s.logger.Info("파일 백업 시작")
	if err := s.repo.Backup(); err != nil {
		s.logger.Error("파일 백업 실패: %v", err)
		return fmt.Errorf("파일 백업 실패: %w", err)
	}
	s.logger.Info("파일 백업 완료")
	return nil
}

// Restore는 백업된 파일들을 복원합니다
func (s *FileService) Restore() error {
	s.logger.Info("파일 복원 시작")
	if err := s.repo.Restore(); err != nil {
		s.logger.Error("파일 복원 실패: %v", err)
		return fmt.Errorf("파일 복원 실패: %w", err)
	}
	s.logger.Info("파일 복원 완료")
	return nil
}

// ExtractZip은 ZIP 파일을 압축 해제합니다
func (s *FileService) ExtractZip(zipFile string) error {
	s.logger.Info("ZIP 파일 압축 해제 시작: %s", zipFile)
	if err := s.repo.ExtractZip(zipFile); err != nil {
		s.logger.Error("ZIP 파일 압축 해제 실패: %v", err)
		return fmt.Errorf("ZIP 파일 압축 해제 실패: %w", err)
	}
	s.logger.Info("ZIP 파일 압축 해제 완료")
	return nil
}

// DeleteFile은 지정된 파일을 삭제합니다
func (s *FileService) DeleteFile(path string) error {
	s.logger.Info("파일 삭제 시작: %s", path)
	if err := s.repo.DeleteFile(path); err != nil {
		s.logger.Error("파일 삭제 실패: %v", err)
		return fmt.Errorf("파일 삭제 실패: %w", err)
	}
	s.logger.Info("파일 삭제 완료")
	return nil
}

// GetBackupDir은 백업 디렉토리 경로를 반환합니다
func (s *FileService) GetBackupDir() string {
	return s.repo.GetBackupDir()
}
