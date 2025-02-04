package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/kihyun1998/dupdater/internal/logger/domain/entity"
	"github.com/kihyun1998/dupdater/internal/logger/domain/repository"
)

// FileLoggerRepository는 파일 기반 로그 저장소를 구현합니다
type FileLoggerRepository struct {
	mu          sync.Mutex
	file        *os.File
	config      repository.LogConfig
	currentSize int64
}

// NewFileLoggerRepository는 새로운 FileLoggerRepository를 생성합니다
func NewFileLoggerRepository(config repository.LogConfig) (*FileLoggerRepository, error) {
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
		return nil, fmt.Errorf("로그 파일 생성 실패: %w", err)
	}

	// 현재 파일 크기 확인
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("파일 정보 가져오기 실패: %w", err)
	}

	return &FileLoggerRepository{
		file:        file,
		config:      config,
		currentSize: info.Size(),
	}, nil
}

// Log implements repository.LoggerRepository
func (r *FileLoggerRepository) Log(entry *entity.LogEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 로그 순환 체크
	if err := r.checkRotate(); err != nil {
		return fmt.Errorf("로그 순환 실패: %w", err)
	}

	// 로그 메시지 포맷팅 및 쓰기
	message := entry.FormatMessage() + "\n"
	n, err := r.file.WriteString(message)
	if err != nil {
		return fmt.Errorf("로그 쓰기 실패: %w", err)
	}

	r.currentSize += int64(n)
	return nil
}

// Close implements repository.LoggerRepository
func (r *FileLoggerRepository) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.file != nil {
		if err := r.file.Sync(); err != nil {
			return fmt.Errorf("로그 파일 동기화 실패: %w", err)
		}
		if err := r.file.Close(); err != nil {
			return fmt.Errorf("로그 파일 닫기 실패: %w", err)
		}
		r.file = nil
	}
	return nil
}

// Rotate implements repository.LoggerRepository
func (r *FileLoggerRepository) Rotate() error {
	// 현재 파일 닫기
	if err := r.file.Close(); err != nil {
		return fmt.Errorf("현재 로그 파일 닫기 실패: %w", err)
	}

	// 백업 파일 순환
	for i := r.config.MaxBackups - 1; i >= 0; i-- {
		oldPath := fmt.Sprintf("%s.%d", r.config.LogPath, i)
		newPath := fmt.Sprintf("%s.%d", r.config.LogPath, i+1)

		if i == 0 {
			oldPath = r.config.LogPath
		}

		// 마지막 백업 파일 삭제
		if i == r.config.MaxBackups-1 {
			os.Remove(newPath)
			continue
		}

		// 나머지 파일들 이름 변경
		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				return fmt.Errorf("파일 이름 변경 실패: %w", err)
			}
		}
	}

	// 새 로그 파일 생성
	newFile, err := os.OpenFile(
		r.config.LogPath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0644,
	)
	if err != nil {
		return fmt.Errorf("새 로그 파일 생성 실패: %w", err)
	}

	r.file = newFile
	r.currentSize = 0

	return nil
}

// checkRotate는 로그 파일 크기를 확인하고 필요시 순환합니다
func (r *FileLoggerRepository) checkRotate() error {
	if r.currentSize >= r.config.MaxSize {
		return r.Rotate()
	}
	return nil
}
