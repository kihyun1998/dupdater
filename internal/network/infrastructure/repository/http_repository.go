// Package repository는 네트워크 작업의 실제 구현체를 제공합니다
package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/kihyun1998/dupdater/internal/network/domain/entity"
	"github.com/kihyun1998/dupdater/internal/network/domain/repository"
)

// HTTPRepository는 HTTP 기반의 네트워크 작업을 구현합니다
type HTTPRepository struct {
	client *http.Client
	logger Logger
}

// Logger는 로깅을 위한 인터페이스입니다
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}

// NewHTTPRepository는 새로운 HTTPRepository 인스턴스를 생성합니다
func NewHTTPRepository(logger Logger) repository.NetworkRepository {
	return &HTTPRepository{
		client: &http.Client{},
		logger: logger,
	}
}

// LoadServerConfig는 설정 파일에서 서버 설정을 로드합니다
func (r *HTTPRepository) LoadServerConfig(configPath string) (*entity.ServerConfig, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("설정 파일 읽기 실패: %w", err)
	}

	var config struct {
		ProfileList []struct {
			Name string `json:"name"`
			IP   string `json:"ip"`
		} `json:"profileList"`
	}
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("설정 파일 파싱 실패: %w", err)
	}

	for _, profile := range config.ProfileList {
		if profile.Name == "server1" {
			return &entity.ServerConfig{
				Profile: entity.Profile{
					Name: profile.Name,
					IP:   profile.IP,
				},
			}, nil
		}
	}

	return nil, fmt.Errorf("서버 프로필을 찾을 수 없음: server1")
}

// FetchUpdateFileName은 서버로부터 업데이트 파일명을 조회합니다
func (r *HTTPRepository) FetchUpdateFileName(serverIP string) (*entity.UpdateFile, error) {
	url := fmt.Sprintf("%s/update/updatefilename", serverIP)
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("서버 요청 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("서버 응답 오류: %d", resp.StatusCode)
	}

	var result struct {
		Filename string `json:"filename"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("응답 데이터 파싱 실패: %w", err)
	}

	return &entity.UpdateFile{
		Filename: result.Filename,
	}, nil
}

// DownloadFile은 서버로부터 파일을 다운로드합니다
func (r *HTTPRepository) DownloadFile(serverIP string, filename string) (io.ReadCloser, *entity.FileInfo, error) {
	url := fmt.Sprintf("%s/update/file", serverIP)
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("파일 다운로드 요청 실패: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, nil, fmt.Errorf("파일 다운로드 응답 오류: %d", resp.StatusCode)
	}

	fileInfo := &entity.FileInfo{
		ContentLength: resp.ContentLength,
		ContentType:   resp.Header.Get("Content-Type"),
	}

	return resp.Body, fileInfo, nil
}
