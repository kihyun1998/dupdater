package repository

import (
	"io"

	"github.com/kihyun1998/dupdater/internal/network/domain/entity"
)

// NetworkRepository는 네트워크 작업을 추상화하는 인터페이스입니다
type NetworkRepository interface {
	// LoadServerConfig는 설정 파일에서 서버 설정을 로드합니다
	LoadServerConfig(configPath, serverName string) (*entity.ServerConfig, error)

	// FetchUpdateFileName은 서버로부터 업데이트 파일명을 조회합니다
	FetchUpdateFileName(serverIP string) (*entity.UpdateFile, error)

	// DownloadFile은 서버로부터 파일을 다운로드합니다
	// io.ReadCloser를 반환하여 스트림 처리가 가능하도록 합니다
	DownloadFile(serverIP string, filename string) (io.ReadCloser, error)
}
