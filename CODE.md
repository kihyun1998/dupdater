# dupdater
## Project Structure

```
dupdater/
├── cmd/
    └── dupdater/
    │   └── main.go
├── internal/
    ├── file/
    │   ├── domain/
    │   │   ├── entity/
    │   │   │   └── dir_info.go
    │   │   ├── repository/
    │   │   │   └── file_repository.go
    │   │   └── usecase/
    │   │   │   └── file_service.go
    │   ├── infrastructure/
    │   │   └── local_file_repository.go
    │   └── factory.go
    ├── hash/
    │   ├── domain/
    │   │   ├── entity/
    │   │   │   └── hash_info.go
    │   │   ├── repository/
    │   │   │   └── hash_repo.go
    │   │   └── usecase/
    │   │   │   └── hash_service.go
    │   ├── infrastructure/
    │   │   └── hash_manager.go
    │   └── factory.go
    ├── i18n/
    │   ├── domain/
    │   │   ├── entity/
    │   │   │   ├── locale.go
    │   │   │   └── message.go
    │   │   ├── repository/
    │   │   │   └── i18n_repository.go
    │   │   └── usecase/
    │   │   │   └── i18n_service.go
    │   ├── infrastructure/
    │   │   ├── providers/
    │   │   │   ├── en_provider.go
    │   │   │   └── ko_provider.go
    │   │   └── i18n_store.go
    │   └── factory.go
    ├── logger/
    │   ├── domain/
    │   │   ├── entity/
    │   │   │   ├── log_entry.go
    │   │   │   └── log_level.go
    │   │   ├── repository/
    │   │   │   └── log_repository.go
    │   │   └── usecase/
    │   │   │   └── logger_service.go
    │   ├── infrastructure/
    │   │   └── file_logger.go
    │   └── factory.go
    ├── network/
    │   ├── domain/
    │   │   ├── entity/
    │   │   │   ├── server_config.go
    │   │   │   └── update_file.go
    │   │   ├── repository/
    │   │   │   └── network_repository.go
    │   │   └── usecase/
    │   │   │   └── network_service.go
    │   ├── infrastructure/
    │   │   └── repository/
    │   │   │   └── http_repository.go
    │   └── factory.go
    ├── scenario/
    │   ├── entity/
    │   │   ├── scenario.go
    │   │   └── step.go
    │   ├── scenarios/
    │   │   ├── base.go
    │   │   ├── error.go
    │   │   └── success.go
    │   ├── usecase/
    │   │   └── scenario_service.go
    │   └── factory.go
    ├── ui/
    │   ├── components/
    │   │   └── status_card.go
    │   ├── domain/
    │   │   ├── entity/
    │   │   │   └── ui_state.go
    │   │   ├── repository/
    │   │   │   └── ui_repository.go
    │   │   └── usecase/
    │   │   │   └── ui_service.go
    │   ├── infrastructure/
    │   │   └── fyne_manager.go
    │   └── factory.go
    ├── updater/
    │   ├── domain/
    │   │   ├── entity/
    │   │   │   ├── update_config.go
    │   │   │   └── update_status.go
    │   │   ├── repository/
    │   │   │   └── update_repository.go
    │   │   └── usecase/
    │   │   │   └── update_service.go
    │   ├── infrastructure/
    │   │   └── update_manager.go
    │   └── factory.go
    └── version/
    │   ├── domain/
    │       ├── entity/
    │       │   └── version.go
    │       ├── repository/
    │       │   └── version_repo.go
    │       └── usecase/
    │       │   └── version_service.go
    │   ├── infrastructure/
    │       └── version_store.go
    │   └── factory.go
└── pkg/
    └── utils/
        ├── fonts/
            └── fonts.go
        ├── theme/
            ├── dark_variant.go
            ├── light_variant.go
            ├── theme_variant.go
            └── types.go
        ├── count.go
        └── system.go
```

