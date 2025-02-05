package repository

import (
	"github.com/kihyun1998/dupdater/internal/file/domain/entity"
)

// FileRepository는 파일 작업을 위한 저장소 인터페이스입니다
type FileRepository interface {
	// CreateDirectory는 디렉토리를 생성합니다
	CreateDirectory(path string) error

	// DeleteFile은 파일을 삭제합니다
	DeleteFile(path string) error

	// CopyFile은 파일을 복사합니다
	CopyFile(src, dest string) error

	// MoveFile은 파일을 이동합니다
	MoveFile(src, dest string) error

	// GetFileInfo는 파일 정보를 가져옵니다
	GetFileInfo(path string) (*entity.FileInfo, error)

	// ListFiles는 디렉토리 내의 모든 파일 목록을 반환합니다
	ListFiles(dir string) ([]*entity.FileInfo, error)

	// ExtractZip은 ZIP 파일을 압축 해제합니다
	ExtractZip(zipFile, destDir string) error
}

// BackupRepository는 백업 작업을 위한 저장소 인터페이스입니다
type BackupRepository interface {
	// SaveBackup은 백업 정보를 저장합니다
	SaveBackup(backup *entity.BackupInfo) error

	// GetBackup은 백업 정보를 조회합니다
	GetBackup(sourceDir string) (*entity.BackupInfo, error)

	// DeleteBackup은 백업을 삭제합니다
	DeleteBackup(sourceDir string) error
}

// Config는 저장소 설정을 정의합니다
type Config struct {
	// CurrentDir은 현재 작업 디렉토리입니다
	CurrentDir string

	// BackupDir은 백업 디렉토리입니다
	BackupDir string
}

// NewConfig는 새로운 저장소 설정을 생성합니다
func NewConfig(currentDir, backupDir string) *Config {
	return &Config{
		CurrentDir: currentDir,
		BackupDir:  backupDir,
	}
}
