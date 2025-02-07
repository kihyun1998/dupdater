// Package app은 업데이터의 핵심 어플리케이션 로직을 포함합니다
package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/kihyun1998/dupdater/internal/i18n"
	"github.com/kihyun1998/dupdater/pkg/utils"
	"golang.org/x/sys/windows"
)

const (
	appName  = "simple_update_test.exe"
	waitTime = 5 * time.Second // 충분한 대기 시간 설정
)

// Updater는 업데이트 프로세스의 전체 흐름을 제어하는 구조체입니다
type Updater struct {
	// 기본 설정
	appName          string // 업데이트할 애플리케이션의 이름
	fromVersion      string // 현재 애플리케이션의 버전
	serverName       string // 서버 프로필 이름
	backupCompleted  bool   // 백업 상태 추적을 위한 필드
	restoreCompleted bool   // 복원 상태 추적을 위한 필드

	// 의존성들
	ui          UIManager      // UI 관리자
	logger      Logger         // 로깅 시스템
	network     NetworkManager // 네트워크 관리자
	fileManager FileManager    // 파일 관리자
	hashManager HashManager    // 해시 관리자
	i18n        i18n.Manager   // 다국어 관리자

	// 상태 정보
	serverIP   string // 조회된 서버 IP
	updateFile string // 다운로드된 업데이트 파일 경로
}

// Config는 새로운 Updater를 생성할 때 필요한 설정을 담는 구조체입니다
type Config struct {
	AppName          string
	FromVersion      string
	ServerName       string
	BackupCompleted  bool
	RestoreCompleted bool
	UIManager        UIManager
	Logger           Logger
	NetworkManager   NetworkManager
	FileManager      FileManager
	HashManager      HashManager
	I18n             i18n.Manager
}

// New는 새로운 Updater 인스턴스를 생성합니다
func New(config Config) *Updater {
	return &Updater{
		appName:          config.AppName,
		fromVersion:      config.FromVersion,
		serverName:       config.ServerName,
		backupCompleted:  config.BackupCompleted,
		restoreCompleted: config.RestoreCompleted,
		ui:               config.UIManager,
		logger:           config.Logger,
		network:          config.NetworkManager,
		fileManager:      config.FileManager,
		hashManager:      config.HashManager,
		i18n:             config.I18n,
	}
}

// Start는 업데이트 프로세스를 시작합니다
func (u *Updater) Start() {
	// UI 복구 핸들러 설정
	u.ui.SetRestoreHandler(func() {
		if err := u.restoreFiles(); err != nil {
			u.logger.Error("파일 복원 실패: %v", err)
			u.ui.ShowError(fmt.Errorf("복원 실패: %w", err))
		} else {
			u.ui.UpdateDetail(u.i18n.GetMessage("update.restore.completed"))
		}
	})

	// 업데이트 프로세스 시작
	go func() {
		if err := u.processUpdate(); err != nil {
			u.logger.Error("업데이트 실패: %v", err)
			u.ui.ShowError(err)
		}
	}()

	// UI 실행 (메인 스레드에서 실행)
	u.ui.Run()
}

func (u *Updater) processUpdate() error {

	defer func() {
		if r := recover(); r != nil {
			u.logger.Error("업데이트 프로세스 중 패닉 발생: %v", r)
			u.handleError("예기치 못한 오류 발생", fmt.Errorf("%v", r))
		}
	}()

	// 1. 애플리케이션 실행 상태 확인
	if err := u.checkRunningApp(); err != nil {
		return u.handleError("애플리케이션 상태 확인 실패: %w", err)
	}

	// 2. 서버 IP 가져오기
	if err := u.getServerIP(); err != nil {
		return u.handleError("서버 IP 가져오기 실패: %w", err)
	}

	// 3. 파일 백업
	if err := u.backupFiles(); err != nil {
		return u.handleError("파일 백업 실패: %w", err)
	}
	u.backupCompleted = true // 백업 완료 상태 설정
	// 4. 업데이트 파일 다운로드
	if err := u.downloadUpdateFile(); err != nil {
		return u.handleError("업데이트 파일 다운로드 실패", err)
	}

	// 5. 업데이트 파일 검증
	if err := u.verifyUpdateFile(); err != nil {
		return u.handleError("업데이트 파일 검증 실패", err)
	}

	// 6. 파일 압축해제
	if err := u.extractUpdateFile(); err != nil {
		return u.handleError("파일 압축해제 실패", err)
	}

	// 7. 압축해제된 파일들 검증
	if err := u.verifyExtractedFiles(); err != nil {
		return u.handleError("압축해제된 파일 검증 실패", err)
	}

	// 8. 애플리케이션 재시작
	if err := u.restartApplication(); err != nil {
		return u.handleError("애플리케이션 재시작 실패", err)
	}

	return nil
}