## cmd/dupdater/main.go
```go
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kihyun1998/dupdater/internal/file"
	"github.com/kihyun1998/dupdater/internal/hash"
	"github.com/kihyun1998/dupdater/internal/i18n"
	"github.com/kihyun1998/dupdater/internal/logger"
	logEntity "github.com/kihyun1998/dupdater/internal/logger/domain/entity"
	"github.com/kihyun1998/dupdater/internal/network"
	"github.com/kihyun1998/dupdater/internal/scenario"
	"github.com/kihyun1998/dupdater/internal/ui"
	"github.com/kihyun1998/dupdater/internal/updater"
	"github.com/kihyun1998/dupdater/internal/version"
	"github.com/kihyun1998/dupdater/pkg/utils/theme"
)

const (
	AppName       = "dupdater"
	TargetAppName = "simple_update_test.exe"
	HashSumTxt    = "hash_sum.txt"
	TotalSteps    = 8 // 총 업데이트 단계 수
)

var (
	fromVersion = flag.String("fromVersion", "", "현재 앱 버전")
	toVersion   = flag.String("toVersion", "", "업데이트할 버전")
	serverName  = flag.String("server", "", "서버 프로필 이름")
	testMode    = flag.Bool("test", false, "테스트 모드 활성화")
	testType    = flag.String("testType", "", "테스트 시나리오 유형 (success, error1, error2, ..., error8)")
	themeMode   = flag.String("theme", "light", "테마 모드 (light/dark)")
	langMode    = flag.String("lang", "ko", "언어 설정 (ko/en)")
)

func main() {
	// 1. 커맨드라인 플래그 파싱
	flag.Parse()

	// 2. 로거 초기화
	logPath := getLogPath()
	logger, err := logger.New(logger.Config{
		LogPath:    logPath,
		LogLevel:   logEntity.INFO,
		MaxSize:    10 * 1024 * 1024, // 10MB
		MaxBackups: 5,
	})
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Info("Starting updater - version: %s, server: %s", *fromVersion, *serverName)

	// 3. 테마 초기화
	if err := theme.InitTheme(*themeMode); err != nil {
		fmt.Printf("테마 초기화 실패: %v\n", err)
		os.Exit(1)
	}

	// 4. i18n 매니저 초기화
	i18nManager, err := i18n.New(i18n.Config{
		DefaultLanguage: *langMode,
		Logger:          logger, // logger도 주입해야 합니다
	})
	if err != nil {
		fmt.Printf("다국어 지원 초기화 실패: %v\n", err)
		os.Exit(1)
	}

	// 테스트 모드 체크
	if *testMode {
		if *testType != "" {
			runScenarioTest(i18nManager)
			return
		}
		runTestMode(i18nManager)
		return
	}

	// 필수 인자 체크
	if *fromVersion == "" || *toVersion == "" {
		fmt.Println("Error: fromVersion and toVersion are required")
		flag.Usage()
		os.Exit(1)
	}

	// 5. 버전 매니저 초기화
	versionManager, err := version.New(version.Config{
		Logger:      logger,
		FromVersion: *fromVersion,
		ToVersion:   *toVersion,
	})
	if err != nil {
		logger.Error("버전 관리자 초기화 실패: %v", err)
		os.Exit(1)
	}

	// 버전 매니저 초기화 후 로깅
	logger.Info("업데이트 진행: %s -> %s", versionManager.GetFromVersion(), versionManager.GetToVersion())

	// 6. UI 매니저 초기화

	uiManager := ui.New(ui.Config{
		AppName:     AppName,
		TotalSteps:  TotalSteps,
		Logger:      logger,
		FromVersion: versionManager.GetFromVersion(),
		ToVersion:   versionManager.GetToVersion(),
		Theme:       theme.GetCurrentVariant(),
		I18n:        i18nManager,
	})

	// 7. 네트워크 매니저 초기화
	networkManager := network.New(network.Config{
		Logger: logger,
	})

	// 8. 파일 매니저 초기화
	fileManager, err := file.New(file.Config{
		Logger:     logger,
		BackupDir:  filepath.Join(os.TempDir(), "ACRABACK"),
		CurrentDir: getCurrentDir(),
	})
	if err != nil {
		logger.Error("Failed to initialize file manager: %v", err)
		uiManager.ShowError(fmt.Errorf("failed to initialize file manager: %w", err))
		os.Exit(1)
	}

	// 9. 해시 매니저 초기화
	hashManager, err := hash.New(hash.Config{
		Logger:      logger,
		CurrentDir:  getCurrentDir(),
		HashSumPath: filepath.Join(getCurrentDir(), HashSumTxt),
	})
	if err != nil {
		logger.Error("Failed to initialize hash manager: %v", err)
		uiManager.ShowError(fmt.Errorf("failed to initialize hash manager: %w", err))
		os.Exit(1)
	}

	// 10. Updater 생성 및 시작
	updater, err := updater.New(updater.Config{
		AppName:         TargetAppName,
		FromVersion:     *fromVersion,
		ServerName:      *serverName,
		BackupCompleted: false,
		UIManager:       uiManager,
		Logger:          logger,
		NetworkManager:  networkManager,
		FileManager:     fileManager,
		HashManager:     hashManager,
		I18n:            i18nManager,
	})
	if err != nil {
		logger.Error("Failed to initialize updater: %v", err)
		os.Exit(1)
	}

	// 11. 업데이트 프로세스 시작
	updater.Start()
}

// getLogPath는 로그 파일의 경로를 반환합니다.
func getLogPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Failed to get home directory: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(homeDir, ".testfolder", "logs", "updater.log")
}

// getCurrentDir는 현재 실행 파일의 디렉토리 경로를 반환합니다.
func getCurrentDir() string {
	dir, err := os.Executable()
	if err != nil {
		fmt.Printf("Failed to get executable path: %v\n", err)
		os.Exit(1)
	}
	return filepath.Dir(dir)
}

// runTestMode는 UI 테스트를 위한 모드를 실행합니다
func runTestMode(i18nManager i18n.Manager) {
	// 로거 초기화
	logPath := getLogPath()
	logger, err := logger.New(logger.Config{
		LogPath:    logPath,
		LogLevel:   logEntity.INFO,
		MaxSize:    10 * 1024 * 1024,
		MaxBackups: 5,
	})
	if err != nil {
		fmt.Printf("로거 초기화 실패: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	// 테스트 모드용 버전 매니저 초기화
	testVersionManager, err := version.New(version.Config{
		Logger:      logger,
		FromVersion: "V3.0.0(2024-01-01)",
		ToVersion:   "V3.0.1(2024-02-01)",
	})
	if err != nil {
		logger.Error("버전 관리자 초기화 실패: %v", err)
		os.Exit(1)
	}

	// UI 매니저 초기화
	uiManager := ui.New(ui.Config{
		AppName:     AppName,
		TotalSteps:  TotalSteps,
		Logger:      logger,
		FromVersion: testVersionManager.GetFromVersion(),
		ToVersion:   testVersionManager.GetToVersion(),
		Theme:       theme.GetCurrentVariant(),
		I18n:        i18nManager,
	})

	// UI 실행
	uiManager.Run()
}

func runScenarioTest(i18nManager i18n.Manager) {
	// 로거 초기화
	logPath := getLogPath()
	logger, err := logger.New(logger.Config{
		LogPath:    logPath,
		LogLevel:   logEntity.INFO,
		MaxSize:    10 * 1024 * 1024,
		MaxBackups: 5,
	})
	if err != nil {
		fmt.Printf("로거 초기화 실패: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	// 테스트용 버전 매니저 초기화
	testVersionManager, err := version.New(version.Config{
		Logger:      logger,
		FromVersion: "V3.0.0(2024-01-01)",
		ToVersion:   "V3.0.1(2024-02-01)",
	})
	if err != nil {
		fmt.Printf("버전 매니저 초기화 실패: %v\n", err)
		os.Exit(1)
	}

	// UI 매니저 생성 - 인터페이스로 받음
	uiManager := ui.New(ui.Config{
		AppName:     AppName,
		TotalSteps:  TotalSteps,
		Logger:      logger,
		FromVersion: testVersionManager.GetFromVersion(),
		ToVersion:   testVersionManager.GetToVersion(),
		Theme:       theme.GetCurrentVariant(),
		I18n:        i18nManager,
	})

	// 시나리오 매니저 생성
	scenarioManager, err := scenario.New(scenario.Config{
		Logger:       logger,
		UIManager:    uiManager,
		I18n:         i18nManager,
		ScenarioType: *testType,
	})
	if err != nil {
		fmt.Printf("시나리오 매니저 생성 실패: %v\n", err)
		os.Exit(1)
	}

	// 시나리오 실행
	scenarioManager.Run()
}

```
## internal/file/domain/entity/dir_info.go
```go
package entity

// DirInfo는 디렉토리 정보를 담는 도메인 엔티티입니다
type DirInfo struct {
	BackupDir  string // 백업 디렉토리 경로
	CurrentDir string // 현재 작업 디렉토리 경로
}

// NewDirInfo는 새로운 DirInfo 인스턴스를 생성합니다
func NewDirInfo(backupDir, currentDir string) *DirInfo {
	return &DirInfo{
		BackupDir:  backupDir,
		CurrentDir: currentDir,
	}
}

// GetBackupDir은 백업 디렉토리 경로를 반환합니다
func (d *DirInfo) GetBackupDir() string {
	return d.BackupDir
}

// GetCurrentDir은 현재 작업 디렉토리 경로를 반환합니다
func (d *DirInfo) GetCurrentDir() string {
	return d.CurrentDir
}

```
## internal/file/domain/repository/file_repository.go
```go
// Package repository는 파일 시스템 작업을 위한 인터페이스를 정의합니다
package repository

import (
	"github.com/kihyun1998/dupdater/internal/file/domain/entity"
	"github.com/kihyun1998/dupdater/internal/logger"
)

// FileRepository는 파일 시스템 작업을 추상화하는 인터페이스입니다
type FileRepository interface {
	// Backup은 파일을 백업합니다
	Backup() error

	// Restore는 백업된 파일을 복원합니다
	Restore() error

	// ExtractZip은 ZIP 파일을 압축 해제합니다
	ExtractZip(zipFile string) error

	// DeleteFile은 파일을 삭제합니다
	DeleteFile(path string) error

	// GetBackupDir은 백업 디렉토리 경로를 반환합니다
	GetBackupDir() string
}

// Config는 FileRepository 생성에 필요한 설정을 정의합니다
type Config struct {
	DirInfo *entity.DirInfo
	Logger  logger.Logger // 로거 인터페이스
}

```
## internal/file/domain/usecase/file_service.go
```go
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

```
## internal/file/factory.go
```go
// Package file은 파일 시스템 작업의 진입점을 제공합니다
package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kihyun1998/dupdater/internal/file/domain/entity"
	"github.com/kihyun1998/dupdater/internal/file/domain/repository"
	"github.com/kihyun1998/dupdater/internal/file/domain/usecase"
	"github.com/kihyun1998/dupdater/internal/file/infrastructure"
	"github.com/kihyun1998/dupdater/internal/logger"
)

// Manager는 파일 관리를 위한 인터페이스입니다
type Manager interface {
	Backup() error
	Restore() error
	ExtractZip(zipFile string) error
	DeleteFile(path string) error
	GetBackupDir() string
}

// Config는 Manager 생성에 필요한 설정입니다
type Config struct {
	Logger     logger.Logger
	BackupDir  string
	CurrentDir string
}

// manager는 Manager 인터페이스의 구현체입니다
type manager struct {
	service *usecase.FileService
}

// New는 새로운 Manager 인스턴스를 생성합니다
func New(config Config) (Manager, error) {
	// 기본값 설정
	if config.BackupDir == "" {
		config.BackupDir = filepath.Join(os.TempDir(), "ACRABACK")
	}

	if config.CurrentDir == "" {
		dir, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("현재 디렉토리 확인 실패: %w", err)
		}
		config.CurrentDir = dir
	}

	dirInfo := entity.NewDirInfo(config.BackupDir, config.CurrentDir)

	// Repository 설정
	repoConfig := &repository.Config{
		DirInfo: dirInfo,
		Logger:  config.Logger,
	}

	// Repository 생성
	repo := infrastructure.NewLocalFileRepository(repoConfig)

	// Service 생성
	service := usecase.NewFileService(repo, config.Logger)

	return &manager{
		service: service,
	}, nil
}

// Manager 인터페이스 구현
func (m *manager) Backup() error {
	return m.service.Backup()
}

func (m *manager) Restore() error {
	return m.service.Restore()
}

func (m *manager) ExtractZip(zipFile string) error {
	return m.service.ExtractZip(zipFile)
}

func (m *manager) DeleteFile(path string) error {
	return m.service.DeleteFile(path)
}

func (m *manager) GetBackupDir() string {
	return m.service.GetBackupDir()
}

```
## internal/file/infrastructure/local_file_repository.go
```go
// Package infrastructure는 파일 시스템 작업의 실제 구현체를 제공합니다
package infrastructure

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/kihyun1998/dupdater/internal/file/domain/repository"
)

// LocalFileRepository는 로컬 파일 시스템 기반의 저장소 구현체입니다
type LocalFileRepository struct {
	config *repository.Config
	mu     sync.RWMutex
}

// NewLocalFileRepository는 새로운 LocalFileRepository 인스턴스를 생성합니다
func NewLocalFileRepository(config *repository.Config) repository.FileRepository {
	return &LocalFileRepository{
		config: config,
	}
}

// cleanCurrentDirectory는 현재 디렉토리를 정리합니다
func (r *LocalFileRepository) cleanCurrentDirectory() error {
	files, err := os.ReadDir(r.config.DirInfo.CurrentDir)
	if err != nil {
		return fmt.Errorf("디렉토리 읽기 실패: %w", err)
	}

	for _, file := range files {
		path := filepath.Join(r.config.DirInfo.CurrentDir, file.Name())

		// 여러 번 삭제 시도
		for i := 0; i < 3; i++ {
			var err error
			if file.IsDir() {
				r.config.Logger.Info("디렉토리 삭제 시도: %s", path)
				err = os.RemoveAll(path)
			} else {
				r.config.Logger.Info("파일 삭제 시도: %s", path)
				err = os.Remove(path)
			}

			if err == nil {
				break
			}

			// 마지막 시도에서는 .old를 붙여서 시도
			if i == 2 {
				newPath := path + ".old"
				r.config.Logger.Error("삭제 실패, 이름 변경 후 재시도: %s -> %s", path, newPath)
				if os.Rename(path, newPath) == nil {
					os.Remove(newPath)
				}
			}

			time.Sleep(time.Second)
		}
	}
	return nil
}

// DeleteFile은 지정된 파일을 삭제합니다
func (r *LocalFileRepository) DeleteFile(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("파일 삭제 실패: %w", err)
	}
	return nil
}

// Extract ----------------------------------------------------------
// ExtractZip은 ZIP 파일을 압축 해제합니다
func (r *LocalFileRepository) ExtractZip(zipFile string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("ZIP 파일 열기 실패: %w", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		if err := r.extractFile(file); err != nil {
			return fmt.Errorf("파일 압축해제 실패 (%s): %w", file.Name, err)
		}
	}

	return nil
}

// extractFile은 ZIP 파일의 항목을 압축 해제합니다
func (r *LocalFileRepository) extractFile(file *zip.File) error {
	filePath := filepath.Join(r.config.DirInfo.CurrentDir, file.Name)

	if file.FileInfo().IsDir() {
		return os.MkdirAll(filePath, os.ModePerm)
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return fmt.Errorf("디렉토리 생성 실패: %w", err)
	}

	dest, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return fmt.Errorf("대상 파일 생성 실패: %w", err)
	}
	defer dest.Close()

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("ZIP 파일 열기 실패: %w", err)
	}
	defer src.Close()

	if _, err := io.Copy(dest, src); err != nil {
		return fmt.Errorf("파일 복사 실패: %w", err)
	}

	return nil
}

// ------------------------------------------------------------------

// Restore ----------------------------------------------------------
// Restore는 백업된 파일들을 복원합니다
func (r *LocalFileRepository) Restore() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 백업 디렉토리 존재 확인
	if _, err := os.Stat(r.config.DirInfo.BackupDir); os.IsNotExist(err) {
		return fmt.Errorf("백업 디렉토리가 존재하지 않습니다: %s", r.config.DirInfo.BackupDir)
	}

	// 현재 디렉토리 정리
	if err := r.cleanCurrentDirectory(); err != nil {
		return fmt.Errorf("현재 디렉토리 정리 실패: %w", err)
	}

	// 백업 파일 복원
	files, err := os.ReadDir(r.config.DirInfo.BackupDir)
	if err != nil {
		return fmt.Errorf("백업 디렉토리 읽기 실패: %w", err)
	}

	for _, file := range files {
		srcPath := filepath.Join(r.config.DirInfo.BackupDir, file.Name())
		destPath := filepath.Join(r.config.DirInfo.CurrentDir, file.Name())

		if file.IsDir() {
			if err := r.restoreDirectory(srcPath, destPath); err != nil {
				return fmt.Errorf("디렉토리 복원 실패 (%s): %w", file.Name(), err)
			}
		} else {
			if err := r.restoreFile(srcPath, destPath); err != nil {
				return fmt.Errorf("파일 복원 실패 (%s): %w", file.Name(), err)
			}
		}
	}

	return nil
}

// restoreDirectory는 디렉토리를 복구합니다.
func (r *LocalFileRepository) restoreDirectory(src, dest string) error {
	// 디렉토리 존재하면 생성
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return fmt.Errorf("디렉토리 생성 실패: %w", err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("디렉토리 읽기 실패: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			if err := r.restoreDirectory(srcPath, destPath); err != nil {
				return err
			}
		} else {
			if err := r.restoreFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}
	r.config.Logger.Info("디렉토리 복원 완료: %s", dest)
	return os.RemoveAll(src)
}

// restoreFile은 파일을 복원합니다
func (r *LocalFileRepository) restoreFile(src, dest string) error {
	// 기존 파일이 존재하면 삭제
	if _, err := os.Stat(dest); err == nil {
		if err := os.Remove(dest); err != nil {
			r.config.Logger.Info("기존 파일 삭제 시도: %s", dest)
			if err := os.Remove(dest); err != nil {
				return fmt.Errorf("기존 파일 삭제 실패: %w", err)
			}
		}
	}

	// Rename 시도
	if err := os.Rename(src, dest); err == nil {
		r.config.Logger.Info("파일 복원 완료(Rename): %s -> %s", src, dest)
		return nil
	}

	// Rename 실패시 복사
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("소스 파일 열기 실패: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("대상 파일 생성 실패: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("파일 복사 실패: %w", err)
	}

	r.config.Logger.Info("파일 복원 완료(Copy): %s -> %s", src, dest)
	return os.Remove(src)
}

// ------------------------------------------------------------------

// Backup------------------------------------------------------------
// Backup은 현재 디렉토리의 파일들을 백업합니다
func (r *LocalFileRepository) Backup() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 백업 디렉토리 생성
	if err := os.MkdirAll(r.config.DirInfo.BackupDir, os.ModePerm); err != nil {
		return fmt.Errorf("백업 디렉토리 생성 실패: %w", err)
	}

	// 현재 디렉토리 파일 목록 조회
	files, err := os.ReadDir(r.config.DirInfo.CurrentDir)
	if err != nil {
		return fmt.Errorf("디렉토리 읽기 실패: %w", err)
	}

	// 각 파일 백업
	for _, file := range files {
		oldPath := filepath.Join(r.config.DirInfo.CurrentDir, file.Name())
		newPath := filepath.Join(r.config.DirInfo.BackupDir, file.Name())

		if file.IsDir() {
			if err := r.backupDirectory(oldPath, newPath); err != nil {
				return fmt.Errorf("디렉토리 백업 실패 (%s): %w", file.Name(), err)
			}
		} else {
			if err := r.backupFile(oldPath, newPath); err != nil {
				return fmt.Errorf("파일 백업 실패 (%s): %w", file.Name(), err)
			}
		}
	}

	return nil
}

// GetBackupDir은 백업 디렉토리 경로를 반환합니다
func (r *LocalFileRepository) GetBackupDir() string {
	return r.config.DirInfo.BackupDir
}

// backupDirectory는 디렉토리를 백업합니다
func (r *LocalFileRepository) backupDirectory(src, dest string) error {
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return fmt.Errorf("디렉토리 생성 실패: %w", err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("디렉토리 읽기 실패: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			if err := r.backupDirectory(srcPath, destPath); err != nil {
				return err
			}
		} else {
			if err := r.backupFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return os.RemoveAll(src)
}

// backupFile은 단일 파일을 백업합니다
func (r *LocalFileRepository) backupFile(src, dest string) error {
	// 먼저 Rename 시도
	if err := os.Rename(src, dest); err == nil {
		return nil
	}

	// Rename 실패시 복사 후 삭제
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("소스 파일 열기 실패: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("대상 파일 생성 실패: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("파일 복사 실패: %w", err)
	}

	// 원본 파일 삭제 시도
	for i := 0; i < 3; i++ {
		if err := os.Remove(src); err == nil {
			return nil
		}
		time.Sleep(time.Second)
	}

	return os.Remove(src)
}

// ------------------------------------------------------------------

```
## internal/hash/domain/entity/hash_info.go
```go
// Package entity는 해시 도메인의 핵심 개념을 정의합니다
package entity

import (
	"fmt"
)

// FileHash는 파일의 해시 정보를 담는 도메인 엔티티입니다
type FileHash struct {
	Path     string // 파일 경로
	HashSum  string // 해시 값
	FileType string // 파일 유형 (파일/디렉토리)
}

// NewFileHash는 새로운 FileHash 인스턴스를 생성합니다
func NewFileHash(path string, hashSum string, fileType string) *FileHash {
	return &FileHash{
		Path:     path,
		HashSum:  hashSum,
		FileType: fileType,
	}
}

// Validate는 해시 정보가 유효한지 검증합니다
func (f *FileHash) Validate() error {
	if f.Path == "" {
		return fmt.Errorf("파일 경로가 비어있습니다")
	}
	if f.HashSum == "" {
		return fmt.Errorf("해시값이 비어있습니다")
	}
	if f.FileType == "" {
		return fmt.Errorf("파일 유형이 비어있습니다")
	}
	return nil
}

```
## internal/hash/domain/repository/hash_repo.go
```go
package repository

import (
	"github.com/kihyun1998/dupdater/internal/hash/domain/entity"
	"github.com/kihyun1998/dupdater/internal/logger"
)

// HashRepository는 해시 검증을 위한 저장소 인터페이스입니다
type HashRepository interface {
	// VerifyFile은 단일 파일의 해시를 검증합니다
	VerifyFile(filePath string, expectedHash string) error

	// VerifyUpdateFile은 업데이트 파일의 해시를 검증합니다
	VerifyUpdateFile(filePath string) error

	// VerifyHashSum은 모든 파일의 해시섬을 검증합니다
	VerifyHashSum() error

	// GetFileHash는 파일의 해시 정보를 조회합니다
	GetFileHash(filePath string) (*entity.FileHash, error)
}

// Config는 저장소 설정을 정의합니다
type Config struct {
	CurrentDir  string
	HashSumPath string
	Logger      logger.Logger
}

```
## internal/hash/domain/usecase/hash_service.go
```go
// Package usecase는 해시 도메인의 비즈니스 로직을 구현합니다
package usecase

import (
	"fmt"

	"github.com/kihyun1998/dupdater/internal/hash/domain/repository"
	"github.com/kihyun1998/dupdater/internal/logger"
)

// HashService는 해시 검증 관련 비즈니스 로직을 구현합니다
type HashService struct {
	repo   repository.HashRepository
	logger logger.Logger
}

// NewHashService는 새로운 HashService 인스턴스를 생성합니다
func NewHashService(repo repository.HashRepository, logger logger.Logger) *HashService {
	return &HashService{
		repo:   repo,
		logger: logger,
	}
}

// VerifyFile은 단일 파일의 해시를 검증합니다
func (s *HashService) VerifyFile(filePath string, expectedHash string) error {
	s.logger.Info("파일 해시 검증 시작: %s", filePath)

	if err := s.repo.VerifyFile(filePath, expectedHash); err != nil {
		s.logger.Error("파일 해시 검증 실패: %v", err)
		return fmt.Errorf("파일 해시 검증 실패: %w", err)
	}

	s.logger.Info("파일 해시 검증 완료: %s", filePath)
	return nil
}

// VerifyUpdateFile은 업데이트 파일의 해시를 검증합니다
func (s *HashService) VerifyUpdateFile(filePath string) error {
	s.logger.Info("업데이트 파일 해시 검증 시작: %s", filePath)

	if err := s.repo.VerifyUpdateFile(filePath); err != nil {
		s.logger.Error("업데이트 파일 해시 검증 실패: %v", err)
		return fmt.Errorf("업데이트 파일 해시 검증 실패: %w", err)
	}

	s.logger.Info("업데이트 파일 해시 검증 완료: %s", filePath)
	return nil
}

// VerifyHashSum은 모든 파일의 해시섬을 검증합니다
func (s *HashService) VerifyHashSum() error {
	s.logger.Info("전체 파일 해시섬 검증 시작")

	if err := s.repo.VerifyHashSum(); err != nil {
		s.logger.Error("전체 파일 해시섬 검증 실패: %v", err)
		return fmt.Errorf("전체 파일 해시섬 검증 실패: %w", err)
	}

	s.logger.Info("전체 파일 해시섬 검증 완료")
	return nil
}

```
## internal/hash/factory.go
```go
// Package hash는 해시 도메인의 진입점을 제공합니다
package hash

import (
	"github.com/kihyun1998/dupdater/internal/hash/domain/repository"
	"github.com/kihyun1998/dupdater/internal/hash/domain/usecase"
	"github.com/kihyun1998/dupdater/internal/hash/infrastructure"
	"github.com/kihyun1998/dupdater/internal/logger"
)

// Config는 해시 매니저 생성에 필요한 설정입니다
type Config struct {
	CurrentDir  string
	HashSumPath string
	Logger      logger.Logger
}

// Manager는 해시 검증을 위한 인터페이스입니다
type Manager interface {
	VerifyFile(filePath string, expectedHash string) error
	VerifyUpdateFile(filePath string) error
	VerifyHashSum() error
}

// 실제 구현체
type manager struct {
	service *usecase.HashService
}

// New는 새로운 해시 매니저를 생성합니다
func New(config Config) (Manager, error) {
	// 리포지토리 설정
	repoConfig := &repository.Config{
		CurrentDir:  config.CurrentDir,
		HashSumPath: config.HashSumPath,
		Logger:      config.Logger,
	}

	// 해시 매니저 생성
	repo := infrastructure.NewFileHashManager(repoConfig)

	// 서비스 생성
	service := usecase.NewHashService(repo, config.Logger)

	return &manager{
		service: service,
	}, nil
}

// 인터페이스 구현
func (m *manager) VerifyFile(filePath string, expectedHash string) error {
	return m.service.VerifyFile(filePath, expectedHash)
}

func (m *manager) VerifyUpdateFile(filePath string) error {
	return m.service.VerifyUpdateFile(filePath)
}

func (m *manager) VerifyHashSum() error {
	return m.service.VerifyHashSum()
}

```
## internal/hash/infrastructure/hash_manager.go
```go
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

	sumFilePath := m.config.HashSumPath
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

```
## internal/i18n/domain/entity/locale.go
```go
package entity

// Language는 지원되는 언어 코드를 정의합니다
type Language string

const (
	// Korean 한국어
	Korean Language = "ko"

	// English 영어
	English Language = "en"

	// DefaultLanguage 기본 언어를 한국어로 설정
	DefaultLanguage = Korean
)

// Locale은 로케일 정보를 담는 도메인 엔티티입니다
type Locale struct {
	Language Language
	Messages map[string]string
}

// NewLocale은 새로운 Locale 인스턴스를 생성합니다
func NewLocale(lang Language, messages map[string]string) *Locale {
	if messages == nil {
		messages = make(map[string]string)
	}
	return &Locale{
		Language: lang,
		Messages: messages,
	}
}

// GetMessage는 지정된 키에 해당하는 메시지를 반환합니다
func (l *Locale) GetMessage(key string) string {
	if msg, ok := l.Messages[key]; ok {
		return msg
	}
	return key
}

// IsValid는 로케일이 유효한지 검증합니다
func (l *Locale) IsValid() bool {
	return l.Language == Korean || l.Language == English
}

```
## internal/i18n/domain/entity/message.go
```go
package entity

// MessageProvider는 특정 언어의 메시지를 제공하는 인터페이스입니다
type MessageProvider interface {
	// GetLanguageCode는 제공자가 지원하는 언어 코드를 반환합니다
	GetLanguage() Language

	// GetMessages는 해당 언어의 전체 메시지 맵을 반환합니다
	GetMessages() map[string]string
}

// Message는 다국어 메시지를 나타내는 값 객체입니다
type Message struct {
	Key   string
	Value string
	Lang  Language
}

// NewMessage는 새로운 Message 인스턴스를 생성합니다
func NewMessage(key, value string, lang Language) *Message {
	return &Message{
		Key:   key,
		Value: value,
		Lang:  lang,
	}
}

// String은 메시지의 문자열 표현을 반환합니다
func (m *Message) String() string {
	return m.Value
}

```
## internal/i18n/domain/repository/i18n_repository.go
```go
package repository

import (
	"github.com/kihyun1998/dupdater/internal/i18n/domain/entity"
	"github.com/kihyun1998/dupdater/internal/logger"
)

// I18nRepository는 다국어 지원을 위한 저장소 인터페이스입니다
type I18nRepository interface {
	// GetMessage는 현재 언어의 메시지를 조회합니다
	GetMessage(key string) string

	// SetLanguage는 현재 언어를 설정합니다
	SetLanguage(lang entity.Language) error

	// GetCurrentLanguage는 현재 설정된 언어를 반환합니다
	GetCurrentLanguage() entity.Language

	// RegisterProvider는 새로운 메시지 제공자를 등록합니다
	RegisterProvider(provider entity.MessageProvider) error
}

// Config는 저장소 설정을 정의합니다
type Config struct {
	DefaultLanguage entity.Language
	Logger          logger.Logger
}

```
## internal/i18n/domain/usecase/i18n_service.go
```go
package usecase

import (
	"fmt"
	"sync"

	"github.com/kihyun1998/dupdater/internal/i18n/domain/entity"
	"github.com/kihyun1998/dupdater/internal/i18n/domain/repository"
	"github.com/kihyun1998/dupdater/internal/logger"
)

// I18nService는 다국어 지원의 비즈니스 로직을 구현합니다
type I18nService struct {
	repo   repository.I18nRepository
	logger logger.Logger
	mu     sync.RWMutex
}

// NewI18nService는 새로운 I18nService 인스턴스를 생성합니다
func NewI18nService(repo repository.I18nRepository, logger logger.Logger) *I18nService {
	return &I18nService{
		repo:   repo,
		logger: logger,
	}
}

// GetMessage는 지정된 키에 해당하는 메시지를 현재 언어로 반환합니다
func (s *I18nService) GetMessage(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	msg := s.repo.GetMessage(key)
	if msg == key {
		s.logger.Error("메시지를 찾을 수 없음: %s", key)
	}
	return msg
}

// SetLanguage는 현재 언어를 설정합니다
func (s *I18nService) SetLanguage(lang string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.repo.SetLanguage(entity.Language(lang)); err != nil {
		s.logger.Error("언어 설정 실패: %v", err)
		return fmt.Errorf("언어 설정 실패: %w", err)
	}

	s.logger.Info("언어가 변경됨: %s", lang)
	return nil
}

// GetCurrentLanguage는 현재 설정된 언어를 반환합니다
func (s *I18nService) GetCurrentLanguage() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return string(s.repo.GetCurrentLanguage())
}

// RegisterProvider는 새로운 메시지 제공자를 등록합니다
func (s *I18nService) RegisterProvider(provider entity.MessageProvider) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.repo.RegisterProvider(provider); err != nil {
		s.logger.Error("메시지 제공자 등록 실패: %v", err)
		return fmt.Errorf("메시지 제공자 등록 실패: %w", err)
	}

	s.logger.Info("메시지 제공자가 등록됨: %s", provider.GetLanguage())
	return nil
}

```
## internal/i18n/factory.go
```go
// factory.go
package i18n

import (
	"fmt"

	"github.com/kihyun1998/dupdater/internal/i18n/domain/entity"
	"github.com/kihyun1998/dupdater/internal/i18n/domain/repository"
	"github.com/kihyun1998/dupdater/internal/i18n/domain/usecase"
	"github.com/kihyun1998/dupdater/internal/i18n/infrastructure"
	"github.com/kihyun1998/dupdater/internal/i18n/infrastructure/providers"
	"github.com/kihyun1998/dupdater/internal/logger"
)

// Manager는 다국어 지원을 위한 인터페이스입니다
type Manager interface {
	GetMessage(key string) string
	SetLanguage(lang string) error
	GetCurrentLanguage() string
}

// Config는 Manager 생성에 필요한 설정입니다
type Config struct {
	DefaultLanguage string // 기본 언어 설정
	Logger          logger.Logger
}

// managerImpl은 Manager 인터페이스 구현체입니다
type managerImpl struct {
	service *usecase.I18nService
}

// New는 새로운 Manager 인스턴스를 생성합니다
func New(config Config) (Manager, error) {
	// 설정 검증
	if config.DefaultLanguage == "" {
		config.DefaultLanguage = string(entity.DefaultLanguage)
	}

	// 저장소 설정
	repoConfig := &repository.Config{
		DefaultLanguage: entity.Language(config.DefaultLanguage),
		Logger:          config.Logger,
	}

	// 저장소 생성
	store := infrastructure.NewI18nStore(repoConfig)

	// 서비스 생성
	service := usecase.NewI18nService(store, config.Logger)

	// 기본 메시지 제공자 등록
	defaultProviders := []entity.MessageProvider{
		providers.NewKoreanProvider(),
		providers.NewEnglishProvider(),
	}

	for _, provider := range defaultProviders {
		if err := service.RegisterProvider(provider); err != nil {
			return nil, fmt.Errorf("메시지 제공자 등록 실패: %w", err)
		}
	}

	// 초기 언어 설정
	if err := service.SetLanguage(config.DefaultLanguage); err != nil {
		return nil, fmt.Errorf("초기 언어 설정 실패: %w", err)
	}

	return &managerImpl{
		service: service,
	}, nil
}

// Manager 인터페이스 구현
func (m *managerImpl) GetMessage(key string) string {
	return m.service.GetMessage(key)
}

func (m *managerImpl) SetLanguage(lang string) error {
	return m.service.SetLanguage(lang)
}

func (m *managerImpl) GetCurrentLanguage() string {
	return m.service.GetCurrentLanguage()
}

```
## internal/i18n/infrastructure/i18n_store.go
```go
package infrastructure

import (
	"fmt"
	"sync"

	"github.com/kihyun1998/dupdater/internal/i18n/domain/entity"
	"github.com/kihyun1998/dupdater/internal/i18n/domain/repository"
	"github.com/kihyun1998/dupdater/internal/logger"
)

// I18nStore는 메모리 기반 다국어 저장소입니다
type I18nStore struct {
	currentLang entity.Language
	messages    map[entity.Language]map[string]string
	logger      logger.Logger
	mu          sync.RWMutex
}

// NewI18nStore는 새로운 I18nStore 인스턴스를 생성합니다
func NewI18nStore(config *repository.Config) *I18nStore {
	return &I18nStore{
		currentLang: config.DefaultLanguage,
		messages:    make(map[entity.Language]map[string]string),
		logger:      config.Logger,
	}
}

// GetMessage는 현재 언어의 메시지를 조회합니다
func (s *I18nStore) GetMessage(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if messages, ok := s.messages[s.currentLang]; ok {
		if msg, ok := messages[key]; ok {
			return msg
		}
	}
	return key
}

// SetLanguage는 현재 언어를 설정합니다
func (s *I18nStore) SetLanguage(lang entity.Language) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.messages[lang]; !ok {
		return fmt.Errorf("지원하지 않는 언어: %s", lang)
	}

	s.currentLang = lang
	return nil
}

// GetCurrentLanguage는 현재 설정된 언어를 반환합니다
func (s *I18nStore) GetCurrentLanguage() entity.Language {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentLang
}

// RegisterProvider는 새로운 메시지 제공자를 등록합니다
func (s *I18nStore) RegisterProvider(provider entity.MessageProvider) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	lang := provider.GetLanguage()
	if !entity.NewLocale(lang, nil).IsValid() {
		return fmt.Errorf("유효하지 않은 언어: %s", lang)
	}

	s.messages[lang] = provider.GetMessages()
	return nil
}

```
## internal/i18n/infrastructure/providers/en_provider.go
```go
package providers

import "github.com/kihyun1998/dupdater/internal/i18n/domain/entity"

// EnglishProvider는 영어 메시지를 제공하는 구현체입니다
type EnglishProvider struct{}

// NewEnglishProvider는 새로운 EnglishProvider 인스턴스를 생성합니다
func NewEnglishProvider() entity.MessageProvider {
	return &EnglishProvider{}
}

// GetLanguage는 제공하는 언어 코드를 반환합니다
func (p *EnglishProvider) GetLanguage() entity.Language {
	return entity.English
}

// GetMessages는 영어 메시지 맵을 반환합니다
func (p *EnglishProvider) GetMessages() map[string]string {
	return map[string]string{
		// Update Status Messages
		"update.status.checking":     "Checking application status...",
		"update.status.getting_info": "Getting update information...",
		"update.status.preparing":    "Preparing update files...",
		"update.status.downloading":  "Downloading update files...",
		"update.status.verifying":    "Verifying update files...",
		"update.status.installing":   "Installing update...",
		"update.status.finalizing":   "Finalizing installation...",
		"update.status.completed":    "Update completed.",

		// Titles and Headers
		"update.title":             "Update",
		"update.header.new_update": "New Update Available",
		"update.header.version":    "Version",

		// Progress Status
		"update.progress.downloading": "Downloading...",
		"update.progress.percentage":  "%d%%",

		// Error Messages
		"update.error.generic":      "An error occurred during update",
		"update.error.connection":   "Failed to connect to server",
		"update.error.download":     "Failed to download file",
		"update.error.verification": "File verification failed",

		// Restoration Related
		"update.restore.in_progress": "Restoring previous version...",
		"update.restore.completed":   "Restoration completed",

		// Others
		"update.info.restart":    "The application will restart automatically after update.",
		"update.button.minimize": "Minimize",
		"update.button.close":    "Close",
	}
}

```
## internal/i18n/infrastructure/providers/ko_provider.go
```go
package providers

import "github.com/kihyun1998/dupdater/internal/i18n/domain/entity"

// KoreanProvider는 한국어 메시지를 제공하는 구현체입니다
type KoreanProvider struct{}

// NewKoreanProvider는 새로운 KoreanProvider 인스턴스를 생성합니다
func NewKoreanProvider() entity.MessageProvider {
	return &KoreanProvider{}
}

// GetLanguage는 제공하는 언어 코드를 반환합니다
func (p *KoreanProvider) GetLanguage() entity.Language {
	return entity.Korean
}

// GetMessages는 한국어 메시지 맵을 반환합니다
func (p *KoreanProvider) GetMessages() map[string]string {
	return map[string]string{
		// 업데이트 상태 메시지
		"update.status.checking":     "앱 상태를 확인하고 있습니다...",
		"update.status.getting_info": "업데이트 정보를 확인하고 있습니다...",
		"update.status.preparing":    "업데이트 파일을 준비하고 있습니다...",
		"update.status.downloading":  "업데이트 파일을 다운로드하고 있습니다...",
		"update.status.verifying":    "업데이트 파일을 검증하고 있습니다...",
		"update.status.installing":   "업데이트를 설치하고 있습니다...",
		"update.status.finalizing":   "설치를 확인하고 있습니다...",
		"update.status.completed":    "업데이트가 완료되었습니다.",

		// 타이틀 및 헤더
		"update.title":             "업데이트",
		"update.header.new_update": "새로운 업데이트가 있습니다",
		"update.header.version":    "버전",

		// 진행 상태
		"update.progress.downloading": "다운로드 중...",
		"update.progress.percentage":  "%d%%",

		// 에러 메시지
		"update.error.generic":      "업데이트 중 오류가 발생했습니다",
		"update.error.connection":   "서버 연결에 실패했습니다",
		"update.error.download":     "파일 다운로드에 실패했습니다",
		"update.error.verification": "파일 검증에 실패했습니다",

		// 복원 관련
		"update.restore.in_progress": "이전 버전으로 복원 중...",
		"update.restore.completed":   "복원이 완료되었습니다",

		// 기타
		"update.info.restart":    "업데이트가 완료되면 자동으로 앱이 다시 시작됩니다.",
		"update.button.minimize": "최소화",
		"update.button.close":    "닫기",
	}
}

```
## internal/logger/domain/entity/log_entry.go
```go
package entity

import (
	"fmt"
	"time"
)

// LogEntry는 하나의 로그 항목을 나타냅니다.
type LogEntry struct {
	// Level은 로그의 중요도를 나타냅니다.
	Level LogLevel
	// Message는 로그 메시지 내용입니다.
	Message string
	// Timestamp는 로그가 생성된 시간입니다.
	Timestamp time.Time
	// CallerInfo는 로그를 생성한 코드의 위치 정보입니다.
	CallerInfo string
	// Fields는 로그에 추가되는 구조화된 데이터입니다.
	Fields map[string]interface{}
}

// NewLogEntry는 새로운 LogEntry를 생성합니다.
func NewLogEntry(level LogLevel, message string, callerInfo string) *LogEntry {
	return &LogEntry{
		Level:      level,
		Message:    message,
		Timestamp:  time.Now(),
		CallerInfo: callerInfo,
		Fields:     make(map[string]interface{}),
	}
}

// Format은 로그 엔트리를 문자열로 포맷팅합니다.
func (e *LogEntry) Format() string {
	base := fmt.Sprintf("[%s] %s [%s] %s",
		e.Timestamp.Format("2006-01-02 15:04:05"),
		e.Level.String(),
		e.CallerInfo,
		e.Message,
	)

	// 추가 필드가 있는 경우 포함
	if len(e.Fields) > 0 {
		fields := ""
		for k, v := range e.Fields {
			fields += fmt.Sprintf(" %s=%v", k, v)
		}
		base += fields
	}

	return base
}

// IsValid는 로그 엔트리가 유효한지 검사합니다.
func (e *LogEntry) IsValid() bool {
	return e.Level.IsValid() &&
		e.Message != "" &&
		!e.Timestamp.IsZero() &&
		e.CallerInfo != ""
}

```
## internal/logger/domain/entity/log_level.go
```go
package entity

// LogLevel은 로그의 중요도를 나타냅니다.
type LogLevel int

const (
	// DEBUG는 디버깅 목적의 상세 정보를 나타냅니다.
	DEBUG LogLevel = iota
	// INFO는 일반적인 정보성 메시지를 나타냅니다.
	INFO
	// WARN은 잠재적인 문제를 나타냅니다.
	WARN
	// ERROR는 오류 상황을 나타냅니다.
	ERROR
	// FATAL은 애플리케이션을 중단시킬 수 있는 심각한 오류를 나타냅니다.
	FATAL
)

// String은 로그 레벨을 문자열로 변환합니다.
func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}[l]
}

// IsValid는 로그 레벨이 유효한지 검사합니다.
func (l LogLevel) IsValid() bool {
	return l >= DEBUG && l <= FATAL
}

```
## internal/logger/domain/repository/log_repository.go
```go
package repository

import (
	"path/filepath"

	"github.com/kihyun1998/dupdater/internal/logger/domain/entity"
)

// LogRepository는 로그 저장소의 인터페이스를 정의합니다.
type LogRepository interface {
	// Write는 로그 엔트리를 저장합니다.
	Write(entry *entity.LogEntry) error

	// Rotate는 로그 파일을 순환합니다.
	Rotate() error

	// Close는 로그 저장소를 정리합니다.
	Close() error
}

// LogConfig는 로그 저장소의 설정을 정의합니다.
type LogConfig struct {
	// LogPath는 로그 파일의 경로입니다.
	LogPath string
	// MaxSize는 로그 파일의 최대 크기(바이트)입니다.
	MaxSize int64
	// MaxBackups는 보관할 최대 백업 파일 수입니다.
	MaxBackups int
}

// Validate는 로그 설정이 유효한지 검사합니다.
func (c *LogConfig) Validate() error {
	if c.LogPath == "" {
		return filepath.ErrBadPattern
	}
	if c.MaxSize <= 0 {
		c.MaxSize = 10 * 1024 * 1024 // 기본값 10MB
	}
	if c.MaxBackups <= 0 {
		c.MaxBackups = 5 // 기본값 5개
	}
	return nil
}

// NewLogConfig는 새로운 LogConfig를 생성합니다.
func NewLogConfig(path string, maxSize int64, maxBackups int) *LogConfig {
	return &LogConfig{
		LogPath:    path,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
	}
}

```
## internal/logger/domain/usecase/logger_service.go
```go
package usecase

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/kihyun1998/dupdater/internal/logger/domain/entity"
	"github.com/kihyun1998/dupdater/internal/logger/domain/repository"
)

// LoggerService는 로깅 서비스를 제공합니다.
type LoggerService struct {
	mu         sync.RWMutex
	repository repository.LogRepository
	level      entity.LogLevel
}

// NewLoggerService는 새로운 LoggerService 인스턴스를 생성합니다.
func NewLoggerService(repo repository.LogRepository, level entity.LogLevel) *LoggerService {
	return &LoggerService{
		repository: repo,
		level:      level,
	}
}

// log는 실제 로깅을 수행하는 내부 메서드입니다.
func (l *LoggerService) log(level entity.LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// 호출자 정보 가져오기
	_, file, line, ok := runtime.Caller(2)
	callerInfo := "unknown"
	if ok {
		callerInfo = fmt.Sprintf("%s:%d", file, line)
	}

	// 로그 엔트리 생성
	entry := entity.NewLogEntry(
		level,
		fmt.Sprintf(format, args...),
		callerInfo,
	)

	// 로그 저장
	if err := l.repository.Write(entry); err != nil {
		fmt.Printf("로그 저장 실패: %v\n", err)
	}

	// FATAL 레벨인 경우 프로그램 종료
	if level == entity.FATAL {
		l.Close()
		os.Exit(1)
	}
}

// 공개 메서드들 구현
func (l *LoggerService) Debug(format string, args ...interface{}) {
	l.log(entity.DEBUG, format, args...)
}

func (l *LoggerService) Info(format string, args ...interface{}) {
	l.log(entity.INFO, format, args...)
}

func (l *LoggerService) Warn(format string, args ...interface{}) {
	l.log(entity.WARN, format, args...)
}

func (l *LoggerService) Error(format string, args ...interface{}) {
	l.log(entity.ERROR, format, args...)
}

func (l *LoggerService) Fatal(format string, args ...interface{}) {
	l.log(entity.FATAL, format, args...)
}

func (l *LoggerService) Close() error {
	return l.repository.Close()
}

```
## internal/logger/factory.go
```go
package logger

import (
	"github.com/kihyun1998/dupdater/internal/logger/domain/entity"
	"github.com/kihyun1998/dupdater/internal/logger/domain/repository"
	"github.com/kihyun1998/dupdater/internal/logger/domain/usecase"
	"github.com/kihyun1998/dupdater/internal/logger/infrastructure"
)

// Config는 로거 생성에 필요한 설정을 정의합니다.
type Config struct {
	LogPath    string          // 로그 파일 경로
	LogLevel   entity.LogLevel // 로그 레벨
	MaxSize    int64           // 최대 파일 크기 (바이트)
	MaxBackups int             // 최대 백업 파일 수
}

// Logger는 로깅을 위한 인터페이스입니다.
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Fatal(format string, v ...interface{})
	Close() error
}

// New는 새로운 로거 인스턴스를 생성합니다.
func New(config Config) (Logger, error) {
	// 저장소 설정 생성
	repoConfig := repository.NewLogConfig(
		config.LogPath,
		config.MaxSize,
		config.MaxBackups,
	)

	// 파일 로거 생성
	fileLogger, err := infrastructure.NewFileLogger(repoConfig)
	if err != nil {
		return nil, err
	}

	// 로깅 서비스 생성 및 반환
	return usecase.NewLoggerService(fileLogger, config.LogLevel), nil
}

```
## internal/logger/infrastructure/file_logger.go
```go
package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/kihyun1998/dupdater/internal/logger/domain/entity"
	"github.com/kihyun1998/dupdater/internal/logger/domain/repository"
)

// FileLogger는 파일 기반 로그 저장소입니다.
type FileLogger struct {
	mu       sync.RWMutex
	config   *repository.LogConfig
	file     *os.File
	fileSize int64
}

// NewFileLogger는 새로운 FileLogger를 생성합니다.
func NewFileLogger(config *repository.LogConfig) (repository.LogRepository, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("로거 설정 검증 실패: %w", err)
	}

	// 로그 디렉토리 생성
	if err := os.MkdirAll(filepath.Dir(config.LogPath), 0755); err != nil {
		return nil, fmt.Errorf("로그 디렉토리 생성 실패: %w", err)
	}

	// 로그 파일 생성 또는 열기
	file, err := os.OpenFile(
		config.LogPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, fmt.Errorf("로그 파일 열기 실패: %w", err)
	}

	// 현재 파일 크기 확인
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("파일 정보 가져오기 실패: %w", err)
	}

	return &FileLogger{
		config:   config,
		file:     file,
		fileSize: info.Size(),
	}, nil
}

// Write는 로그 엔트리를 파일에 기록합니다.
func (f *FileLogger) Write(entry *entity.LogEntry) error {
	if !entry.IsValid() {
		return fmt.Errorf("유효하지 않은 로그 엔트리")
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	// 로그 문자열 생성
	logString := entry.Format() + "\n"
	logSize := int64(len(logString))

	// 파일 크기 체크 및 순환
	if f.fileSize+logSize > f.config.MaxSize {
		if err := f.rotate(); err != nil {
			return fmt.Errorf("로그 파일 순환 실패: %w", err)
		}
	}

	// 로그 기록
	n, err := f.file.WriteString(logString)
	if err != nil {
		return fmt.Errorf("로그 쓰기 실패: %w", err)
	}

	f.fileSize += int64(n)
	return nil
}

// Rotate는 로그 파일을 순환합니다.
func (f *FileLogger) Rotate() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.rotate()
}

// rotate는 내부적으로 로그 파일을 순환하는 헬퍼 메서드입니다.
func (f *FileLogger) rotate() error {
	// 현재 파일 닫기
	if err := f.file.Close(); err != nil {
		return fmt.Errorf("현재 로그 파일 닫기 실패: %w", err)
	}

	// 백업 파일 순환
	for i := f.config.MaxBackups - 1; i >= 0; i-- {
		oldPath := fmt.Sprintf("%s.%d", f.config.LogPath, i)
		newPath := fmt.Sprintf("%s.%d", f.config.LogPath, i+1)

		if i == 0 {
			oldPath = f.config.LogPath
		}

		// 마지막 백업 파일은 삭제
		if i == f.config.MaxBackups-1 {
			os.Remove(newPath)
			continue
		}

		// 파일 이름 변경
		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				return fmt.Errorf("파일 이름 변경 실패: %w", err)
			}
		}
	}

	// 새 로그 파일 생성
	file, err := os.OpenFile(
		f.config.LogPath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0644,
	)
	if err != nil {
		return fmt.Errorf("새 로그 파일 생성 실패: %w", err)
	}

	f.file = file
	f.fileSize = 0
	return nil
}

// Close는 로거를 정리합니다.
func (f *FileLogger) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.file != nil {
		if err := f.file.Sync(); err != nil {
			return fmt.Errorf("파일 동기화 실패: %w", err)
		}
		if err := f.file.Close(); err != nil {
			return fmt.Errorf("파일 닫기 실패: %w", err)
		}
		f.file = nil
	}
	return nil
}

```
## internal/network/domain/entity/server_config.go
```go
// Package entity는 네트워크 도메인의 핵심 개념을 정의합니다
package entity

import "fmt"

// ServerConfig는 서버 설정 정보를 담는 엔티티입니다
type ServerConfig struct {
	Profile Profile // 현재 사용중인 서버 프로필
}

// Profile은 서버 프로필 정보를 담는 값 객체입니다
type Profile struct {
	Name string // 서버 프로필명 (예: "server1")
	IP   string // 서버 IP 주소
}

// GetServerIP는 프로필의 서버 IP를 반환합니다
func (sc *ServerConfig) GetServerIP() (string, error) {
	if sc.Profile.IP == "" {
		return "", fmt.Errorf("서버 IP가 설정되지 않았습니다")
	}
	return sc.Profile.IP, nil
}

```
## internal/network/domain/entity/update_file.go
```go
package entity

// UpdateFile은 업데이트 파일 정보를 담는 엔티티입니다
type UpdateFile struct {
	Filename string // 업데이트 파일명
}

```
## internal/network/domain/repository/network_repository.go
```go
package repository

import (
	"io"

	"github.com/kihyun1998/dupdater/internal/network/domain/entity"
)

// NetworkRepository는 네트워크 작업을 추상화하는 인터페이스입니다
type NetworkRepository interface {
	// LoadServerConfig는 설정 파일에서 서버 설정을 로드합니다
	LoadServerConfig(configPath, serverName string) (*entity.ServerConfig, error)

	// FetchUpdateFileName은 서버로부터 업데이트 파일명을 조회합니다
	FetchUpdateFileName(serverIP string) (*entity.UpdateFile, error)

	// DownloadFile은 서버로부터 파일을 다운로드합니다
	// io.ReadCloser를 반환하여 스트림 처리가 가능하도록 합니다
	DownloadFile(serverIP string, filename string) (io.ReadCloser, error)
}

```
## internal/network/domain/usecase/network_service.go
```go
// Package usecase는 네트워크 작업의 비즈니스 로직을 구현합니다
package usecase

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/network/domain/repository"
)

// NetworkService는 네트워크 작업의 비즈니스 로직을 구현합니다
type NetworkService struct {
	repo   repository.NetworkRepository
	logger logger.Logger
}

// NewNetworkService는 새로운 NetworkService 인스턴스를 생성합니다
func NewNetworkService(repo repository.NetworkRepository, logger logger.Logger) *NetworkService {
	return &NetworkService{
		repo:   repo,
		logger: logger,
	}
}

// GetServerIP는 프로필명을 통해 서버 IP를 조회합니다
func (s *NetworkService) GetServerIP(serverName string) (string, error) {
	// 패닉 복구
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("GetServerIP 함수에서 패닉 발생: %v", r)
		}
	}()

	// 홈 디렉토리 경로 가져오기
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("홈 디렉토리 경로 가져오기 실패: %w", err)
	}

	// 설정 파일 경로 설정 및 로드
	configPath := filepath.Join(homeDir, ".testfolder", "config.json")
	config, err := s.repo.LoadServerConfig(configPath, serverName)
	if err != nil {
		return "", fmt.Errorf("서버 설정 로드 실패: %w", err)
	}

	serverIP, err := config.GetServerIP()
	if err != nil {
		return "", fmt.Errorf("서버 IP 가져오기 실패: %w", err)
	}

	return serverIP, nil
}

// GetUpdateFileName은 서버로부터 업데이트 파일명을 가져옵니다
func (s *NetworkService) GetUpdateFileName(serverIP string) (string, error) {
	// 패닉 복구
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("GetUpdateFileName 함수에서 패닉 발생: %v", r)
		}
	}()

	updateFile, err := s.repo.FetchUpdateFileName(serverIP)
	if err != nil {
		return "", fmt.Errorf("업데이트 파일명 가져오기 실패: %w", err)
	}

	if updateFile.Filename == "" {
		return "", fmt.Errorf("업데이트 파일명이 비어있습니다")
	}

	return updateFile.Filename, nil
}

// DownloadFile은 업데이트 파일을 다운로드합니다
func (s *NetworkService) DownloadFile(serverIP string, filename string) (io.ReadCloser, error) {
	reader, err := s.repo.DownloadFile(serverIP, filename)
	if err != nil {
		return nil, fmt.Errorf("파일 다운로드 실패: %w", err)
	}

	s.logger.Info("다운로드 시작 - 파일: %s", filename)
	return reader, nil
}

```
## internal/network/factory.go
```go
package network

import (
	"net/http"

	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/network/domain/usecase"
	"github.com/kihyun1998/dupdater/internal/network/infrastructure/repository"
)

// Manager는 네트워크 작업을 관리하는 인터페이스입니다
type Manager interface {
	GetServerIP(serverName string) (string, error)
	GetUpdateFileName(serverIP string) (string, error)
	DownloadFile(serverIP, filename string) (*http.Response, error)
}

// Config는 Manager 생성에 필요한 설정을 담는 구조체입니다
type Config struct {
	Logger logger.Logger
}

// networkManager는 Manager 인터페이스를 구현하는 구조체입니다
type networkManager struct {
	service *usecase.NetworkService
}

// New는 새로운 Manager 인스턴스를 생성합니다
func New(config Config) Manager {
	// HTTP Repository 생성
	repo := repository.NewHTTPRepository(config.Logger)

	// Network Service 생성
	service := usecase.NewNetworkService(repo, config.Logger)

	return &networkManager{
		service: service,
	}
}

// GetServerIP는 프로필명을 통해 서버 IP를 조회합니다
func (m *networkManager) GetServerIP(serverName string) (string, error) {
	return m.service.GetServerIP(serverName)
}

// GetUpdateFileName은 서버로부터 업데이트 파일명을 가져옵니다
func (m *networkManager) GetUpdateFileName(serverIP string) (string, error) {
	return m.service.GetUpdateFileName(serverIP)
}

// DownloadFile은 업데이트 파일을 다운로드합니다
func (m *networkManager) DownloadFile(serverIP, filename string) (*http.Response, error) {
	reader, err := m.service.DownloadFile(serverIP, filename)
	if err != nil {
		return nil, err
	}

	// io.ReadCloser를 http.Response로 변환
	return &http.Response{
		Body: reader,
		// 기존 코드와의 호환성을 위해 Response 객체로 감싸서 반환
	}, nil
}

```
## internal/network/infrastructure/repository/http_repository.go
```go
// Package repository는 네트워크 작업의 실제 구현체를 제공합니다
package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/network/domain/entity"
	"github.com/kihyun1998/dupdater/internal/network/domain/repository"
)

// HTTPRepository는 HTTP 기반의 네트워크 작업을 구현합니다
type HTTPRepository struct {
	client *http.Client
	logger logger.Logger
}

// NewHTTPRepository는 새로운 HTTPRepository 인스턴스를 생성합니다
func NewHTTPRepository(logger logger.Logger) repository.NetworkRepository {
	return &HTTPRepository{
		client: &http.Client{},
		logger: logger,
	}
}

// LoadServerConfig는 설정 파일에서 서버 설정을 로드합니다
func (r *HTTPRepository) LoadServerConfig(configPath, serverName string) (*entity.ServerConfig, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("설정 파일 읽기 실패: %w", err)
	}

	var config struct {
		ProfileList map[string]struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"profileList"`
	}

	if err := json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("설정 파일 파싱 실패: %w", err)
	}

	if profile, exists := config.ProfileList[serverName]; exists {
		return &entity.ServerConfig{
			Profile: entity.Profile{
				Name: profile.Name,
				IP:   profile.Address,
			},
		}, nil
	}

	return nil, fmt.Errorf("서버 프로필을 찾을 수 없음: %s", serverName)
}

