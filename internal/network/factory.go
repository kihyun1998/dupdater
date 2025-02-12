package network

import (
	"net/http"

	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/network/domain/usecase"
	"github.com/kihyun1998/dupdater/internal/network/infrastructure/repository"
)

// Manager는 네트워크 작업을 관리하는 인터페이스입니다
type Manager interface {
	GetServerIP(serverName string) (string, error)
	GetUpdateFileName(serverIP string) (string, error)
	DownloadFile(serverIP, filename string) (*http.Response, error)
}

// Config는 Manager 생성에 필요한 설정을 담는 구조체입니다
type Config struct {
	Logger logger.Logger
}

// networkManager는 Manager 인터페이스를 구현하는 구조체입니다
type networkManager struct {
	service *usecase.NetworkService
}

// New는 새로운 Manager 인스턴스를 생성합니다
func New(config Config) Manager {
	// HTTP Repository 생성
	repo := repository.NewHTTPRepository(config.Logger)

	// Network Service 생성
	service := usecase.NewNetworkService(repo, config.Logger)

	return &networkManager{
		service: service,
	}
}

// GetServerIP는 프로필명을 통해 서버 IP를 조회합니다
func (m *networkManager) GetServerIP(serverName string) (string, error) {
	return m.service.GetServerIP(serverName)
}

// GetUpdateFileName은 서버로부터 업데이트 파일명을 가져옵니다
func (m *networkManager) GetUpdateFileName(serverIP string) (string, error) {
	return m.service.GetUpdateFileName(serverIP)
}

// DownloadFile은 업데이트 파일을 다운로드합니다
func (m *networkManager) DownloadFile(serverIP, filename string) (*http.Response, error) {
	reader, err := m.service.DownloadFile(serverIP, filename)
	if err != nil {
		return nil, err
	}

	// io.ReadCloser를 http.Response로 변환
	return &http.Response{
		Body: reader,
		// 기존 코드와의 호환성을 위해 Response 객체로 감싸서 반환
	}, nil
}
