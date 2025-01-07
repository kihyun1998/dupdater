package network

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// NetworkManager는 네트워크 통신을 관리하는 구조체
type NetworkManager struct {
	logger Logger
	client *http.Client
}

// Logger는 로깅을 위한 인터페이스
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}

// ServerProfile은 서버 설정 정보를 담는 구조체
type ServerProfile struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}

// Config는 설정 파일의 구조를 정의하는 구조체
type Config struct {
	ProfileList []ServerProfile `json:"profileList"`
}

// DownloadProgress는 다운로드 진행 상황을 추적하는 구조체
type DownloadProgress struct {
	Total      int64
	Downloaded int64
}

// NewNetworkManager는 새로운 NetworkManager 인스턴스를 생성
func NewNetworkManager(logger Logger) *NetworkManager {
	return &NetworkManager{
		logger: logger,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetServerIP는 프로필명을 통해 서버 IP를 조회
func (nm *NetworkManager) GetServerIP(serverName string) (string, error) {
	nm.logger.Info("서버 IP 조회 시작: %s", serverName)

	// 홈 디렉토리 획득
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("홈 디렉토리 획득 실패: %w", err)
	}

	// 설정 파일 경로
	configPath := filepath.Join(homeDir, ".acrapoint", "config.json")

	// 설정 파일 읽기
	file, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("설정 파일 읽기 실패: %w", err)
	}

	// JSON 파싱
	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return "", fmt.Errorf("설정 파일 파싱 실패: %w", err)
	}

	// 프로필 검색
	for _, profile := range config.ProfileList {
		if profile.Name == serverName {
			nm.logger.Info("서버 IP 찾음: %s -> %s", serverName, profile.IP)
			return profile.IP, nil
		}
	}

	return "", fmt.Errorf("프로필을 찾을 수 없음: %s", serverName)
}

// GetUpdateFileName은 서버로부터 업데이트 파일명을 조회
func (nm *NetworkManager) GetUpdateFileName(serverIP string) (string, error) {
	nm.logger.Info("업데이트 파일명 조회 시작: %s", serverIP)

	url := fmt.Sprintf("%s/update/updatefilename", serverIP)

	resp, err := nm.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("서버 요청 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("서버 응답 오류: %d", resp.StatusCode)
	}

	var result struct {
		Filename string `json:"filename"`
		Error    string `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("응답 파싱 실패: %w", err)
	}

	if result.Error != "" {
		return "", fmt.Errorf("서버 오류: %s", result.Error)
	}

	if result.Filename == "" {
		return "", fmt.Errorf("파일명이 비어있음")
	}

	nm.logger.Info("업데이트 파일명 조회 완료: %s", result.Filename)
	return result.Filename, nil
}

// DownloadFile은 파일을 다운로드하고 진행 상황을 보고
func (nm *NetworkManager) DownloadFile(url, filepath string, progress chan<- DownloadProgress) error {
	nm.logger.Info("파일 다운로드 시작: %s -> %s", url, filepath)

	// GET 요청
	resp, err := nm.client.Get(url)
	if err != nil {
		return fmt.Errorf("다운로드 요청 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("서버 응답 오류: %d", resp.StatusCode)
	}

	// 파일 생성
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("파일 생성 실패: %w", err)
	}
	defer out.Close()

	// 진행 상황을 추적하는 래퍼
	progressReader := &ProgressReader{
		Reader:   resp.Body,
		Total:    resp.ContentLength,
		Progress: progress,
	}

	// 파일 복사
	_, err = io.Copy(out, progressReader)
	if err != nil {
		return fmt.Errorf("파일 쓰기 실패: %w", err)
	}

	nm.logger.Info("파일 다운로드 완료: %s", filepath)
	return nil
}

// ProgressReader는 다운로드 진행 상황을 추적하는 Reader 구현체
type ProgressReader struct {
	Reader   io.Reader
	Total    int64
	Current  int64
	Progress chan<- DownloadProgress
}

// Read는 io.Reader 인터페이스 구현
func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	pr.Current += int64(n)

	// 진행 상황 전송
	if pr.Progress != nil {
		pr.Progress <- DownloadProgress{
			Total:      pr.Total,
			Downloaded: pr.Current,
		}
	}

	return
}