// FetchUpdateFileName은 서버로부터 업데이트 파일명을 조회합니다
func (r *HTTPRepository) FetchUpdateFileName(serverIP string) (*entity.UpdateFile, error) {
	url := fmt.Sprintf("%s/update/updatefilename", serverIP)
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("서버 요청 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("서버 응답 오류: %d", resp.StatusCode)
	}

	var result struct {
		Filename string `json:"filename"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("응답 데이터 파싱 실패: %w", err)
	}

	return &entity.UpdateFile{
		Filename: result.Filename,
	}, nil
}

// DownloadFile은 서버로부터 파일을 다운로드합니다
func (r *HTTPRepository) DownloadFile(serverIP string, filename string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/update/file", serverIP)
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("파일 다운로드 요청 실패: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("파일 다운로드 응답 오류: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

```
## internal/scenario/entity/scenario.go
```go
// internal/scenario/entity/scenario.go
package entity

import "fmt"

// ScenarioType은 시나리오 유형을 정의합니다
type ScenarioType string

const (
	ScenarioSuccess ScenarioType = "success"
	ScenarioError1  ScenarioType = "error1"
	ScenarioError2  ScenarioType = "error2"
	ScenarioError3  ScenarioType = "error3"
	ScenarioError4  ScenarioType = "error4"
	ScenarioError5  ScenarioType = "error5"
	ScenarioError6  ScenarioType = "error6"
	ScenarioError7  ScenarioType = "error7"
	ScenarioError8  ScenarioType = "error8"
)

// Scenario는 테스트 시나리오의 기본 정보를 정의합니다
type Scenario struct {
	Type        ScenarioType // 시나리오 유형
	TargetStep  *Step        // 목표 단계 (에러 발생 단계)
	Description string       // 시나리오 설명
}

// NewScenario는 새로운 시나리오를 생성합니다
func NewScenario(scenarioType ScenarioType) *Scenario {
	s := &Scenario{
		Type: scenarioType,
	}

	// 시나리오 유형에 따른 설정
	switch scenarioType {
	case ScenarioSuccess:
		s.Description = "정상 업데이트 시나리오"
		s.TargetStep = nil
	default:
		// error1 ~ error8 처리
		stepIndex := int(scenarioType[5] - '1') // "error1"에서 숫자 추출
		if step := GetStep(stepIndex); step != nil {
			s.TargetStep = step
			s.Description = fmt.Sprintf("%s 단계 실패 시나리오", step.Name)
		}
	}

	return s
}

// IsErrorScenario는 현재 시나리오가 에러 시나리오인지 확인합니다
func (s *Scenario) IsErrorScenario() bool {
	return s.Type != ScenarioSuccess
}

```
## internal/scenario/entity/step.go
```go
package entity

// Step은 업데이트 프로세스의 각 단계를 정의합니다
type Step struct {
	Index       int    // 단계 인덱스 (0-7)
	Name        string // 단계 이름
	MessageKey  string // i18n 메시지 키
	NeedRestore bool   // 실패시 복원 필요 여부
}

// NewStep은 새로운 단계 정보를 생성합니다
func NewStep(index int, name string, messageKey string, needRestore bool) *Step {
	return &Step{
		Index:       index,
		Name:        name,
		MessageKey:  messageKey,
		NeedRestore: needRestore,
	}
}

// Steps는 전체 업데이트 단계 정보를 제공합니다
var Steps = []*Step{
	NewStep(0, "AppCheck", "update.status.checking", false),
	NewStep(1, "ServerInfo", "update.status.getting_info", false),
	NewStep(2, "Backup", "update.status.preparing", true),
	NewStep(3, "Download", "update.status.downloading", true),
	NewStep(4, "Verify", "update.status.verifying", true),
	NewStep(5, "Install", "update.status.installing", true),
	NewStep(6, "FinalVerify", "update.status.finalizing", true),
	NewStep(7, "Restart", "update.status.completed", true),
}

// GetStep은 인덱스에 해당하는 단계 정보를 반환합니다
func GetStep(index int) *Step {
	if index < 0 || index >= len(Steps) {
		return nil
	}
	return Steps[index]
}

```
## internal/scenario/factory.go
```go
// internal/scenario/factory.go
package scenario

import (
	"fmt"

	"github.com/kihyun1998/dupdater/internal/i18n"
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/scenario/entity"
	"github.com/kihyun1998/dupdater/internal/scenario/usecase"
	"github.com/kihyun1998/dupdater/internal/ui"
)

// Manager는 시나리오 실행을 위한 인터페이스입니다
type Manager interface {
	Run()
}

// Config는 시나리오 매니저 생성에 필요한 설정입니다
type Config struct {
	Logger       logger.Logger
	UIManager    ui.Manager
	I18n         i18n.Manager
	ScenarioType string
}

// manager는 시나리오 매니저의 구현체입니다
type manager struct {
	service *usecase.ScenarioService
}

// New는 새로운 시나리오 매니저를 생성합니다
func New(config Config) (Manager, error) {
	// 시나리오 타입 검증 및 생성
	scenarioType := entity.ScenarioType(config.ScenarioType)
	scenario := entity.NewScenario(scenarioType)

	if scenario == nil {
		return nil, fmt.Errorf("잘못된 시나리오 타입: %s", config.ScenarioType)
	}

	// 시나리오 서비스 생성
	service := usecase.NewScenarioService(
		config.Logger,
		config.UIManager,
		scenario,
	)

	return &manager{
		service: service,
	}, nil
}

// Run은 시나리오를 실행합니다
func (m *manager) Run() {
	m.service.Run()
}

```
## internal/scenario/scenarios/base.go
```go
package scenarios

import (
	"time"

	"github.com/kihyun1998/dupdater/internal/scenario/entity"
	"github.com/kihyun1998/dupdater/internal/ui"
)

// StepExecutor는 각 단계의 실행을 담당하는 인터페이스입니다
type StepExecutor interface {
	Execute(uiManager ui.Manager) error
	GetStep() *entity.Step
}

// BaseScenario는 기본 시나리오 구현을 제공합니다
type BaseScenario struct {
	executors []StepExecutor
}

// Execute는 시나리오의 단계들을 순차적으로 실행합니다
func (b *BaseScenario) Execute(uiManager ui.Manager) error {
	for _, executor := range b.executors {
		step := executor.GetStep()
		uiManager.SetCurrentStep(step.Index)
		time.Sleep(2 * time.Second) // 단계 실행 시뮬레이션

		if err := executor.Execute(uiManager); err != nil {
			return err
		}
	}
	return nil
}

```
## internal/scenario/scenarios/error.go
```go
// internal/scenario/scenarios/error.go
package scenarios

import (
	"fmt"

	"github.com/kihyun1998/dupdater/internal/scenario/entity"
	"github.com/kihyun1998/dupdater/internal/ui"
)

// ErrorStep은 에러 발생 단계를 구현합니다
type ErrorStep struct {
	step      *entity.Step
	shouldErr bool
}

func NewErrorStep(step *entity.Step, shouldErr bool) *ErrorStep {
	return &ErrorStep{
		step:      step,
		shouldErr: shouldErr,
	}
}

func (s *ErrorStep) Execute(uiManager ui.Manager) error {
	if s.shouldErr {
		return fmt.Errorf("%s 단계에서 오류 발생", s.step.Name)
	}
	return nil
}

func (s *ErrorStep) GetStep() *entity.Step {
	return s.step
}

// ErrorScenario는 에러 시나리오를 구현합니다
type ErrorScenario struct {
	BaseScenario
	errorStep *entity.Step
}

// NewErrorScenario는 새로운 에러 시나리오를 생성합니다
func NewErrorScenario(errorStep *entity.Step) *ErrorScenario {
	scenario := &ErrorScenario{
		errorStep: errorStep,
	}

	// 에러 발생 단계까지의 실행자 추가
	for _, step := range entity.Steps {
		shouldErr := step.Index == errorStep.Index
		scenario.executors = append(scenario.executors, NewErrorStep(step, shouldErr))
		if shouldErr {
			break
		}
	}

	return scenario
}

```
## internal/scenario/scenarios/success.go
```go
package scenarios

import (
	"github.com/kihyun1998/dupdater/internal/scenario/entity"
	"github.com/kihyun1998/dupdater/internal/ui"
)

// SuccessStep은 성공 시나리오의 단계를 구현합니다
type SuccessStep struct {
	step *entity.Step
}

func NewSuccessStep(step *entity.Step) *SuccessStep {
	return &SuccessStep{step: step}
}

func (s *SuccessStep) Execute(uiManager ui.Manager) error {
	return nil // 성공 시나리오는 항상 성공
}

func (s *SuccessStep) GetStep() *entity.Step {
	return s.step
}

// SuccessScenario는 성공 시나리오를 구현합니다
type SuccessScenario struct {
	BaseScenario
}

// NewSuccessScenario는 새로운 성공 시나리오를 생성합니다
func NewSuccessScenario() *SuccessScenario {
	scenario := &SuccessScenario{}

	// 모든 단계를 성공 단계로 추가
	for _, step := range entity.Steps {
		scenario.executors = append(scenario.executors, NewSuccessStep(step))
	}

	return scenario
}

```
## internal/scenario/usecase/scenario_service.go
```go
package usecase

import (
	"fmt"
	"time"

	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/scenario/entity"
	"github.com/kihyun1998/dupdater/internal/ui"
)

// ScenarioService는 시나리오 실행을 담당하는 서비스입니다
type ScenarioService struct {
	logger    logger.Logger
	uiManager ui.Manager
	scenario  *entity.Scenario
}

// NewScenarioService는 새로운 ScenarioService를 생성합니다
func NewScenarioService(
	logger logger.Logger,
	uiManager ui.Manager,
	scenario *entity.Scenario,
) *ScenarioService {
	return &ScenarioService{
		logger:    logger,
		uiManager: uiManager,
		scenario:  scenario,
	}
}

// Run은 시나리오를 실행합니다
func (s *ScenarioService) Run() {
	s.logger.Info("시나리오 시작: %s", s.scenario.Description)

	// UI 복원 핸들러 설정
	s.uiManager.SetRestoreHandler(func() {
		s.handleRestore()
	})

	// 시나리오 실행
	go func() {
		if s.scenario.IsErrorScenario() {
			s.runErrorScenario()
		} else {
			s.runSuccessScenario()
		}
	}()

	// UI 실행
	s.uiManager.Run()
}

// runSuccessScenario는 성공 시나리오를 실행합니다
func (s *ScenarioService) runSuccessScenario() {
	for i, _ := range entity.Steps {
		s.uiManager.SetCurrentStep(i)
		time.Sleep(2 * time.Second) // 단계별 지연
	}
}

// runErrorScenario는 에러 시나리오를 실행합니다
func (s *ScenarioService) runErrorScenario() {
	targetStep := s.scenario.TargetStep

	// 목표 단계까지 정상 진행
	for i := 0; i <= targetStep.Index; i++ {
		s.uiManager.SetCurrentStep(i)

		if i == targetStep.Index {
			// 에러 발생 단계
			time.Sleep(1 * time.Second)
			errMsg := fmt.Sprintf("%s 단계에서 오류 발생", targetStep.Name)
			s.uiManager.ShowError(fmt.Errorf(errMsg))
			return
		}

		time.Sleep(2 * time.Second)
	}
}

// handleRestore는 복원 프로세스를 처리합니다
func (s *ScenarioService) handleRestore() {
	if !s.scenario.TargetStep.NeedRestore {
		return
	}

	s.logger.Info("복원 프로세스 시작")
	time.Sleep(3 * time.Second) // 복원 시간 시뮬레이션
	s.uiManager.ShowRestoreComplete()
}

```
## internal/ui/components/status_card.go
```go
package components

import (
	"fmt"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kihyun1998/dupdater/internal/i18n"
	"github.com/kihyun1998/dupdater/pkg/utils/theme"
)

// StatusCard는 현재 상태를 표시하는 컴포넌트입니다
type StatusCard struct {
	widget.BaseWidget
	container   *fyne.Container
	progressBar *widget.ProgressBar

	titleText    *canvas.Text
	versionText  *canvas.Text
	subtitleText *canvas.Text

	currentTheme theme.ThemeVariant
	i18n         i18n.Manager

	targetProgress  float64
	currentProgress float64
	animating       bool
	ticker          *time.Ticker

	onComplete func()
	onRestore  func()
}

// NewStatusCard는 새로운 StatusCard를 생성합니다
func NewStatusCard(fromVersion, toVersion string, themeVariant theme.ThemeVariant, i18n i18n.Manager) *StatusCard {
	card := &StatusCard{
		targetProgress:  0,
		currentProgress: 0,
		animating:       false,
		currentTheme:    themeVariant,
		i18n:            i18n,
	}
	card.ExtendBaseWidget(card)
	card.setupUI(fromVersion, toVersion)
	return card
}

// CreateRenderer는 Fyne 위젯 인터페이스를 구현합니다
func (s *StatusCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(s.container)
}

// setupUI는 UI 컴포넌트를 초기화합니다
func (s *StatusCard) setupUI(fromVersion, toVersion string) {
	// 타이틀 텍스트
	s.titleText = canvas.NewText(
		s.i18n.GetMessage("update.header.new_update"),
		s.currentTheme.TextColor(),
	)
	s.titleText.TextSize = s.currentTheme.FontSizeLarge()
	s.titleText.TextStyle = fyne.TextStyle{Bold: true}

	// 버전 정보
	versionFormat := fmt.Sprintf("%s → %s", fromVersion, toVersion)
	s.versionText = canvas.NewText(versionFormat, s.currentTheme.SubTextColor())
	s.versionText.TextSize = s.currentTheme.FontSizeMedium()

	// 진행 바
	s.progressBar = widget.NewProgressBar()
	s.progressBar.Resize(fyne.NewSize(350, 20))

	// 상태 메시지
	s.subtitleText = canvas.NewText(
		s.i18n.GetMessage("update.info.restart"),
		s.currentTheme.SubTextColor(),
	)
	s.subtitleText.TextSize = s.currentTheme.FontSizeSmall()

	// 레이아웃 구성
	s.container = container.NewVBox(
		container.NewVBox(
			s.titleText,
			s.versionText,
		),
		container.NewPadded(s.progressBar),
		container.NewPadded(s.subtitleText),
	)
}

// SetError는 에러 상태를 표시합니다
func (s *StatusCard) SetError(errMsg string) {
	s.titleText.Text = s.i18n.GetMessage("update.error.generic")
	s.titleText.Color = s.currentTheme.ErrorColor()
	s.subtitleText.Text = errMsg
	s.subtitleText.Color = s.currentTheme.ErrorColor()
	s.progressBar.Hide()
	s.Refresh()
}

// SetRestoring는 복원 진행 상태를 표시합니다
func (s *StatusCard) SetRestoring(msg string) {
	s.titleText.Text = s.i18n.GetMessage("update.restore.in_progress")
	s.titleText.Color = s.currentTheme.WarningColor()
	s.subtitleText.Text = msg
	s.subtitleText.Color = s.currentTheme.WarningColor()
	s.progressBar.Show()
	s.Refresh()
}

// SetRestoreComplete는 복원 완료 상태를 표시합니다
func (s *StatusCard) SetRestoreComplete() {
	s.titleText.Text = s.i18n.GetMessage("update.restore.completed")
	s.titleText.Color = s.currentTheme.SuccessColor()
	s.subtitleText.Text = s.i18n.GetMessage("update.info.restart")
	s.subtitleText.Color = s.currentTheme.SuccessColor()
	s.progressBar.Hide()
	s.Refresh()
}

// UpdateStatus는 상태와 진행률을 업데이트합니다
func (s *StatusCard) UpdateStatus(progress float64, message string) {
	s.titleText.Color = s.currentTheme.TextColor()
	s.subtitleText.Text = message
	s.subtitleText.Color = s.currentTheme.SubTextColor()

	if progress >= 0 {
		s.targetProgress = progress
		s.animateProgress()
	}

	s.progressBar.Show()
	s.Refresh()
}

// SetProgress는 다운로드 진행률을 업데이트합니다
func (s *StatusCard) SetProgress(current, total int64) {
	progress := float64(current) / float64(total)
	s.progressBar.SetValue(progress)
	s.Refresh()
}

// SetCompletionCallback은 진행률 100% 도달 시 실행될 콜백을 설정합니다
func (s *StatusCard) SetCompletionCallback(callback func()) {
	s.onComplete = callback
}

// SetRestoreHandler는 복원 핸들러를 설정합니다
func (s *StatusCard) SetRestoreHandler(handler func()) {
	s.onRestore = handler
}

// Refresh는 위젯을 새로고침합니다
func (s *StatusCard) Refresh() {
	s.container.Refresh()
}

// 부드러운 진행률 업데이트를 위한 메서드
func (s *StatusCard) animateProgress() {
	if s.ticker != nil {
		s.ticker.Stop()
	}

	s.animating = true
	s.ticker = time.NewTicker(16 * time.Millisecond) // 약 60FPS

	go func() {
		defer s.ticker.Stop()

		for range s.ticker.C {
			if !s.animating {
				return
			}

			diff := s.targetProgress - s.currentProgress
			step := diff * 0.1

			if math.Abs(diff) < 0.001 {
				s.currentProgress = s.targetProgress
				s.progressBar.SetValue(s.currentProgress)
				s.progressBar.Refresh()
				s.animating = false

				// 애니메이션이 100%에 도달했을 때 콜백 실행
				if s.currentProgress >= 0.999 {
					if s.onComplete != nil {
						s.onComplete()
					}
				}
				return
			}

			s.currentProgress += step
			s.progressBar.SetValue(s.currentProgress)
			s.progressBar.Refresh()
		}
	}()
}

```
## internal/ui/domain/entity/ui_state.go
```go
package entity

import "fmt"

// UIState는 UI의 상태를 나타내는 도메인 엔티티입니다
type UIState struct {
	CurrentStep      int     // 현재 진행 단계
	TotalSteps       int     // 전체 단계 수
	Progress         float64 // 진행률 (0-1 사이 값)
	DetailMessage    string  // 상세 메시지
	FromVersion      string  // 시작 버전
	ToVersion        string  // 목표 버전
	IsError          bool    // 에러 상태 여부
	IsRestoring      bool    // 복원 중 상태 여부
	RestoreCompleted bool    // 복원 완료 상태
}

// NewUIState는 새로운 UIState 인스턴스를 생성합니다
func NewUIState(totalSteps int, fromVersion, toVersion string) *UIState {
	return &UIState{
		CurrentStep:      0,
		TotalSteps:       totalSteps,
		Progress:         0.0,
		FromVersion:      fromVersion,
		ToVersion:        toVersion,
		IsError:          false,
		IsRestoring:      false,
		RestoreCompleted: false,
	}
}

// UpdateProgress는 현재 진행 상태를 업데이트합니다
func (s *UIState) UpdateProgress(step int, message string) {
	s.CurrentStep = step
	s.DetailMessage = message
	if s.TotalSteps > 0 {
		s.Progress = float64(step) / float64(s.TotalSteps-1)
	}
}

// SetError는 에러 상태로 변경합니다
func (s *UIState) SetError() {
	s.IsError = true
}

// SetRestoring는 복원 중 상태로 변경합니다
func (s *UIState) SetRestoring() {
	s.IsRestoring = true
	s.IsError = false
}

// SetRestoreCompleted는 복원 완료 상태로 변경합니다
func (s *UIState) SetRestoreCompleted() {
	s.RestoreCompleted = true
	s.IsRestoring = false
}

// IsCompleted는 업데이트가 완료되었는지 확인합니다
func (s *UIState) IsCompleted() bool {
	return s.Progress >= 1.0
}

// Validate는 상태가 유효한지 검증합니다
func (s *UIState) Validate() error {
	if s.TotalSteps <= 0 {
		return fmt.Errorf("전체 단계 수는 0보다 커야 합니다")
	}
	if s.FromVersion == "" || s.ToVersion == "" {
		return fmt.Errorf("버전 정보가 누락되었습니다")
	}
	return nil
}

```
## internal/ui/domain/repository/ui_repository.go
```go
package repository

import (
	"github.com/kihyun1998/dupdater/internal/i18n"
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/ui/domain/entity"
	"github.com/kihyun1998/dupdater/pkg/utils/theme"
)

// UIRepository는 UI 상태 관리를 위한 저장소 인터페이스입니다
type UIRepository interface {
	// UIManager 인터페이스와 호환되는 메서드들
	SetCurrentStep(step int)
	UpdateDetail(message string)
	ShowError(err error)
	Run()
	Close()
	SetRestoreHandler(handler func())
	ShowRestoring()
	ShowRestoreComplete()
	GetTotalSteps() int
	SetCompletionCallback(func())

	// 내부 상태 관리를 위한 추가 메서드
	GetState() *entity.UIState
	UpdateState(state *entity.UIState)
}

// Config는 UI Repository 생성에 필요한 설정입니다
type Config struct {
	AppName     string
	TotalSteps  int
	FromVersion string
	ToVersion   string
	Logger      logger.Logger
	Theme       theme.ThemeVariant
	I18n        i18n.Manager
}

```
## internal/ui/domain/usecase/ui_service.go
```go
// internal/ui/domain/usecase/ui_service.go

package usecase

import (
	"sync"

	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/ui/domain/repository"
)

// UIService는 UI 관련 비즈니스 로직을 처리합니다
type UIService struct {
	repo   repository.UIRepository
	logger logger.Logger
	mu     sync.RWMutex
}

// NewUIService는 새로운 UIService 인스턴스를 생성합니다
func NewUIService(repo repository.UIRepository, logger logger.Logger) *UIService {
	return &UIService{
		repo:   repo,
		logger: logger,
	}
}

// SetCurrentStep은 현재 진행 단계를 설정합니다
func (s *UIService) SetCurrentStep(step int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state := s.repo.GetState()
	state.UpdateProgress(step, state.DetailMessage)
	s.repo.UpdateState(state)
	s.repo.SetCurrentStep(step)
	s.logger.Info("현재 단계 업데이트: %d", step)
}

// UpdateDetail은 상세 메시지를 업데이트합니다
func (s *UIService) UpdateDetail(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state := s.repo.GetState()
	state.DetailMessage = message
	s.repo.UpdateState(state)
	s.repo.UpdateDetail(message)
	s.logger.Info("상세 메시지 업데이트: %s", message)
}

// ShowError는 에러 상태를 표시합니다
func (s *UIService) ShowError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state := s.repo.GetState()
	state.SetError()
	s.repo.UpdateState(state)
	s.repo.ShowError(err)
	s.logger.Error("에러 발생: %v", err)
}

// Run은 UI를 실행합니다
func (s *UIService) Run() {
	s.logger.Info("UI 실행")
	s.repo.Run()
}

// Close는 UI를 종료합니다
func (s *UIService) Close() {
	s.logger.Info("UI 종료")
	s.repo.Close()
}

// SetRestoreHandler는 복원 핸들러를 설정합니다
func (s *UIService) SetRestoreHandler(handler func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.repo.SetRestoreHandler(handler)
	s.logger.Info("복원 핸들러 설정됨")
}

// ShowRestoring은 복원 중 상태를 표시합니다
func (s *UIService) ShowRestoring() {
	s.mu.Lock()
	defer s.mu.Unlock()

	state := s.repo.GetState()
	state.SetRestoring()
	s.repo.UpdateState(state)
	s.repo.ShowRestoring()
	s.logger.Info("복원 중 상태로 변경")
}

// ShowRestoreComplete는 복원 완료 상태를 표시합니다
func (s *UIService) ShowRestoreComplete() {
	s.mu.Lock()
	defer s.mu.Unlock()

	state := s.repo.GetState()
	state.SetRestoreCompleted()
	s.repo.UpdateState(state)
	s.repo.ShowRestoreComplete()
	s.logger.Info("복원 완료 상태로 변경")
}

// GetTotalSteps는 전체 단계 수를 반환합니다
func (s *UIService) GetTotalSteps() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.repo.GetTotalSteps()
}

// SetCompletionCallback은 완료 콜백을 설정합니다
func (s *UIService) SetCompletionCallback(callback func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.repo.SetCompletionCallback(callback)
	s.logger.Info("완료 콜백 설정됨")
}

```
## internal/ui/factory.go
```go
// internal/ui/factory.go

package ui

import (
	"github.com/kihyun1998/dupdater/internal/i18n"
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/ui/domain/repository"
	"github.com/kihyun1998/dupdater/internal/ui/domain/usecase"
	"github.com/kihyun1998/dupdater/internal/ui/infrastructure"
	"github.com/kihyun1998/dupdater/pkg/utils/theme"
)

// Config는 UI 생성에 필요한 설정입니다
type Config struct {
	AppName     string
	TotalSteps  int
	FromVersion string
	ToVersion   string
	Logger      logger.Logger
	Theme       theme.ThemeVariant
	I18n        i18n.Manager
}

// Manager는 updater.UIManager와 동일한 인터페이스를 제공합니다
type Manager interface {
	SetCurrentStep(step int)
	UpdateDetail(message string)
	ShowError(err error)
	Run()
	Close()
	SetRestoreHandler(handler func())
	ShowRestoring()
	ShowRestoreComplete()
	GetTotalSteps() int
	SetCompletionCallback(func())
}

// manager는 UI 관리자의 실제 구현체입니다
type manager struct {
	service *usecase.UIService
}

// New는 새로운 UI Manager를 생성합니다
func New(config Config) Manager {
	// repository 계층 초기화
	repo := infrastructure.NewFyneManager(&repository.Config{
		AppName:     config.AppName,
		TotalSteps:  config.TotalSteps,
		FromVersion: config.FromVersion,
		ToVersion:   config.ToVersion,
		Logger:      config.Logger,
		Theme:       config.Theme,
		I18n:        config.I18n,
	})

	// service 계층 초기화
	service := usecase.NewUIService(repo, config.Logger)

	return &manager{
		service: service,
	}
}

// Manager 인터페이스 구현
func (m *manager) SetCurrentStep(step int) {
	m.service.SetCurrentStep(step)
}

func (m *manager) UpdateDetail(message string) {
	m.service.UpdateDetail(message)
}

func (m *manager) ShowError(err error) {
	m.service.ShowError(err)
}

func (m *manager) Run() {
	m.service.Run()
}

func (m *manager) Close() {
	m.service.Close()
}

func (m *manager) SetRestoreHandler(handler func()) {
	m.service.SetRestoreHandler(handler)
}

func (m *manager) ShowRestoring() {
	m.service.ShowRestoring()
}

func (m *manager) ShowRestoreComplete() {
	m.service.ShowRestoreComplete()
}

func (m *manager) GetTotalSteps() int {
	return m.service.GetTotalSteps()
}

func (m *manager) SetCompletionCallback(callback func()) {
	m.service.SetCompletionCallback(callback)
}

```
## internal/ui/infrastructure/fyne_manager.go
```go
// internal/ui/infrastructure/fyne_manager.go

package infrastructure

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/kihyun1998/dupdater/internal/i18n"
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/ui/components"
	"github.com/kihyun1998/dupdater/internal/ui/domain/entity"
	"github.com/kihyun1998/dupdater/internal/ui/domain/repository"
	"github.com/kihyun1998/dupdater/pkg/utils/theme"
)

// FyneManager는 Fyne 기반 UI 구현체입니다
type FyneManager struct {
	app        fyne.App
	mainWindow fyne.Window
	statusCard *components.StatusCard
	state      *entity.UIState
	logger     logger.Logger
	i18n       i18n.Manager
	theme      theme.ThemeVariant
	mu         sync.RWMutex
}

// NewFyneManager는 새로운 FyneManager를 생성합니다
func NewFyneManager(config *repository.Config) *FyneManager {
	// Fyne 앱 생성
	fyneApp := app.New()
	customTheme := theme.NewCustomTheme(config.Theme)
	fyneApp.Settings().SetTheme(customTheme)

	// 상태 초기화
	state := entity.NewUIState(
		config.TotalSteps,
		config.FromVersion,
		config.ToVersion,
	)

	manager := &FyneManager{
		app:    fyneApp,
		state:  state,
		logger: config.Logger,
		i18n:   config.I18n,
		theme:  config.Theme,
	}

	// 메인 윈도우 초기화
	manager.initWindow()
	manager.initUI()

	return manager
}

// initWindow은 메인 윈도우를 초기화합니다
func (m *FyneManager) initWindow() {
	m.mainWindow = m.app.NewWindow(m.i18n.GetMessage("update.title"))
	m.mainWindow.Resize(fyne.NewSize(400, 150))
	m.mainWindow.CenterOnScreen()
	m.mainWindow.SetFixedSize(true)
}

// initUI는 UI 컴포넌트를 초기화합니다
func (m *FyneManager) initUI() {
	m.statusCard = components.NewStatusCard(
		m.state.FromVersion,
		m.state.ToVersion,
		m.theme,
		m.i18n,
	)

	content := container.NewPadded(m.statusCard)
	m.mainWindow.SetContent(content)
}

// GetState는 현재 UI 상태를 반환합니다
func (m *FyneManager) GetState() *entity.UIState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}

// UpdateState는 UI 상태를 업데이트합니다
func (m *FyneManager) UpdateState(state *entity.UIState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state = state
}

// SetCurrentStep은 현재 진행 단계를 설정합니다
func (m *FyneManager) SetCurrentStep(step int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	progress := float64(step) / float64(m.state.TotalSteps-1)
	m.statusCard.UpdateStatus(progress, m.state.DetailMessage)
}

// UpdateDetail은 상세 메시지를 업데이트합니다
func (m *FyneManager) UpdateDetail(message string) {
	m.statusCard.UpdateStatus(-1, message)
}

// ShowError는 에러 상태를 표시합니다
func (m *FyneManager) ShowError(err error) {
	m.statusCard.SetError(err.Error())
}

// Run은 UI를 실행합니다
func (m *FyneManager) Run() {
	m.mainWindow.ShowAndRun()
}

// Close는 UI를 종료합니다
func (m *FyneManager) Close() {
	m.mainWindow.Close()
}

// SetRestoreHandler는 복원 핸들러를 설정합니다
func (m *FyneManager) SetRestoreHandler(handler func()) {
	m.statusCard.SetRestoreHandler(handler)
}

// ShowRestoring은 복원 중 상태를 표시합니다
func (m *FyneManager) ShowRestoring() {
	m.statusCard.SetRestoring(m.i18n.GetMessage("update.restore.in_progress"))
}

// ShowRestoreComplete는 복원 완료 상태를 표시합니다
func (m *FyneManager) ShowRestoreComplete() {
	m.statusCard.SetRestoreComplete()
}

// GetTotalSteps는 전체 단계 수를 반환합니다
func (m *FyneManager) GetTotalSteps() int {
	return m.state.TotalSteps
}

// SetCompletionCallback은 완료 콜백을 설정합니다
func (m *FyneManager) SetCompletionCallback(callback func()) {
	m.statusCard.SetCompletionCallback(callback)
}

```
## internal/updater/domain/entity/update_config.go
```go
package entity

import "fmt"

// UpdateConfig는 업데이트에 필요한 기본 설정을 담는 엔티티입니다
type UpdateConfig struct {
	AppName     string // 업데이트할 애플리케이션 이름
	FromVersion string // 현재 버전
	ServerName  string // 서버 프로필 이름
}

// NewUpdateConfig는 새로운 UpdateConfig 인스턴스를 생성합니다
func NewUpdateConfig(appName, fromVersion, serverName string) (*UpdateConfig, error) {
	if appName == "" {
		return nil, fmt.Errorf("앱 이름은 필수입니다")
	}
	if fromVersion == "" {
		return nil, fmt.Errorf("현재 버전은 필수입니다")
	}
	if serverName == "" {
		return nil, fmt.Errorf("서버 이름은 필수입니다")
	}

	return &UpdateConfig{
		AppName:     appName,
		FromVersion: fromVersion,
		ServerName:  serverName,
	}, nil
}

// Validate는 설정값의 유효성을 검증합니다
func (c *UpdateConfig) Validate() error {
	if c.AppName == "" {
		return fmt.Errorf("앱 이름이 비어있습니다")
	}
	if c.FromVersion == "" {
		return fmt.Errorf("현재 버전이 비어있습니다")
	}
	if c.ServerName == "" {
		return fmt.Errorf("서버 이름이 비어있습니다")
	}
	return nil
}

```
## internal/updater/domain/entity/update_status.go
```go
package entity

// UpdateStatus는 업데이트 진행 상태를 추적하는 엔티티입니다
type UpdateStatus struct {
	BackupCompleted  bool   // 백업 완료 여부
	RestoreCompleted bool   // 복원 완료 여부
	ServerIP         string // 서버 IP 주소
	UpdateFile       string // 업데이트 파일 경로
}

// NewUpdateStatus는 새로운 UpdateStatus 인스턴스를 생성합니다
func NewUpdateStatus() *UpdateStatus {
	return &UpdateStatus{
		BackupCompleted:  false,
		RestoreCompleted: false,
	}
}

// SetServerIP는 서버 IP를 설정합니다
func (s *UpdateStatus) SetServerIP(ip string) {
	s.ServerIP = ip
}

// SetUpdateFile은 업데이트 파일 경로를 설정합니다
func (s *UpdateStatus) SetUpdateFile(path string) {
	s.UpdateFile = path
}

// MarkBackupCompleted는 백업 완료 상태를 설정합니다
func (s *UpdateStatus) MarkBackupCompleted() {
	s.BackupCompleted = true
}

// MarkRestoreCompleted는 복원 완료 상태를 설정합니다
func (s *UpdateStatus) MarkRestoreCompleted() {
	s.RestoreCompleted = true
}

// NeedsRestore는 복원이 필요한 상태인지 확인합니다
func (s *UpdateStatus) NeedsRestore() bool {
	return s.BackupCompleted && !s.RestoreCompleted
}

```
## internal/updater/domain/repository/update_repository.go
```go
package repository

import (
	"github.com/kihyun1998/dupdater/internal/file"
	"github.com/kihyun1998/dupdater/internal/hash"
	"github.com/kihyun1998/dupdater/internal/i18n"
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/network"
	"github.com/kihyun1998/dupdater/internal/ui"
	"github.com/kihyun1998/dupdater/internal/updater/domain/entity"
)

// UpdateRepository는 업데이트 작업의 인터페이스를 정의합니다
type UpdateRepository interface {
	// CheckRunningApp은 대상 애플리케이션의 실행 상태를 확인합니다
	CheckRunningApp() error

	// GetServerIP는 서버 IP를 조회합니다
	GetServerIP(serverName string) (string, error)

	// BackupFiles는 현재 파일들을 백업합니다
	BackupFiles() error

	// DownloadUpdateFile은 업데이트 파일을 다운로드합니다
	DownloadUpdateFile(serverIP string) (string, error)

	// VerifyUpdateFile은 다운로드된 업데이트 파일을 검증합니다
	VerifyUpdateFile(filePath string) error

	// ExtractUpdateFile은 업데이트 파일을 압축 해제합니다
	ExtractUpdateFile(filePath string) error

	// VerifyExtractedFiles는 압축 해제된 파일들을 검증합니다
	VerifyExtractedFiles() error

	// RestartApplication은 애플리케이션을 재시작합니다
	RestartApplication(config *entity.UpdateConfig) error

	// RestartAfterRestore는 복원 완료 후 애플리케이션을 재시작합니다
	RestartAfterRestore(config *entity.UpdateConfig) error

	// RestoreFiles는 백업된 파일들을 복원합니다
	RestoreFiles() error
}

// Config는 저장소 생성에 필요한 설정을 정의합니다
type Config struct {
	Logger      logger.Logger        // 로거
	Status      *entity.UpdateStatus // 업데이트 상태
	Config      *entity.UpdateConfig // 업데이트 설정
	UI          ui.Manager
	Network     network.Manager
	FileManager file.Manager
	HashManager hash.Manager
	I18n        i18n.Manager
}

```
## internal/updater/domain/usecase/update_service.go
```go
package usecase

import (
	"fmt"
	"time"

	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/ui"
	"github.com/kihyun1998/dupdater/internal/updater/domain/entity"
	"github.com/kihyun1998/dupdater/internal/updater/domain/repository"
)

// UpdateService는 업데이트 프로세스의 비즈니스 로직을 구현합니다
type UpdateService struct {
	repo   repository.UpdateRepository
	ui     ui.Manager
	logger logger.Logger
	status *entity.UpdateStatus
	config *entity.UpdateConfig
}

// NewUpdateService는 새로운 UpdateService 인스턴스를 생성합니다
func NewUpdateService(
	repo repository.UpdateRepository,
	ui ui.Manager,
	logger logger.Logger,
	status *entity.UpdateStatus,
	config *entity.UpdateConfig,
) *UpdateService {
	return &UpdateService{
		repo:   repo,
		ui:     ui,
		logger: logger,
		status: status,
		config: config,
	}
}

// Start는 업데이트 프로세스를 시작합니다
func (s *UpdateService) Start() {
	// UI 복구 핸들러 설정
	s.ui.SetRestoreHandler(func() {
		if err := s.restoreFiles(); err != nil {
			s.logger.Error("파일 복원 실패: %v", err)
			s.ui.ShowError(fmt.Errorf("복원 실패: %w", err))
		}
	})

	// 업데이트 프로세스 시작
	go func() {
		if err := s.processUpdate(); err != nil {
			s.logger.Error("업데이트 실패: %v", err)
			s.handleError("업데이트 실패", err)
		}
	}()

	// UI 실행
	s.ui.Run()
}

// processUpdate는 실제 업데이트 작업을 수행합니다
func (s *UpdateService) processUpdate() error {
	s.logger.Info("업데이트 프로세스 시작")

	// 1. 애플리케이션 실행 상태 확인
	if err := s.repo.CheckRunningApp(); err != nil {
		return s.handleError("애플리케이션 상태 확인 실패", err)
	}

	// 2. 서버 IP 가져오기
	serverIP, err := s.repo.GetServerIP(s.config.ServerName)
	if err != nil {
		return s.handleError("서버 IP 가져오기 실패", err)
	}
	s.status.SetServerIP(serverIP)

	// 3. 파일 백업
	if err := s.repo.BackupFiles(); err != nil {
		return s.handleError("파일 백업 실패", err)
	}
	s.status.MarkBackupCompleted()

	// 4. 업데이트 파일 다운로드
	updateFile, err := s.repo.DownloadUpdateFile(serverIP)
	if err != nil {
		return s.handleError("업데이트 파일 다운로드 실패", err)
	}
	s.status.SetUpdateFile(updateFile)

	// 5. 업데이트 파일 검증
	if err := s.repo.VerifyUpdateFile(updateFile); err != nil {
		return s.handleError("업데이트 파일 검증 실패", err)
	}

	// 6. 파일 압축해제
	if err := s.repo.ExtractUpdateFile(updateFile); err != nil {
		return s.handleError("파일 압축해제 실패", err)
	}

	// 7. 압축해제된 파일들 검증
	if err := s.repo.VerifyExtractedFiles(); err != nil {
		return s.handleError("압축해제된 파일 검증 실패", err)
	}

	// 8. 애플리케이션 재시작
	if err := s.repo.RestartApplication(s.config); err != nil {
		return s.handleError("애플리케이션 재시작 실패", err)
	}

	return nil
}

// restoreFiles는 백업된 파일들을 복원합니다
func (s *UpdateService) restoreFiles() error {
	if s.status.RestoreCompleted {
		s.logger.Info("이미 복원이 완료되었습니다")
		return nil
	}

	s.logger.Info("파일 복원 시작")
	if err := s.repo.RestoreFiles(); err != nil {
		return err
	}

	s.status.MarkRestoreCompleted()
	time.Sleep(2 * time.Second) // UI 메시지 표시를 위한 대기
	s.ui.ShowRestoreComplete()

	return nil
}

// handleError는 에러 상황을 처리합니다
func (s *UpdateService) handleError(message string, err error) error {
	s.logger.Error("%s: %v", message, err)

	if s.status.NeedsRestore() {
		s.ui.ShowRestoring()
		if restoreErr := s.restoreFiles(); restoreErr != nil {
			s.logger.Error("파일 복원 실패: %v", restoreErr)
			s.ui.ShowError(fmt.Errorf("복원 실패: %v", restoreErr))
			return fmt.Errorf("%s 및 복원 실패: %v", message, err)
		}

		// 복원 성공시 기존 버전으로 앱 재시작
		if err := s.repo.RestartAfterRestore(s.config); err != nil {
			s.logger.Error("복원 후 앱 재시작 실패: %v", err)
		}

		return fmt.Errorf("%s, 파일이 복원됨: %v", message, err)
	}

	s.ui.ShowError(fmt.Errorf("%s: %v", message, err))
	return fmt.Errorf("%s: %v", message, err)
}

```
## internal/updater/factory.go
```go
package updater

import (
	"fmt"

	"github.com/kihyun1998/dupdater/internal/file"
	"github.com/kihyun1998/dupdater/internal/hash"
	"github.com/kihyun1998/dupdater/internal/i18n"
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/network"
	"github.com/kihyun1998/dupdater/internal/ui"
	"github.com/kihyun1998/dupdater/internal/updater/domain/entity"
	"github.com/kihyun1998/dupdater/internal/updater/domain/repository"
	"github.com/kihyun1998/dupdater/internal/updater/domain/usecase"
	"github.com/kihyun1998/dupdater/internal/updater/infrastructure"
)

// Manager는 업데이트 관리를 위한 인터페이스입니다
type Manager interface {
	Start()
}

// Config는 업데이터 생성에 필요한 설정입니다
type Config struct {
	AppName          string
	FromVersion      string
	ServerName       string
	BackupCompleted  bool
	RestoreCompleted bool
	UIManager        ui.Manager
	Logger           logger.Logger
	NetworkManager   network.Manager
	FileManager      file.Manager
	HashManager      hash.Manager
	I18n             i18n.Manager
}

// manager는 Manager 인터페이스의 구현체입니다
type manager struct {
	service *usecase.UpdateService
}

// New는 새로운 업데이트 Manager를 생성합니다
func New(config Config) (Manager, error) {
	// 1. 엔티티 생성
	updateConfig, err := entity.NewUpdateConfig(
		config.AppName,
		config.FromVersion,
		config.ServerName,
	)
	if err != nil {
		return nil, fmt.Errorf("업데이트 설정 생성 실패: %w", err)
	}

	updateStatus := entity.NewUpdateStatus()
	updateStatus.BackupCompleted = config.BackupCompleted
	updateStatus.RestoreCompleted = config.RestoreCompleted

	// 2. Repository 계층 설정
	repoConfig := &repository.Config{
		Config:      updateConfig,
		Status:      updateStatus,
		UI:          config.UIManager,
		Logger:      config.Logger,
		Network:     config.NetworkManager,
		FileManager: config.FileManager,
		HashManager: config.HashManager,
		I18n:        config.I18n,
	}

	// 3. Repository 구현체 생성
	repo, err := infrastructure.NewUpdateManager(repoConfig)
	if err != nil {
		return nil, fmt.Errorf("업데이트 매니저 생성 실패: %w", err)
	}

	// 4. Service 계층 생성
	service := usecase.NewUpdateService(
		repo,
		config.UIManager,
		config.Logger,
		updateStatus,
		updateConfig,
	)

	return &manager{
		service: service,
	}, nil
}

// Start는 업데이트 프로세스를 시작합니다
func (m *manager) Start() {
	m.service.Start()
}

```
## internal/updater/infrastructure/update_manager.go
```go
package infrastructure

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kihyun1998/dupdater/internal/file"
	"github.com/kihyun1998/dupdater/internal/hash"
	"github.com/kihyun1998/dupdater/internal/i18n"
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/network"
	"github.com/kihyun1998/dupdater/internal/ui"
	"github.com/kihyun1998/dupdater/internal/updater/domain/entity"
	"github.com/kihyun1998/dupdater/internal/updater/domain/repository"
	"github.com/kihyun1998/dupdater/pkg/utils"
)

// UpdateManager는 실제 업데이트 작업을 수행하는 구현체입니다
type UpdateManager struct {
	config      *entity.UpdateConfig
	status      *entity.UpdateStatus
	ui          ui.Manager
	logger      logger.Logger
	network     network.Manager
	fileManager file.Manager
	hashManager hash.Manager
	i18n        i18n.Manager
}

// NewUpdateManager는 새로운 UpdateManager 인스턴스를 생성합니다
func NewUpdateManager(config *repository.Config) (*UpdateManager, error) {
	if err := config.Config.Validate(); err != nil {
		return nil, fmt.Errorf("설정 검증 실패: %w", err)
	}

	return &UpdateManager{
		config:      config.Config,
		status:      config.Status,
		ui:          config.UI,
		logger:      config.Logger,
		network:     config.Network,
		fileManager: config.FileManager,
		hashManager: config.HashManager,
		i18n:        config.I18n,
	}, nil
}

// CheckRunningApp은 대상 애플리케이션의 실행 상태를 확인합니다
func (m *UpdateManager) CheckRunningApp() error {
	m.ui.SetCurrentStep(0)
	m.ui.UpdateDetail(m.i18n.GetMessage("update.status.checking"))

	const maxAttempts = 30
	for attempt := 0; attempt < maxAttempts; attempt++ {
		isRunning, _ := utils.CheckApplicationRunning(m.config.AppName)
		if !isRunning {
			return nil
		}
		time.Sleep(time.Second)
	}

	return fmt.Errorf("앱 종료 대기 시간 초과")
}

// GetServerIP는 서버 IP를 조회합니다
func (m *UpdateManager) GetServerIP(serverName string) (string, error) {
	m.ui.SetCurrentStep(1)
	m.ui.UpdateDetail(m.i18n.GetMessage("update.status.getting_info"))

	return m.network.GetServerIP(serverName)
}

// BackupFiles는 현재 파일들을 백업합니다
func (m *UpdateManager) BackupFiles() error {
	m.ui.SetCurrentStep(2)
	m.ui.UpdateDetail(m.i18n.GetMessage("update.status.preparing"))
	return m.fileManager.Backup()
}

// DownloadUpdateFile은 업데이트 파일을 다운로드합니다
func (m *UpdateManager) DownloadUpdateFile(serverIP string) (string, error) {
	m.ui.SetCurrentStep(3)
	m.ui.UpdateDetail(m.i18n.GetMessage("update.status.downloading"))

	filename, err := m.network.GetUpdateFileName(serverIP)
	if err != nil {
		return "", err
	}

	resp, err := m.network.DownloadFile(serverIP, filename)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("파일 생성 실패: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("파일 쓰기 실패: %w", err)
	}

	return filename, nil
}

// VerifyUpdateFile은 다운로드된 업데이트 파일을 검증합니다
func (m *UpdateManager) VerifyUpdateFile(filePath string) error {
	m.ui.SetCurrentStep(4)
	m.ui.UpdateDetail(m.i18n.GetMessage("update.status.verifying"))
	return m.hashManager.VerifyUpdateFile(filePath)
}

// ExtractUpdateFile은 업데이트 파일을 압축 해제합니다
func (m *UpdateManager) ExtractUpdateFile(filePath string) error {
	m.ui.SetCurrentStep(5)
	m.ui.UpdateDetail(m.i18n.GetMessage("update.status.installing"))

	if err := m.fileManager.ExtractZip(filePath); err != nil {
		return err
	}

	return m.fileManager.DeleteFile(filePath)
}

// VerifyExtractedFiles는 압축 해제된 파일들을 검증합니다
func (m *UpdateManager) VerifyExtractedFiles() error {
	m.ui.SetCurrentStep(6)
	m.ui.UpdateDetail(m.i18n.GetMessage("update.status.finalizing"))
	return m.hashManager.VerifyHashSum()
}

// RestartApplication은 애플리케이션을 재시작합니다
func (m *UpdateManager) RestartApplication(config *entity.UpdateConfig) error {
	m.ui.SetCurrentStep(m.ui.GetTotalSteps() - 1)
	m.ui.UpdateDetail(m.i18n.GetMessage("update.status.completed"))

	completionCh := make(chan struct{})

	m.ui.SetCompletionCallback(func() {
		// 상대 경로 대신 실제 경로 사용
		appPath := fmt.Sprintf("./%s", config.AppName)
		args := []string{"--patch", "--fromVersion", config.FromVersion}

		if err := utils.LaunchApplication(appPath, args); err != nil {
			m.logger.Error("애플리케이션 재시작 실패: %v", err)
			return
		}

		time.Sleep(2 * time.Second)
		close(completionCh)
	})

	<-completionCh
	m.ui.Close()

	return nil
}

// RestartAfterRestore는 복원 완료 후 기존 버전의 애플리케이션을 재시작합니다
func (m *UpdateManager) RestartAfterRestore(config *entity.UpdateConfig) error {
	m.ui.SetCurrentStep(0)
	m.ui.UpdateDetail(m.i18n.GetMessage("update.restore.completed"))

	// 복원 완료 메시지 표시를 위한 대기
	time.Sleep(2 * time.Second)

	// 기존 버전으로 앱 실행
	appPath := fmt.Sprintf("./%s", config.AppName)

	if err := utils.LaunchApplication(appPath, nil); err != nil {
		m.logger.Error("복원 후 애플리케이션 실행 실패: %v", err)
		return err
	}

	// UI 종료 전 잠시 대기
	time.Sleep(2 * time.Second)
	m.ui.Close()

	return nil
}

// RestoreFiles는 백업된 파일들을 복원합니다
func (m *UpdateManager) RestoreFiles() error {
	if err := m.fileManager.Restore(); err != nil {
		m.logger.Error("파일 복원 실패: %v", err)
		return err
	}
	return nil
}

```
## internal/version/domain/entity/version.go
```go
package entity

import (
	"fmt"
	"regexp"
	"time"
)

// Version은 버전 정보를 나타내는 도메인 엔티티입니다
type Version struct {
	Major       int       // 주 버전
	Minor       int       // 부 버전
	Patch       int       // 패치 버전
	Date        time.Time // 배포 날짜
	FullVersion string    // 전체 버전 문자열
}

// NewVersion은 새로운 Version 엔티티를 생성합니다
func NewVersion(major, minor, patch int, date time.Time, fullVersion string) *Version {
	return &Version{
		Major:       major,
		Minor:       minor,
		Patch:       patch,
		Date:        date,
		FullVersion: fullVersion,
	}
}

// ParseVersion은 문자열을 Version 엔티티로 파싱합니다
// 예: V3.0.0(2024-01-01)
func ParseVersion(version string) (*Version, error) {
	pattern := `^V(\d+)\.(\d+)\.(\d+)\((\d{4}-\d{2}-\d{2})\)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(version)

	if matches == nil {
		return nil, fmt.Errorf("잘못된 버전 형식: %s (예: V3.0.0(2024-01-01))", version)
	}

	// 버전 번호 파싱
	var major, minor, patch int
	fmt.Sscanf(matches[1], "%d", &major)
	fmt.Sscanf(matches[2], "%d", &minor)
	fmt.Sscanf(matches[3], "%d", &patch)

	// 날짜 파싱
	date, err := time.Parse("2006-01-02", matches[4])
	if err != nil {
		return nil, fmt.Errorf("날짜 파싱 실패: %w", err)
	}

	return NewVersion(major, minor, patch, date, version), nil
}

// String은 Version을 문자열로 변환합니다
func (v *Version) String() string {
	return v.FullVersion
}

// Equal은 두 Version이 동일한지 비교합니다
func (v *Version) Equal(other *Version) bool {
	return v.Major == other.Major &&
		v.Minor == other.Minor &&
		v.Patch == other.Patch
}

// IsNewer는 현재 버전이 다른 버전보다 새로운지 확인합니다
func (v *Version) IsNewer(other *Version) bool {
	if v.Major != other.Major {
		return v.Major > other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor > other.Minor
	}
	if v.Patch != other.Patch {
		return v.Patch > other.Patch
	}
	if v.Date != other.Date {
		return v.Date.After(other.Date)
	}
	return false
}

// Validate는 버전 정보가 유효한지 검증합니다
func (v *Version) Validate() error {
	if v.Major < 0 || v.Minor < 0 || v.Patch < 0 {
		return fmt.Errorf("버전 번호는 음수일 수 없습니다")
	}
	if v.Date.IsZero() {
		return fmt.Errorf("날짜 정보가 없습니다")
	}
	if v.FullVersion == "" {
		return fmt.Errorf("전체 버전 문자열이 비어있습니다")
	}
	return nil
}

```
## internal/version/domain/repository/version_repo.go
```go
package repository

import (
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/version/domain/entity"
)

// VersionRepository는 버전 관리를 위한 저장소 인터페이스입니다
type VersionRepository interface {
	// GetFromVersion은 현재 버전 정보를 조회합니다
	GetFromVersion() *entity.Version

	// GetToVersion은 대상 버전 정보를 조회합니다
	GetToVersion() *entity.Version

	// ValidateVersions는 버전 정보의 유효성을 검증합니다
	ValidateVersions() error
}

// Config는 저장소 설정을 정의합니다
type Config struct {
	FromVersion string        // 현재 버전
	ToVersion   string        // 대상 버전
	Logger      logger.Logger // 로거 인터페이스
}

```
## internal/version/domain/usecase/version_service.go
```go
package usecase

import (
	"fmt"

	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/version/domain/repository"
)

// VersionService는 버전 관리의 비즈니스 로직을 구현합니다
type VersionService struct {
	repo   repository.VersionRepository
	logger logger.Logger
}

// NewVersionService는 새로운 VersionService 인스턴스를 생성합니다
func NewVersionService(repo repository.VersionRepository, logger logger.Logger) *VersionService {
	return &VersionService{
		repo:   repo,
		logger: logger,
	}
}

// GetFromVersion은 현재 버전을 반환합니다
func (s *VersionService) GetFromVersion() string {
	version := s.repo.GetFromVersion()
	if version == nil {
		s.logger.Error("현재 버전 정보를 찾을 수 없습니다")
		return ""
	}
	return version.String()
}

// GetToVersion은 대상 버전을 반환합니다
func (s *VersionService) GetToVersion() string {
	version := s.repo.GetToVersion()
	if version == nil {
		s.logger.Error("대상 버전 정보를 찾을 수 없습니다")
		return ""
	}
	return version.String()
}

// ValidateVersions는 버전 정보의 유효성을 검사합니다
func (s *VersionService) ValidateVersions() error {
	if err := s.repo.ValidateVersions(); err != nil {
		s.logger.Error("버전 정보 검증 실패: %v", err)
		return fmt.Errorf("버전 정보 검증 실패: %w", err)
	}
	return nil
}

```
## internal/version/factory.go
```go
package version

import (
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/version/domain/repository"
	"github.com/kihyun1998/dupdater/internal/version/domain/usecase"
	"github.com/kihyun1998/dupdater/internal/version/infrastructure"
)

// Manager는 버전 관리를 위한 외부 인터페이스입니다
type Manager interface {
	// GetFromVersion은 현재 버전을 반환합니다
	GetFromVersion() string

	// GetToVersion은 대상 버전을 반환합니다
	GetToVersion() string
}

// Config는 Manager 생성에 필요한 설정을 담는 구조체입니다
type Config struct {
	Logger      logger.Logger
	FromVersion string
	ToVersion   string
}

// managerImpl은 Manager 인터페이스의 구현체입니다
type managerImpl struct {
	service *usecase.VersionService
}

// New는 새로운 Manager 인스턴스를 생성합니다
func New(config Config) (Manager, error) {
	// 저장소 설정 생성
	repoConfig := &repository.Config{
		FromVersion: config.FromVersion,
		ToVersion:   config.ToVersion,
		Logger:      config.Logger,
	}

	// 버전 저장소 생성
	store, err := infrastructure.NewVersionStore(repoConfig)
	if err != nil {
		return nil, err
	}

	// 버전 서비스 생성
	service := usecase.NewVersionService(store, config.Logger)

	// 버전 정보 유효성 검증
	if err := service.ValidateVersions(); err != nil {
		return nil, err
	}

	return &managerImpl{
		service: service,
	}, nil
}

// GetFromVersion은 현재 버전을 반환합니다
func (m *managerImpl) GetFromVersion() string {
	return m.service.GetFromVersion()
}

// GetToVersion은 대상 버전을 반환합니다
func (m *managerImpl) GetToVersion() string {
	return m.service.GetToVersion()
}

```
## internal/version/infrastructure/version_store.go
```go
package infrastructure

import (
	"fmt"
	"sync"

	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/version/domain/entity"
	"github.com/kihyun1998/dupdater/internal/version/domain/repository"
)

// VersionStore는 버전 정보를 관리하는 저장소 구현체입니다
type VersionStore struct {
	fromVersion *entity.Version
	toVersion   *entity.Version
	logger      logger.Logger
	mu          sync.RWMutex
}

// NewVersionStore는 새로운 VersionStore 인스턴스를 생성합니다
func NewVersionStore(config *repository.Config) (*VersionStore, error) {
	fromVersion, err := entity.ParseVersion(config.FromVersion)
	if err != nil {
		return nil, fmt.Errorf("현재 버전 파싱 실패: %w", err)
	}

	toVersion, err := entity.ParseVersion(config.ToVersion)
	if err != nil {
		return nil, fmt.Errorf("대상 버전 파싱 실패: %w", err)
	}

	return &VersionStore{
		fromVersion: fromVersion,
		toVersion:   toVersion,
		logger:      config.Logger,
	}, nil
}

// GetFromVersion은 현재 버전을 반환합니다
func (s *VersionStore) GetFromVersion() *entity.Version {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.fromVersion
}

// GetToVersion은 대상 버전을 반환합니다
func (s *VersionStore) GetToVersion() *entity.Version {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.toVersion
}

// ValidateVersions는 버전 정보의 유효성을 검증합니다
func (s *VersionStore) ValidateVersions() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.fromVersion == nil || s.toVersion == nil {
		return fmt.Errorf("버전 정보가 초기화되지 않았습니다")
	}

	if err := s.fromVersion.Validate(); err != nil {
		return fmt.Errorf("현재 버전 검증 실패: %w", err)
	}

	if err := s.toVersion.Validate(); err != nil {
		return fmt.Errorf("대상 버전 검증 실패: %w", err)
	}

	if !s.toVersion.IsNewer(s.fromVersion) {
		return fmt.Errorf("대상 버전(%s)이 현재 버전(%s)보다 낮거나 같습니다",
			s.toVersion.String(), s.fromVersion.String())
	}

	return nil
}

