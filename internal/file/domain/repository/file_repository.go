package repository

import "archive/zip"

// FileRepository는 파일 시스템 작업을 추상화합니다
type FileRepository interface {
	Backup() error
	Restore() error
	ExtractZip(zipFile string) error
	DeleteFile(path string) error

	// 내부 메서드들
	backupFile(src, dest string) error
	backupDirectory(src, dest string) error
	restoreFile(src, dest string) error
	restoreDirectory(src, dest string) error
	cleanCurrentDirectory() error
	extractFile(file *zip.File) error
}