func (u *Updater) checkRunningApp() error {
	const maxAttempts = 30
	u.ui.SetCurrentStep(0)
	u.ui.UpdateDetail(u.i18n.GetMessage("update.status.checking"))

	for attempt := 0; attempt < maxAttempts; attempt++ {
		isRunning, _ := utils.CheckApplicationRunning(appName)
		if !isRunning {
			return nil
		}
		time.Sleep(time.Second)
	}

	return fmt.Errorf("앱 종료 대기 시간 초과")
}

func (u *Updater) getServerIP() error {
	u.ui.SetCurrentStep(1)
	u.ui.UpdateDetail(u.i18n.GetMessage("update.status.getting_info"))

	ip, err := u.network.GetServerIP(u.serverName)
	if err != nil {
		return err
	}

	u.serverIP = ip
	return nil
}

func (u *Updater) backupFiles() error {
	u.ui.SetCurrentStep(2)
	u.ui.UpdateDetail(u.i18n.GetMessage("update.status.preparing"))
	return u.fileManager.Backup()
}

func (u *Updater) downloadUpdateFile() error {
	u.ui.SetCurrentStep(3)
	u.ui.UpdateDetail(u.i18n.GetMessage("update.status.downloading"))

	filename, err := u.network.GetUpdateFileName(u.serverIP)
	if err != nil {
		return err
	}

	resp, err := u.network.DownloadFile(u.serverIP, filename)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 파일 생성 및 저장 로직 추가
	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("파일 생성 실패: %v", err)
	}
	defer out.Close()

	// 파일 쓰기
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("파일 쓰기 실패: %v", err)
	}

	u.updateFile = filename
	return nil
}

func (u *Updater) verifyUpdateFile() error {
	u.ui.SetCurrentStep(4)
	u.ui.UpdateDetail(u.i18n.GetMessage("update.status.verifying"))
	return u.hashManager.VerifyUpdateFile(u.updateFile)
}

func (u *Updater) extractUpdateFile() error {
	u.ui.SetCurrentStep(5)
	u.ui.UpdateDetail(u.i18n.GetMessage("update.status.installing"))

	if err := u.fileManager.ExtractZip(u.updateFile); err != nil {
		return err
	}

	return u.fileManager.DeleteFile(u.updateFile)
}

func (u *Updater) verifyExtractedFiles() error {
	u.ui.SetCurrentStep(6)
	u.ui.UpdateDetail(u.i18n.GetMessage("update.status.finalizing"))
	return u.hashManager.VerifyHashSum()
}

func (u *Updater) restartApplication() error {
	// 마지막 단계 메시지 표시
	u.ui.SetCurrentStep(u.ui.GetTotalSteps() - 1)
	u.ui.UpdateDetail(u.i18n.GetMessage("update.status.completed"))

	// 프로그레스바가 100%까지 도달할 때까지 대기하기 위한 채널
	completionCh := make(chan struct{})

	// UI에 완료 콜백 설정
	u.ui.SetCompletionCallback(func() {
		// 앱 실행 준비
		cmd := exec.Command(fmt.Sprintf("./%s", u.appName), "--patch", "--fromVersion", u.fromVersion)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: windows.CREATE_NEW_CONSOLE,
		}

		// 앱 실행
		if err := cmd.Start(); err != nil {
			u.logger.Error("애플리케이션 실행 실패: %v", err)
			return
		}

		// 2초 대기 후 UI 종료
		time.Sleep(2 * time.Second)
		close(completionCh)
	})

	// 완료 대기
	<-completionCh
	u.ui.Close()

	return nil
}

