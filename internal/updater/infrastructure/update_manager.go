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