```
## pkg/utils/count.go
```go
package utils

import (
	"fmt"
)

// WriteCounter는 다운로드 진행률을 추적하는 구조체입니다.
type WriteCounter struct {
	Total      int64                 // 전체 바이트 수
	Downloaded int64                 // 다운로드된 바이트 수
	OnProgress func(float64, string) // 진행률 업데이트 콜백
}

// Write는 io.Writer 인터페이스를 구현합니다.
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Downloaded += int64(n)
	percentage := float64(wc.Downloaded) / float64(wc.Total) * 100

	// 진행률 및 메시지 생성
	message := fmt.Sprintf("다운로드 중... %.1f%% (%.2f MB / %.2f MB)",
		percentage,
		float64(wc.Downloaded)/(1024*1024),
		float64(wc.Total)/(1024*1024))

	// 콜백 호출
	if wc.OnProgress != nil {
		wc.OnProgress(percentage, message)
	}

	return n, nil
}

// NewWriteCounter는 새로운 WriteCounter를 생성합니다.
func NewWriteCounter(total int64, onProgress func(float64, string)) *WriteCounter {
	return &WriteCounter{
		Total:      total,
		OnProgress: onProgress,
	}
}

```
## pkg/utils/fonts/fonts.go
```go
package fonts

// 폰트 리소스를 외부에서 접근 가능하도록 export
var (
	PretendardRegular = resourcePretendardRegularTtf
	PretendardBold    = resourcePretendardBoldTtf
	PretendardMedium  = resourcePretendardMediumTtf
)

