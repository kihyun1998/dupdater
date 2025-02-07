// Package repository는 파일 시스템 작업을 위한 인터페이스를 정의합니다
package repository

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
	BackupDir  string // 백업 디렉토리 경로
	CurrentDir string // 현재 작업 디렉토리
	Logger     Logger // 로거 인터페이스
}

// Logger는 로깅을 위한 인터페이스입니다
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}
