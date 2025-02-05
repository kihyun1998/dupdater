package ports

// FilePort는 파일 시스템 작업을 위한 외부 인터페이스입니다
type FilePort interface {
	// Backup은 현재 디렉토리의 파일들을 백업합니다
	Backup() error

	// Restore는 백업된 파일들을 복원합니다
	Restore() error

	// ExtractZip은 ZIP 파일을 압축 해제합니다
	ExtractZip(zipFile string) error

	// DeleteFile은 지정된 파일을 삭제합니다
	DeleteFile(path string) error
}

// FileConfig는 파일 시스템 초기화에 필요한 설정을 정의합니다
type FileConfig struct {
	// Logger는 로깅을 위한 인터페이스입니다
	Logger Logger

	// BackupDir은 백업 디렉토리 경로입니다 (기본값: %TEMP%/ACRABACK)
	BackupDir string

	// CurrentDir은 현재 작업 디렉토리입니다
	CurrentDir string
}

// Logger는 로깅을 위한 인터페이스입니다
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}

// Validate는 설정이 유효한지 검증합니다
func (c *FileConfig) Validate() error {
	if c.CurrentDir == "" {
		return ErrInvalidCurrentDir
	}
	return nil
}

// NewConfig는 새로운 FileConfig를 생성합니다
func NewConfig(logger Logger, backupDir, currentDir string) *FileConfig {
	return &FileConfig{
		Logger:     logger,
		BackupDir:  backupDir,
		CurrentDir: currentDir,
	}
}