```
## pkg/utils/system.go
```go
package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
)

// IsProcessRunning은 지정된 프로세스가 실행 중인지 확인합니다.
func IsProcessRunning(processName string) (bool, error) {
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s", processName))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("프로세스 확인 실패: %w", err)
	}

	return strings.Contains(string(output), processName), nil
}

// WaitForProcessToEnd는 프로세스가 종료될 때까지 대기합니다.
func WaitForProcessToEnd(processName string, timeout time.Duration, onWait func(attempt, maxAttempts int)) error {
	maxAttempts := int(timeout.Seconds())
	attempt := 0

	for {
		isRunning, err := IsProcessRunning(processName)
		if err != nil {
			return err
		}

		if !isRunning {
			return nil
		}

		attempt++
		if attempt >= maxAttempts {
			return fmt.Errorf("타임아웃: %s 프로세스가 %v 동안 종료되지 않음", processName, timeout)
		}

		if onWait != nil {
			onWait(attempt, maxAttempts)
		}

		time.Sleep(time.Second)
	}
}

// IsAdmin은 현재 프로세스가 관리자 권한으로 실행 중인지 확인합니다.
func IsAdmin() bool {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token := windows.Token(0)
	member, err := token.IsMember(sid)
	return err == nil && member
}

// CreateDirectoryIfNotExists는 디렉토리가 없으면 생성합니다.
func CreateDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

