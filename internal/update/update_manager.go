package update

import (
	"fmt"
	"time"
)

// UpdateManager는 업데이트 프로세스를 관리하는 구조체
type UpdateManager struct {
	fileManager    FileManager
	networkManager NetworkManager
	hashManager    HashManager
	processManager ProcessManager
	logger         Logger
	ui             UIUpdater
}

// FileManager는 파일 관리 인터페이스
type FileManager interface {
	BackupFiles(srcDir, backupDir string) error
	RestoreFiles(backupDir, targetDir string) error
	UnzipFile(zipPath, destDir string) error
	SafeRemove(path string) error
}

// NetworkManager는 네트워크 관리 인터페이스
type NetworkManager interface {
	GetServerIP(serverName string) (string, error)
	GetUpdateFileName(serverIP string) (string, error)
	DownloadFile(url, filepath string, progress chan<- DownloadProgress) error
}

// HashManager는 해시 관리 인터페이스
type HashManager interface {
	VerifyFile(filePath string) error
	VerifyHashSum(baseDir string) error
}

// ProcessManager는 프로세스 관리 인터페이스
type ProcessManager interface {
	IsProcessRunning(processName string) (bool, error)
	WaitForProcessToEnd(processName string, timeout time.Duration) error
	LaunchAppWithElevation(appPath string, args []string) error
}

// Logger는 로깅 인터페이스
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}

// UIUpdater는 UI 업데이트 인터페이스
type UIUpdater interface {
	SetCurrentStep(step int)
	UpdateDetail(message string)
	ShowError(err error)
}

// DownloadProgress는 다운로드 진행 상황 구조체
type DownloadProgress struct {
	Total      int64
	Downloaded int64
}

// Config는 업데이트 설정 구조체
type Config struct {
	AppName     string
	FromVersion string
	ServerName  string
	BackupDir   string
	DownloadDir string
	CurrentDir  string
}

// NewUpdateManager는 새로운 UpdateManager 인스턴스를 생성
func NewUpdateManager(
	fm FileManager,
	nm NetworkManager,
	hm HashManager,
	pm ProcessManager,
	logger Logger,
	ui UIUpdater,
) *UpdateManager {
	return &UpdateManager{
		fileManager:    fm,
		networkManager: nm,
		hashManager:    hm,
		processManager: pm,
		logger:         logger,
		ui:             ui,
	}
}

// StartUpdate는 업데이트 프로세스를 시작
func (um *UpdateManager) StartUpdate(config Config) error {
	um.logger.Info("업데이트 프로세스 시작")

	// 1. 서버 IP 획득
	um.ui.SetCurrentStep(1)
	um.ui.UpdateDetail("서버 IP 확인 중...")

	serverIP, err := um.networkManager.GetServerIP(config.ServerName)
	if err != nil {
		return fmt.Errorf("서버 IP 획득 실패: %w", err)
	}

	// 2. 앱 종료 대기
	um.ui.SetCurrentStep(2)
	um.ui.UpdateDetail("애플리케이션 종료 대기 중...")

	if err := um.processManager.WaitForProcessToEnd(config.AppName, 30*time.Second); err != nil {
		return fmt.Errorf("애플리케이션 종료 대기 실패: %w", err)
	}

	// 3. 파일 백업
	um.ui.SetCurrentStep(3)
	um.ui.UpdateDetail("파일 백업 중...")

	if err := um.fileManager.BackupFiles(config.CurrentDir, config.BackupDir); err != nil {
		return fmt.Errorf("파일 백업 실패: %w", err)
	}

	// 4. 업데이트 파일 다운로드
	um.ui.SetCurrentStep(4)
	um.ui.UpdateDetail("업데이트 파일 다운로드 중...")

	fileName, err := um.networkManager.GetUpdateFileName(serverIP)
	if err != nil {
		return um.handleError(err, config)
	}

	downloadURL := fmt.Sprintf("%s/update/file", serverIP)
	progress := make(chan DownloadProgress)

	// 다운로드 진행률 모니터링
	go um.monitorDownloadProgress(progress)

	if err := um.networkManager.DownloadFile(downloadURL, fileName, progress); err != nil {
		return um.handleError(err, config)
	}

	// 5. 압축 해제
	um.ui.SetCurrentStep(5)
	um.ui.UpdateDetail("파일 압축 해제 중...")

	if err := um.fileManager.UnzipFile(fileName, config.CurrentDir); err != nil {
		return um.handleError(err, config)
	}

	// 6. 해시 검증
	um.ui.SetCurrentStep(6)
	um.ui.UpdateDetail("파일 무결성 검증 중...")

	if err := um.hashManager.VerifyHashSum(config.CurrentDir); err != nil {
		return um.handleError(err, config)
	}

	// 7. 정리 및 새 버전 실행
	um.ui.SetCurrentStep(7)
	um.ui.UpdateDetail("업데이트 완료. 애플리케이션 실행 중...")

	if err := um.fileManager.SafeRemove(fileName); err != nil {
		um.logger.Error("임시 파일 삭제 실패: %v", err)
		// 임시 파일 삭제 실패는 치명적이지 않으므로 계속 진행
	}

	// 새 버전 실행
	args := []string{"--patch", "--fromVersion", config.FromVersion}
	if err := um.processManager.LaunchAppWithElevation(config.AppName, args); err != nil {
		return um.handleError(err, config)
	}

	um.logger.Info("업데이트 프로세스 완료")
	return nil
}

// monitorDownloadProgress는 다운로드 진행률을 모니터링
func (um *UpdateManager) monitorDownloadProgress(progress <-chan DownloadProgress) {
	for p := range progress {
		percentage := float64(p.Downloaded) / float64(p.Total) * 100
		um.ui.UpdateDetail(fmt.Sprintf("다운로드 중... %.1f%%", percentage))
	}
}

// handleError는 에러 처리 및 복구를 담당
func (um *UpdateManager) handleError(err error, config Config) error {
	um.logger.Error("업데이트 실패: %v", err)
	um.ui.UpdateDetail("오류 발생. 복구 중...")

	// 파일 복원 시도
	if restoreErr := um.fileManager.RestoreFiles(config.BackupDir, config.CurrentDir); restoreErr != nil {
		um.logger.Error("파일 복원 실패: %v", restoreErr)
		return fmt.Errorf("업데이트 및 복원 실패: %v (복원 오류: %v)", err, restoreErr)
	}

	return fmt.Errorf("업데이트 실패 (파일 복원됨): %v", err)
}
