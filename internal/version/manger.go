package version

import "fmt"

// Manager는 버전 관리를 담당하는 구조체입니다
type Manager struct {
	logger      Logger
	fromVersion *Version
	toVersion   *Version
}

// Config는 Manager 생성에 필요한 설정을 담는 구조체입니다
type Config struct {
	Logger      Logger
	FromVersion string
	ToVersion   string
}

// New는 새로운 Manager 인스턴스를 생성합니다
func New(config Config) (*Manager, error) {
	fromVersion, err := ParseVersion(config.FromVersion)
	if err != nil {
		return nil, fmt.Errorf("현재 버전 파싱 실패: %w", err)
	}

	toVersion, err := ParseVersion(config.ToVersion)
	if err != nil {
		return nil, fmt.Errorf("대상 버전 파싱 실패: %w", err)
	}

	return &Manager{
		logger:      config.Logger,
		fromVersion: fromVersion,
		toVersion:   toVersion,
	}, nil
}

// GetFromVersion은 현재 버전을 반환합니다
func (m *Manager) GetFromVersion() string {
	return m.fromVersion.String()
}

// GetToVersion은 대상 버전을 반환합니다
func (m *Manager) GetToVersion() string {
	return m.toVersion.String()
}

// Logger는 로깅을 위한 인터페이스입니다
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}