// LaunchProcessWithElevation은 프로세스를 관리자 권한으로 실행합니다.
func LaunchProcessWithElevation(execPath string, args string) error {
	verb := windows.StringToUTF16Ptr("runas") // verb도 StringToUTF16Ptr로 변환
	exe := windows.StringToUTF16Ptr(execPath)
	params := windows.StringToUTF16Ptr(args)
	dir := windows.StringToUTF16Ptr(filepath.Dir(execPath))

	err := windows.ShellExecute(0, verb, exe, params, dir, windows.SW_NORMAL)
	if err != nil {
		return fmt.Errorf("프로세스 실행 실패: %w", err)
	}

	return nil
}

// CheckApplicationRunning은 지정된 프로세스가 실행 중인지 확인합니다
func CheckApplicationRunning(appName string) (bool, error) {
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s", appName))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("프로세스 확인 중 오류: %v", err)
	}

	return strings.Contains(string(output), appName), nil
}

// WaitForApplicationToClose는 앱이 종료될 때까지 대기합니다
func WaitForApplicationToClose(appName string, onWait func(string)) error {
	const (
		maxAttempts = 30
		waitTime    = time.Second
	)

	for attempt := 0; attempt < maxAttempts; attempt++ {
		isRunning, err := CheckApplicationRunning(appName)
		if err != nil {
			return fmt.Errorf("프로세스 상태 확인 실패: %v", err)
		}

		if !isRunning {
			if onWait != nil {
				onWait("애플리케이션이 성공적으로 종료되었습니다.")
			}
			return nil
		}

		if onWait != nil {
			onWait(fmt.Sprintf("애플리케이션 종료 대기 중... (%d/%d)", attempt+1, maxAttempts))
		}

		time.Sleep(waitTime)
	}

	return fmt.Errorf("타임아웃: %d초 동안 애플리케이션이 종료되지 않았습니다", maxAttempts)
}

