// Package usecase는 네트워크 작업의 비즈니스 로직을 구현합니다
package usecase

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kihyun1998/dupdater/internal/network/domain/repository"
)

// NetworkService는 네트워크 작업의 비즈니스 로직을 구현합니다
type NetworkService struct {
	repo   repository.NetworkRepository
	logger Logger
}

// Logger는 로깅을 위한 인터페이스입니다
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}

// NewNetworkService는 새로운 NetworkService 인스턴스를 생성합니다
func NewNetworkService(repo repository.NetworkRepository, logger Logger) *NetworkService {
	return &NetworkService{
		repo:   repo,
		logger: logger,
	}
}

// GetServerIP는 프로필명을 통해 서버 IP를 조회합니다
func (s *NetworkService) GetServerIP(serverName string) (string, error) {
	// 패닉 복구
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("GetServerIP 함수에서 패닉 발생: %v", r)
		}
	}()

	// 홈 디렉토리 경로 가져오기
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("홈 디렉토리 경로 가져오기 실패: %w", err)
	}

	// 설정 파일 경로 설정 및 로드
	configPath := filepath.Join(homeDir, ".testfolder", "config.json")
	config, err := s.repo.LoadServerConfig(configPath)
	if err != nil {
		return "", fmt.Errorf("서버 설정 로드 실패: %w", err)
	}

	serverIP, err := config.GetServerIP()
	if err != nil {
		return "", fmt.Errorf("서버 IP 가져오기 실패: %w", err)
	}

	return serverIP, nil
}

// GetUpdateFileName은 서버로부터 업데이트 파일명을 가져옵니다
func (s *NetworkService) GetUpdateFileName(serverIP string) (string, error) {
	// 패닉 복구
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("GetUpdateFileName 함수에서 패닉 발생: %v", r)
		}
	}()

	updateFile, err := s.repo.FetchUpdateFileName(serverIP)
	if err != nil {
		return "", fmt.Errorf("업데이트 파일명 가져오기 실패: %w", err)
	}

	if updateFile.Filename == "" {
		return "", fmt.Errorf("업데이트 파일명이 비어있습니다")
	}

	return updateFile.Filename, nil
}

// DownloadFile은 업데이트 파일을 다운로드합니다
func (s *NetworkService) DownloadFile(serverIP string, filename string) (io.ReadCloser, error) {
	reader, fileInfo, err := s.repo.DownloadFile(serverIP, filename)
	if err != nil {
		return nil, fmt.Errorf("파일 다운로드 실패: %w", err)
	}

	s.logger.Info("다운로드 시작 - 파일: %s, 크기: %d bytes", filename, fileInfo.ContentLength)
	return reader, nil
}
