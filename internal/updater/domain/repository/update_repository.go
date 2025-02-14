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

// WithProgress는 진행 상태 전달을 위한 콜백 함수 타입입니다
type WithProgress func(step int, message string)