// LaunchApplication은 새 버전의 애플리케이션을 실행합니다
func LaunchApplication(execPath string, args []string) error {
	cmd := exec.Command(execPath, args...)
	// Windows 특정 설정 추가
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_NEW_CONSOLE,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("애플리케이션 실행 실패: %v", err)
	}

	return nil
}

// ElevateProcess는 프로세스를 관리자 권한으로 실행합니다
func ElevateProcess(execPath string, args string) error {
	verb := windows.StringToUTF16Ptr("runas")
	exe := windows.StringToUTF16Ptr(execPath)
	params := windows.StringToUTF16Ptr(args)
	dir := windows.StringToUTF16Ptr("")

	err := windows.ShellExecute(0, verb, exe, params, dir, windows.SW_NORMAL)
	if err != nil {
		return fmt.Errorf("관리자 권한으로 실행 실패: %v", err)
	}

	return nil
}

```
## pkg/utils/theme/dark_variant.go
```go
package theme

import "image/color"

// DarkVariant는 다크 모드 테마를 구현합니다
type DarkVariant struct{}

// NewDarkVariant는 새로운 다크 테마 인스턴스를 생성합니다
func NewDarkVariant() *DarkVariant {
	return &DarkVariant{}
}

// 기본 색상
func (t *DarkVariant) BackgroundColor() color.Color {
	return color.NRGBA{R: 12, G: 12, B: 19, A: 255} // #0C0C13
}

