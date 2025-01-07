// Package app은 업데이터의 핵심 어플리케이션 로직을 포함합니다
package app

import (
	"fmt"
	"net/http"
	"time"
)

// Updater는 업데이트 프로세스의 전체 흐름을 제어하는 구조체입니다
type Updater struct {
	// 기본 설정
	appName     string // 업데이트할 애플리케이션의 이름
	fromVersion string // 현재 애플리케이션의 버전
	serverName  string // 서버 프로필 이름

	// 의존성들
	ui          UIManager      // UI 관리자
	logger      Logger         // 로깅 시스템
	network     NetworkManager // 네트워크 관리자
	fileManager FileManager    // 파일 관리자
	hashManager HashManager    // 해시 관리자

	// 상태 정보
	serverIP   string // 조회된 서버 IP
	updateFile string // 다운로드된 업데이트 파일 경로
}

// Config는 새로운 Updater를 생성할 때 필요한 설정을 담는 구조체입니다
type Config struct {
	AppName        string
	FromVersion    string
	ServerName     string
	UIManager      UIManager
	Logger         Logger
	NetworkManager NetworkManager
	FileManager    FileManager
	HashManager    HashManager
}

// New는 새로운 Updater 인스턴스를 생성합니다
func New(config Config) *Updater {
	return &Updater{
		appName:     config.AppName,
		fromVersion: config.FromVersion,
		serverName:  config.ServerName,
		ui:          config.UIManager,
		logger:      config.Logger,
		network:     config.NetworkManager,
		fileManager: config.FileManager,
		hashManager: config.HashManager,
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
			u.ui.UpdateDetail("파일 복원이 완료되었습니다")
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
		}
	}()

	// 1. 애플리케이션 실행 상태 확인
	if err := u.checkRunningApp(); err != nil {
		return fmt.Errorf("애플리케이션 상태 확인 실패: %w", err)
	}

	// 2. 서버 IP 가져오기
	if err := u.getServerIP(); err != nil {
		return fmt.Errorf("서버 IP 가져오기 실패: %w", err)
	}

	// 3. 파일 백업
	if err := u.backupFiles(); err != nil {
		return fmt.Errorf("파일 백업 실패: %w", err)
	}

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
	u.ui.SetCurrentStep(0)
	u.ui.UpdateDetail("애플리케이션 실행 상태를 확인하고 있습니다...")
	// 애플리케이션 실행 상태 확인 구현
	time.Sleep(1 * time.Second) // 임시 구현
	return nil
}

func (u *Updater) getServerIP() error {
	u.ui.SetCurrentStep(1)
	u.ui.UpdateDetail("서버 정보를 가져오고 있습니다...")

	ip, err := u.network.GetServerIP(u.serverName)
	if err != nil {
		return err
	}

	u.serverIP = ip
	return nil
}

func (u *Updater) backupFiles() error {
	u.ui.SetCurrentStep(2)
	u.ui.UpdateDetail("파일을 백업하고 있습니다...")
	return u.fileManager.Backup()
}

func (u *Updater) downloadUpdateFile() error {
	u.ui.SetCurrentStep(3)
	u.ui.UpdateDetail("업데이트 파일을 다운로드하고 있습니다...")

	filename, err := u.network.GetUpdateFileName(u.serverIP)
	if err != nil {
		return err
	}

	resp, err := u.network.DownloadFile(u.serverIP, filename)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	u.updateFile = filename
	return nil
}

func (u *Updater) verifyUpdateFile() error {
	u.ui.SetCurrentStep(4)
	u.ui.UpdateDetail("업데이트 파일을 검증하고 있습니다...")
	return u.hashManager.VerifyUpdateFile(u.updateFile)
}

func (u *Updater) extractUpdateFile() error {
	u.ui.SetCurrentStep(5)
	u.ui.UpdateDetail("파일을 압축해제하고 있습니다...")

	if err := u.fileManager.ExtractZip(u.updateFile); err != nil {
		return err
	}

	return u.fileManager.DeleteFile(u.updateFile)
}

func (u *Updater) verifyExtractedFiles() error {
	u.ui.SetCurrentStep(6)
	u.ui.UpdateDetail("압축해제된 파일들을 검증하고 있습니다...")
	return u.hashManager.VerifyHashSum()
}

func (u *Updater) restartApplication() error {
	u.ui.SetCurrentStep(7)
	u.ui.UpdateDetail("애플리케이션을 재시작하고 있습니다...")
	// 애플리케이션 재시작 로직 구현
	return nil
}

func (u *Updater) restoreFiles() error {
	u.ui.UpdateDetail("파일을 복원하고 있습니다...")
	return u.fileManager.Restore()
}

func (u *Updater) handleError(message string, err error) error {
	u.logger.Error("%s: %v", message, err)
	if restoreErr := u.restoreFiles(); restoreErr != nil {
		u.logger.Error("파일 복원 실패: %v", restoreErr)
		return fmt.Errorf("%s, 복원 실패: %v", message, err)
	}
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
}
