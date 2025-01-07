// Package app은 업데이터의 핵심 어플리케이션 로직을 포함합니다
package app

import (
	"fmt"
	"net/http"
)

// Updater는 업데이트 프로세스의 전체 흐름을 제어하는 구조체입니다
type Updater struct {
	// 설정 관련 필드
	appName     string // 업데이트할 애플리케이션의 이름
	fromVersion string // 현재 애플리케이션의 버전
	serverName  string // 서버 프로필 이름

	// 의존성 주입을 위한 필드
	ui          UI             // 사용자 인터페이스
	logger      Logger         // 로깅 시스템
	network     NetworkManager // 네트워크 관리자
	fileManager FileManager    // 파일 관리자
	serverIP    string         // 조회된 서버 IP
}

// Config는 새로운 Updater를 생성할 때 필요한 설정을 담는 구조체입니다
type Config struct {
	AppName     string // 업데이트할 애플리케이션의 이름
	FromVersion string // 현재 애플리케이션의 버전
	ServerName  string // 서버 프로필 이름
	UI          UI     // 사용자 인터페이스 구현체
	Logger      Logger // 로거 구현체
}

// New는 새로운 Updater 인스턴스를 생성합니다
func New(config Config) *Updater {
	return &Updater{
		appName:     config.AppName,
		fromVersion: config.FromVersion,
		serverName:  config.ServerName,
		ui:          config.UI,
		logger:      config.Logger,
	}
}

// Start는 업데이트 프로세스를 시작합니다
func (u *Updater) Start() error {
	// 패닉 복구를 위한 defer 함수
	defer func() {
		if r := recover(); r != nil {
			u.logger.Error("업데이트 프로세스 중 패닉 발생: %v", r)
		}
	}()

	// 업데이트 프로세스 초기화
	u.ui.SetCurrentStep(0)
	u.ui.UpdateDetail("업데이트 프로세스를 시작합니다...")

	// update_flow.markdown에 정의된 업데이트 흐름에 따라 진행
	if err := u.checkRunningApp(); err != nil {
		return fmt.Errorf("애플리케이션 실행 상태 확인 실패: %v", err)
	}

	if err := u.getServerIP(); err != nil {
		return fmt.Errorf("서버 IP 가져오기 실패: %v", err)
	}

	if err := u.backupFiles(); err != nil {
		return fmt.Errorf("파일 백업 실패: %v", err)
	}

	// 나머지 프로세스는 추후 구현 예정입니다

	return nil
}

// 필요한 인터페이스 정의

// UI는 사용자 인터페이스와의 상호작용을 위한 인터페이스입니다
type UI interface {
	SetCurrentStep(step int)     // 현재 진행 단계를 설정합니다
	UpdateDetail(message string) // 상세 메시지를 업데이트합니다
	ShowError(err error)         // 에러를 표시합니다
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

// checkRunningApp은 업데이트할 애플리케이션이 실행 중인지 확인합니다
func (u *Updater) checkRunningApp() error {
	u.ui.SetCurrentStep(1)
	u.ui.UpdateDetail("애플리케이션 실행 상태를 확인하고 있습니다...")
	// 구체적인 구현은 추후 추가될 예정입니다
	return nil
}

// getServerIP
func (u *Updater) getServerIP() error {
	u.ui.SetCurrentStep(2)
	u.ui.UpdateDetail("서버 IP를 가져오고 있습니다...")

	serverIP, err := u.network.GetServerIP(u.serverName)
	if err != nil {
		return fmt.Errorf("서버 IP 가져오기 실패: %w", err)
	}

	u.serverIP = serverIP
	u.logger.Info("서버 IP 확인됨: %s", serverIP)
	return nil
}

// backupFiles
func (u *Updater) backupFiles() error {
	u.ui.SetCurrentStep(3)
	u.ui.UpdateDetail("파일을 백업하고 있습니다...")

	if err := u.fileManager.Backup(); err != nil {
		return fmt.Errorf("파일 백업 실패: %w", err)
	}

	u.logger.Info("파일 백업 완료")
	return nil
}

// 새로운 메서드 추가: 업데이트 파일 조회
func (u *Updater) getUpdateFileName() error {
	u.ui.SetCurrentStep(3)
	u.ui.UpdateDetail("업데이트 파일 정보를 조회하고 있습니다...")

	filename, err := u.network.GetUpdateFileName(u.serverIP)
	if err != nil {
		return fmt.Errorf("업데이트 파일 정보 조회 실패: %w", err)
	}

	u.logger.Info("업데이트 파일 확인됨: %s", filename)
	return nil
}

// 업데이트 파일 압축해제
func (u *Updater) extractUpdateFile(zipPath string) error {
	u.ui.SetCurrentStep(5)
	u.ui.UpdateDetail("업데이트 파일을 압축해제하고 있습니다...")

	if err := u.fileManager.ExtractZip(zipPath); err != nil {
		return fmt.Errorf("압축해제 실패: %w", err)
	}

	// 압축 해제 후 ZIP 파일 삭제
	if err := u.fileManager.DeleteFile(zipPath); err != nil {
		u.logger.Error("ZIP 파일 삭제 실패: %v", err)
	}

	return nil
}

// 복원 메서드
func (u *Updater) restoreFiles() error {
	u.ui.UpdateDetail("파일을 복원하고 있습니다...")

	if err := u.fileManager.Restore(); err != nil {
		return fmt.Errorf("파일 복원 실패: %w", err)
	}

	u.logger.Info("파일 복원 완료")
	return nil
}
