package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/kihyun1998/dupdater/internal/logger/domain/entity"
	"github.com/kihyun1998/dupdater/internal/logger/domain/repository"
)

// FileLogger는 파일 기반 로그 저장소입니다.
type FileLogger struct {
	mu       sync.RWMutex
	config   *repository.LogConfig
	file     *os.File
	fileSize int64
}

// NewFileLogger는 새로운 FileLogger를 생성합니다.
func NewFileLogger(config *repository.LogConfig) (repository.LogRepository, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("로거 설정 검증 실패: %w", err)
	}

	// 로그 디렉토리 생성
	if err := os.MkdirAll(filepath.Dir(config.LogPath), 0755); err != nil {
		return nil, fmt.Errorf("로그 디렉토리 생성 실패: %w", err)
	}

	// 로그 파일 생성 또는 열기
	file, err := os.OpenFile(
		config.LogPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, fmt.Errorf("로그 파일 열기 실패: %w", err)
	}

	// 현재 파일 크기 확인
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("파일 정보 가져오기 실패: %w", err)
	}

	return &FileLogger{
		config:   config,
		file:     file,
		fileSize: info.Size(),
	}, nil
}

// Write는 로그 엔트리를 파일에 기록합니다.
func (f *FileLogger) Write(entry *entity.LogEntry) error {
	if !entry.IsValid() {
		return fmt.Errorf("유효하지 않은 로그 엔트리")
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	// 로그 문자열 생성
	logString := entry.Format() + "\n"
	logSize := int64(len(logString))

	// 파일 크기 체크 및 순환
	if f.fileSize+logSize > f.config.MaxSize {
		if err := f.rotate(); err != nil {
			return fmt.Errorf("로그 파일 순환 실패: %w", err)
		}
	}

	// 로그 기록
	n, err := f.file.WriteString(logString)
	if err != nil {
		return fmt.Errorf("로그 쓰기 실패: %w", err)
	}

	f.fileSize += int64(n)
	return nil
}

// rotate는 로그 파일을 순환합니다.
func (f *FileLogger) rotate() error {
	// 현재 파일 닫기
	if err := f.file.Close(); err != nil {
		return fmt.Errorf("현재 로그 파일 닫기 실패: %w", err)
	}

	// 백업 파일 순환
	for i := f.config.MaxBackups - 1; i >= 0; i-- {
		oldPath := fmt.Sprintf("%s.%d", f.config.LogPath, i)
		newPath := fmt.Sprintf("%s.%d", f.config.LogPath, i+1)

		if i == 0 {
			oldPath = f.config.LogPath
		}

		// 마지막 백업 파일은 삭제
		if i == f.config.MaxBackups-1 {
			os.Remove(newPath)
			continue
		}

		// 파일 이름 변경
		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				return fmt.Errorf("파일 이름 변경 실패: %w", err)
			}
		}
	}

	// 새 로그 파일 생성
	file, err := os.OpenFile(
		f.config.LogPath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0644,
	)
	if err != nil {
		return fmt.Errorf("새 로그 파일 생성 실패: %w", err)
	}

	f.file = file
	f.fileSize = 0
	return nil
}

// Rotate는 수동으로 로그 파일을 순환합니다.
func (f *FileLogger) Rotate() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.rotate()
}

// Close는 로거를 정리합니다.
func (f *FileLogger) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.file != nil {
		if err := f.file.Sync(); err != nil {
			return fmt.Errorf("파일 동기화 실패: %w", err)
		}
		if err := f.file.Close(); err != nil {
			return fmt.Errorf("파일 닫기 실패: %w", err)
		}
		f.file = nil
	}
	return nil
}

// GetCurrentSize는 현재 로그 파일의 크기를 반환합니다.
func (f *FileLogger) GetCurrentSize() int64 {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.fileSize
}

// GetConfig는 현재 로거의 설정을 반환합니다.
func (f *FileLogger) GetConfig() *repository.LogConfig {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.config
}