func (t *DarkVariant) TextColor() color.Color {
	return color.NRGBA{R: 229, G: 231, B: 235, A: 255} // #E5E7EB
}

func (t *DarkVariant) SubTextColor() color.Color {
	return color.NRGBA{R: 156, G: 163, B: 175, A: 255} // #9CA3AF
}

// 강조 색상
func (t *DarkVariant) PrimaryColor() color.Color {
	return color.NRGBA{R: 107, G: 105, B: 232, A: 255} // #6B69E8
}

// 상태 표시 색상
func (t *DarkVariant) SuccessColor() color.Color {
	return color.NRGBA{R: 34, G: 197, B: 94, A: 255} // #22C55E
}

func (t *DarkVariant) WarningColor() color.Color {
	return color.NRGBA{R: 234, G: 179, B: 8, A: 255} // #EAB308
}

func (t *DarkVariant) ErrorColor() color.Color {
	return color.NRGBA{R: 239, G: 68, B: 68, A: 255} // #EF4444
}

// 구분선 색상
func (t *DarkVariant) DividerColor() color.Color {
	return color.NRGBA{R: 31, G: 41, B: 55, A: 255} // #1F2937
}

// 폰트 크기
func (t *DarkVariant) FontSizeSmall() float32 {
	return DefaultFontSizeSmall
}

func (t *DarkVariant) FontSizeMedium() float32 {
	return DefaultFontSizeMedium
}

func (t *DarkVariant) FontSizeLarge() float32 {
	return DefaultFontSizeLarge
}

// 테마 모드
func (t *DarkVariant) IsLight() bool {
	return false
}

```
## pkg/utils/theme/light_variant.go
```go
package theme

import "image/color"

// LightVariant는 라이트 모드 테마를 구현합니다
type LightVariant struct{}

// NewLightVariant는 새로운 라이트 테마 인스턴스를 생성합니다
func NewLightVariant() *LightVariant {
	return &LightVariant{}
}

// 기본 색상
func (t *LightVariant) BackgroundColor() color.Color {
	return color.NRGBA{R: 255, G: 255, B: 255, A: 255} // #FFFFFF
}

func (t *LightVariant) TextColor() color.Color {
	return color.NRGBA{R: 17, G: 24, B: 39, A: 255} // #111827
}

func (t *LightVariant) SubTextColor() color.Color {
	return color.NRGBA{R: 107, G: 114, B: 128, A: 255} // #6B7280
}

// 강조 색상
func (t *LightVariant) PrimaryColor() color.Color {
	return color.NRGBA{R: 107, G: 105, B: 232, A: 255} // #6B69E8
}

// 상태 표시 색상
func (t *LightVariant) SuccessColor() color.Color {
	return color.NRGBA{R: 74, G: 222, B: 128, A: 255} // #4ADE80
}

func (t *LightVariant) WarningColor() color.Color {
	return color.NRGBA{R: 251, G: 191, B: 36, A: 255} // #FBBF24
}

func (t *LightVariant) ErrorColor() color.Color {
	return color.NRGBA{R: 248, G: 113, B: 113, A: 255} // #F87171
}

// 구분선 색상
func (t *LightVariant) DividerColor() color.Color {
	return color.NRGBA{R: 229, G: 231, B: 235, A: 255} // #E5E7EB
}

// 폰트 크기
func (t *LightVariant) FontSizeSmall() float32 {
	return DefaultFontSizeSmall
}

func (t *LightVariant) FontSizeMedium() float32 {
	return DefaultFontSizeMedium
}

func (t *LightVariant) FontSizeLarge() float32 {
	return DefaultFontSizeLarge
}

// 테마 모드
func (t *LightVariant) IsLight() bool {
	return true
}

```
## pkg/utils/theme/theme_variant.go
```go
package theme

import (
	"fmt"
	"image/color"
	"sync"
)

var (
	// currentVariant는 현재 사용 중인 테마 변형입니다
	currentVariant ThemeVariant

	// currentTheme은 현재 사용 중인 Fyne 테마입니다
	currentTheme *CustomTheme

	// themeMutex는 테마 변경 시 동시성을 제어합니다
	themeMutex sync.RWMutex
)

// InitTheme은 앱의 테마를 초기화합니다
func InitTheme(mode string) error {
	themeMutex.Lock()
	defer themeMutex.Unlock()

	switch mode {
	case "light":
		currentVariant = NewLightVariant()
	case "dark":
		currentVariant = NewDarkVariant()
	default:
		return fmt.Errorf("지원하지 않는 테마 모드: %s", mode)
	}

	currentTheme = NewCustomTheme(currentVariant)
	return nil
}

// GetCurrentTheme은 현재 Fyne 테마를 반환합니다
func GetCurrentTheme() *CustomTheme {
	themeMutex.RLock()
	defer themeMutex.RUnlock()

	if currentTheme == nil {
		// 기본값으로 라이트 테마 사용
		currentVariant = NewLightVariant()
		currentTheme = NewCustomTheme(currentVariant)
	}

	return currentTheme
}

// GetCurrentVariant는 현재 테마 변형을 반환합니다
func GetCurrentVariant() ThemeVariant {
	themeMutex.RLock()
	defer themeMutex.RUnlock()

	if currentVariant == nil {
		currentVariant = NewLightVariant()
	}

	return currentVariant
}

// SetTheme은 새로운 테마를 설정합니다
func SetTheme(variant ThemeVariant) {
	themeMutex.Lock()
	defer themeMutex.Unlock()

	currentVariant = variant
	currentTheme = NewCustomTheme(variant)
}

// 테마 모드 변경을 위한 헬퍼 함수들
func ToggleTheme() {
	themeMutex.Lock()
	defer themeMutex.Unlock()

	if currentVariant.IsLight() {
		currentVariant = NewDarkVariant()
	} else {
		currentVariant = NewLightVariant()
	}
	currentTheme = NewCustomTheme(currentVariant)
}

// IsLightMode는 현재 라이트 모드인지 확인합니다
func IsLightMode() bool {
	return GetCurrentVariant().IsLight()
}

// IsDarkMode는 현재 다크 모드인지 확인합니다
func IsDarkMode() bool {
	return !GetCurrentVariant().IsLight()
}

// 테마 색상을 직접 가져오는 유틸리티 함수들
func GetBackgroundColor() color.Color {
	return GetCurrentVariant().BackgroundColor()
}

func GetTextColor() color.Color {
	return GetCurrentVariant().TextColor()
}

func GetPrimaryColor() color.Color {
	return GetCurrentVariant().PrimaryColor()
}

func GetSubTextColor() color.Color {
	return GetCurrentVariant().SubTextColor()
}

func GetSuccessColor() color.Color {
	return GetCurrentVariant().SuccessColor()
}

func GetWarningColor() color.Color {
	return GetCurrentVariant().WarningColor()
}

func GetErrorColor() color.Color {
	return GetCurrentVariant().ErrorColor()
}

func GetDividerColor() color.Color {
	return GetCurrentVariant().DividerColor()
}

```
## pkg/utils/theme/types.go
```go
package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	baseTheme "fyne.io/fyne/v2/theme"
	"github.com/kihyun1998/dupdater/pkg/utils/fonts"
)

// ThemeVariant는 테마의 색상과 스타일을 정의하는 인터페이스입니다
type ThemeVariant interface {
	// 기본 색상
	BackgroundColor() color.Color
	TextColor() color.Color
	SubTextColor() color.Color

	// 강조 색상
	PrimaryColor() color.Color

	// 상태 표시 색상
	SuccessColor() color.Color
	WarningColor() color.Color
	ErrorColor() color.Color

	// 구분선 색상
	DividerColor() color.Color

	// 텍스트 크기
	FontSizeSmall() float32
	FontSizeMedium() float32
	FontSizeLarge() float32

	// 테마 모드
	IsLight() bool
}

// CustomTheme은 Fyne의 테마 인터페이스를 구현하는 커스텀 테마입니다
type CustomTheme struct {
	variant ThemeVariant
}

// CommonFontSizes는 공통으로 사용되는 폰트 크기를 정의합니다
const (
	DefaultFontSizeSmall  float32 = 14
	DefaultFontSizeMedium float32 = 16
	DefaultFontSizeLarge  float32 = 20
)

// NewCustomTheme은 새로운 CustomTheme 인스턴스를 생성합니다
func NewCustomTheme(variant ThemeVariant) *CustomTheme {
	return &CustomTheme{
		variant: variant,
	}
}

// Fyne Theme 인터페이스 구현
func (t *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case baseTheme.ColorNameBackground:
		return t.variant.BackgroundColor()
	case baseTheme.ColorNameForeground:
		return t.variant.TextColor()
	case baseTheme.ColorNamePrimary:
		return t.variant.PrimaryColor()
	case baseTheme.ColorNameError:
		return t.variant.ErrorColor()
	default:
		if t.variant.IsLight() {
			return baseTheme.LightTheme().Color(name, variant)
		}
		return baseTheme.DarkTheme().Color(name, variant)
	}
}

func (t *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	if t.variant.IsLight() {
		return baseTheme.LightTheme().Icon(name)
	}
	return baseTheme.DarkTheme().Icon(name)
}

// Font는 textStyle에 따른 폰트를 반환합니다
func (t *CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Bold {
		return fonts.PretendardBold
	}
	if style.Italic {
		return fonts.PretendardMedium
	}
	return fonts.PretendardRegular
}

func (t *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case baseTheme.SizeNameText:
		return t.variant.FontSizeMedium()
	case baseTheme.SizeNameHeadingText:
		return t.variant.FontSizeLarge()
	case baseTheme.SizeNameCaptionText:
		return t.variant.FontSizeSmall()
	default:
		if t.variant.IsLight() {
			return baseTheme.LightTheme().Size(name)
		}
		return baseTheme.DarkTheme().Size(name)
	}
}

```
