package ports

import "errors"

var (
	// ErrInvalidCurrentDir는 현재 디렉토리가 유효하지 않을 때 반환되는 에러입니다
	ErrInvalidCurrentDir = errors.New("현재 디렉토리가 유효하지 않습니다")

	// ErrInvalidBackupDir는 백업 디렉토리가 유효하지 않을 때 반환되는 에러입니다
	ErrInvalidBackupDir = errors.New("백업 디렉토리가 유효하지 않습니다")

	// ErrBackupNotFound는 백업을 찾을 수 없을 때 반환되는 에러입니다
	ErrBackupNotFound = errors.New("백업을 찾을 수 없습니다")

	// ErrFileNotFound는 파일을 찾을 수 없을 때 반환되는 에러입니다
	ErrFileNotFound = errors.New("파일을 찾을 수 없습니다")
)
