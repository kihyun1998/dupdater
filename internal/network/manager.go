package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// Manager는 네트워크 통신을 관리하는 구조체입니다
type Manager struct {
	logger Logger // 로깅을 위한 인터페이스
}

// Config는 Manager 생성에 필요한 설정을 담는 구조체입니다
type Config struct {
	Logger Logger
}

// New는 새로운 Manager 인스턴스를 생성합니다
func New(config Config) *Manager {
	return &Manager{
		logger: config.Logger,
	}
}

// ServerConfig는 서버 설정 정보를 담는 구조체입니다
type ServerConfig struct {
	ProfileList []Profile `json:"profileList"`
}

// Profile은 서버 프로필 정보를 담는 구조체입니다
type Profile struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}

// UpdateFileInfo는 업데이트 파일 정보를 담는 구조체입니다
type UpdateFileInfo struct {
	Filename string `json:"filename"`
}

// GetServerIP는 프로필명을 통해 서버 IP를 조회합니다
func (m *Manager) GetServerIP(serverName string) (string, error) {
	// 패닉 복구
	defer func() {
		if r := recover(); r != nil {
			m.logger.Error("GetServerIP 함수에서 패닉 발생: %v", r)
		}
	}()

	// 홈 디렉토리 경로 가져오기
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("홈 디렉토리 경로 가져오기 실패: %w", err)
	}

	// 설정 파일 경로 설정 및 읽기
	configPath := filepath.Join(homeDir, ".testfolder", "config.json")
	file, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("설정 파일 읽기 실패: %w", err)
	}

	// JSON 파싱
	var config ServerConfig
	if err := json.Unmarshal(file, &config); err != nil {
		return "", fmt.Errorf("설정 파일 파싱 실패: %w", err)
	}

	// 프로필 찾기
	for _, profile := range config.ProfileList {
		if profile.Name == serverName {
			return profile.IP, nil
		}
	}

	return "", fmt.Errorf("서버 프로필을 찾을 수 없음: %s", serverName)
}

// GetUpdateFileName은 서버로부터 업데이트 파일명을 가져옵니다
func (m *Manager) GetUpdateFileName(serverIP string) (string, error) {
	// 패닉 복구
	defer func() {
		if r := recover(); r != nil {
			m.logger.Error("GetUpdateFileName 함수에서 패닉 발생: %v", r)
		}
	}()

	// URL 구성 및 요청
	url := fmt.Sprintf("%s/update/updatefilename", serverIP)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("서버 요청 실패: %w", err)
	}
	defer resp.Body.Close()

	// 응답 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("서버 응답 오류: %d", resp.StatusCode)
	}

	// JSON 응답 파싱
	var result UpdateFileInfo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("응답 데이터 파싱 실패: %w", err)
	}

	if result.Filename == "" {
		return "", fmt.Errorf("업데이트 파일명이 없습니다")
	}

	return result.Filename, nil
}

// DownloadFile은 업데이트 파일을 다운로드합니다
func (m *Manager) DownloadFile(serverIP, filename string) (*http.Response, error) {
	url := fmt.Sprintf("%s/update/file", serverIP)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("파일 다운로드 요청 실패: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("파일 다운로드 응답 오류: %d", resp.StatusCode)
	}

	return resp, nil
}

// Logger는 로깅을 위한 인터페이스입니다
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}
