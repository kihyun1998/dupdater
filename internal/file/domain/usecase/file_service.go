package usecase

import (
	"fmt"
	"path/filepath"

	"github.com/kihyun1998/dupdater/internal/file/domain/entity"
	"github.com/kihyun1998/dupdater/internal/file/domain/repository"
)

// FileService는 파일 관련 비즈니스 로직을 구현합니다
type FileService struct {
	fileRepo   repository.FileRepository
	backupRepo repository.BackupRepository
	config     *repository.Config
	logger     Logger
}

// NewFileService는 새로운 FileService 인스턴스를 생성합니다
func NewFileService(
	fileRepo repository.FileRepository,
	backupRepo repository.BackupRepository,
	config *repository.Config,
	logger Logger,
) *FileService {
	return &FileService{
		fileRepo:   fileRepo,
		backupRepo: backupRepo,
		config:     config,
		logger:     logger,
	}
}

// Backup은 현재 디렉토리의 파일들을 백업합니다
func (s *FileService) Backup() error {
	// 백업 정보 생성
	backupInfo := entity.NewBackupInfo(s.config.CurrentDir, s.config.BackupDir)

	// 백업 디렉토리 생성
	if err := s.fileRepo.CreateDirectory(s.config.BackupDir); err != nil {
		return fmt.Errorf("백업 디렉토리 생성 실패: %w", err)
	}

	// 현재 디렉토리 파일 목록 조회
	files, err := s.fileRepo.ListFiles(s.config.CurrentDir)
	if err != nil {
		return fmt.Errorf("파일 목록 조회 실패: %w", err)
	}

	// 각 파일 백업
	for _, file := range files {
		destPath := filepath.Join(s.config.BackupDir, file.Name)
		if err := s.fileRepo.CopyFile(file.Path, destPath); err != nil {
			backupInfo.Fail()
			s.logger.Error("파일 백업 실패: %v", err)
			return fmt.Errorf("파일 백업 실패 (%s): %w", file.Name, err)
		}
		backupInfo.AddFile(file)
	}

	// 백업 완료 처리
	backupInfo.Complete()

	// 백업 정보 저장
	if err := s.backupRepo.SaveBackup(backupInfo); err != nil {
		return fmt.Errorf("백업 정보 저장 실패: %w", err)
	}

	return nil
}

// Restore는 백업된 파일들을 복원합니다
func (s *FileService) Restore() error {
	// 백업 정보 조회
	backupInfo, err := s.backupRepo.GetBackup(s.config.CurrentDir)
	if err != nil {
		return fmt.Errorf("백업 정보 조회 실패: %w", err)
	}

	// 현재 디렉토리 정리
	files, err := s.fileRepo.ListFiles(s.config.CurrentDir)
	if err != nil {
		return fmt.Errorf("현재 디렉토리 파일 목록 조회 실패: %w", err)
	}

	// 기존 파일 삭제
	for _, file := range files {
		if err := s.fileRepo.DeleteFile(file.Path); err != nil {
			s.logger.Error("파일 삭제 실패: %v", err)
			continue
		}
	}

	// 백업 파일 복원
	for _, file := range backupInfo.Files {
		srcPath := filepath.Join(s.config.BackupDir, file.Name)
		destPath := filepath.Join(s.config.CurrentDir, file.Name)

		if err := s.fileRepo.CopyFile(srcPath, destPath); err != nil {
			return fmt.Errorf("파일 복원 실패 (%s): %w", file.Name, err)
		}
	}

	// 백업 정보 삭제
	if err := s.backupRepo.DeleteBackup(s.config.CurrentDir); err != nil {
		s.logger.Error("백업 정보 삭제 실패: %v", err)
	}

	return nil
}

// ExtractZip은 ZIP 파일을 지정된 디렉토리에 압축 해제합니다
func (s *FileService) ExtractZip(zipFile string) error {
	if err := s.fileRepo.ExtractZip(zipFile, s.config.CurrentDir); err != nil {
		return fmt.Errorf("ZIP 파일 압축해제 실패: %w", err)
	}
	return s.fileRepo.DeleteFile(zipFile)
}

// DeleteFile은 지정된 파일을 삭제합니다
func (s *FileService) DeleteFile(path string) error {
	return s.fileRepo.DeleteFile(path)
}

// Logger는 로깅을 위한 인터페이스입니다
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}