// restartAfterRestore는 복원 완료 후 애플리케이션을 재시작합니다
func (u *Updater) restartAfterRestore() error {
	u.ui.SetCurrentStep(0)
	u.ui.UpdateDetail(u.i18n.GetMessage("update.restore.completed"))

	// 잠시 대기하여 메시지가 표시되도록 함
	time.Sleep(2 * time.Second)

	// 기존 버전으로 앱 실행 준비
	cmd := exec.Command(fmt.Sprintf("./%s", u.appName))
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_NEW_CONSOLE,
	}

	// 앱 실행
	if err := cmd.Start(); err != nil {
		u.logger.Error("복원 후 애플리케이션 실행 실패: %v", err)
		return err
	}

	// UI 종료 전 잠시 대기
	time.Sleep(2 * time.Second)
	u.ui.Close()

	return nil
}

func (u *Updater) restoreFiles() error {
	if u.restoreCompleted {
		u.logger.Error("restoreFiles()가 중복 실행되려 했으나, 이미 복원이 완료된 상태입니다. 중단함.")
		return nil
	}

	u.logger.Info("파일 복원 시작")
	if err := u.fileManager.Restore(); err != nil {
		u.logger.Error("파일 복원 실패: %v", err)
		return err
	}

	u.restoreCompleted = true

	// 복원 후 디렉토리 확인
	files, err := os.ReadDir(u.fileManager.GetBackupDir())
	if err != nil {
		u.logger.Error("복원 후 백업 디렉토리 확인 실패: %v", err)
	} else {
		for _, file := range files {
			u.logger.Info("복원 후 파일 확인: %s", file.Name())
		}
	}

	return nil
}

func (u *Updater) handleError(message string, err error) error {
	u.logger.Error("%s: %v", message, err)

	if u.backupCompleted && !u.restoreCompleted {
		u.ui.ShowRestoring()
		if restoreErr := u.restoreFiles(); restoreErr != nil {
			u.logger.Error("파일 복원 실패: %v", restoreErr)
			u.ui.ShowError(fmt.Errorf("복원 실패: %v", restoreErr))
			return fmt.Errorf("%s 및 복원 실패: %v", message, err)
		}
		u.ui.ShowRestoreComplete()

		if _, err := os.Stat(u.fileManager.GetBackupDir()); os.IsNotExist(err) {
			u.logger.Error("백업 디렉토리는 이미 삭제됨")
		} else {
			u.logger.Info("복원 완료 후 백업 디렉토리 유지됨: %s", u.fileManager.GetBackupDir())
		}

		// 복원 완료시 복원 완료 상태 저장
		u.restoreCompleted = true

		// 복원 성공 시 애플리케이션 재시작 추가
		u.logger.Info("복원이 완료됨. 애플리케이션을 재시작합니다...")
		if err := u.restartAfterRestore(); err != nil {
			u.logger.Error("복원 후 애플리케이션 재시작 실패: %v", err)
		}

		return fmt.Errorf("%s, 파일이 복원됨: %v", message, err)
	}

	u.ui.ShowError(fmt.Errorf("%s: %v", message, err))
	return fmt.Errorf("%s: %v", message, err)
}

// Logger는 로깅 작업을 위한 인터페이스입니다
type Logger interface {
	Info(format string, v ...interface{})  // 정보 로그를 기록합니다
	Error(format string, v ...interface{}) // 에러 로그를 기록합니다
}

// NetworkManager 인터페이스
type NetworkManager interface {
	GetServerIP(serverName string) (string, error)
	GetUpdateFileName(serverIP string) (string, error)
	DownloadFile(serverIP, filename string) (*http.Response, error)
}

// FileManager 인터페이스
type FileManager interface {
	Backup() error
	Restore() error
	ExtractZip(zipFile string) error
	DeleteFile(path string) error
	GetBackupDir() string
}

// HashManager 인터페이스
type HashManager interface {
	VerifyFile(filePath string, expectedHash string) error
	VerifyUpdateFile(filePath string) error
	VerifyHashSum() error
}

// UIManager 인터페이스
type UIManager interface {
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
