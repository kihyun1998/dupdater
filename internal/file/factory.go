package file

import (
	"fmt"

	"github.com/kihyun1998/dupdater/internal/file/domain/ports"
	"github.com/kihyun1998/dupdater/internal/file/domain/repository"
	"github.com/kihyun1998/dupdater/internal/file/domain/usecase"
	"github.com/kihyun1998/dupdater/internal/file/infrastructure"
)

// Factory는 파일 도메인의 컴포넌트들을 생성하고 조립합니다
type Factory struct {
	config *ports.FileConfig
}

// NewFactory는 새로운 Factory 인스턴스를 생성합니다
func NewFactory(config *ports.FileConfig) (*Factory, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("파일 설정 검증 실패: %w", err)
	}

	return &Factory{
		config: config,
	}, nil
}

// Create는 FilePort 인터페이스를 구현하는 인스턴스를 생성합니다
func (f *Factory) Create() (ports.FilePort, error) {
	// 저장소 설정 생성
	repoConfig := repository.NewConfig(
		f.config.CurrentDir,
		f.config.BackupDir,
	)

	// 파일 저장소 생성
	fileRepo := infrastructure.NewFileManager(repoConfig, f.config.Logger)

	// 백업 저장소 생성
	backupRepo := infrastructure.NewBackupStore(repoConfig, f.config.Logger)

	// 파일 서비스 생성
	fileService := usecase.NewFileService(
		fileRepo,
		backupRepo,
		repoConfig,
		f.config.Logger,
	)

	return &fileAdapter{
		service: fileService,
		logger:  f.config.Logger,
	}, nil
}

// New는 새로운 FilePort 인스턴스를 생성하는 편의 함수입니다
func New(config ports.FileConfig) (ports.FilePort, error) {
	factory, err := NewFactory(&config)
	if err != nil {
		return nil, err
	}
	return factory.Create()
}

// fileAdapter는 FilePort 인터페이스를 구현하고 FileService를 사용합니다
type fileAdapter struct {
	service *usecase.FileService
	logger  ports.Logger
}

// Backup은 현재 디렉토리의 파일들을 백업합니다
func (a *fileAdapter) Backup() error {
	if err := a.service.Backup(); err != nil {
		a.logger.Error("백업 실패: %v", err)
		return fmt.Errorf("백업 실패: %w", err)
	}
	return nil
}

// Restore는 백업된 파일들을 복원합니다
func (a *fileAdapter) Restore() error {
	if err := a.service.Restore(); err != nil {
		a.logger.Error("복원 실패: %v", err)
		return fmt.Errorf("복원 실패: %w", err)
	}
	return nil
}

// ExtractZip은 ZIP 파일을 압축 해제합니다
func (a *fileAdapter) ExtractZip(zipFile string) error {
	if err := a.service.ExtractZip(zipFile); err != nil {
		a.logger.Error("압축 해제 실패: %v", err)
		return fmt.Errorf("압축 해제 실패: %w", err)
	}
	return nil
}

// DeleteFile은 지정된 파일을 삭제합니다
func (a *fileAdapter) DeleteFile(path string) error {
	if err := a.service.DeleteFile(path); err != nil {
		a.logger.Error("파일 삭제 실패: %v", err)
		return fmt.Errorf("파일 삭제 실패: %w", err)
	}
	return nil
}
