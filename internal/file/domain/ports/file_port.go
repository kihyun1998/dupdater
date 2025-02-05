package ports

// FilePort는 파일 시스템 작업의 외부 인터페이스를 정의합니다
type FilePort interface {
	// Backup은 현재 디렉토리의 파일들을 백업합니다
	Backup() error

	// Restore는 백업된 파일들을 복원합니다
	Restore() error

	// ExtractZip은 ZIP 파일을 현재 디렉토리에 압축해제합니다
	ExtractZip(zipFile string) error

	// DeleteFile은 지정된 파일을 삭제합니다
	DeleteFile(path string) error
}
