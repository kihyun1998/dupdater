package entity

import "fmt"

// UpdateConfig는 업데이트에 필요한 기본 설정을 담는 엔티티입니다
type UpdateConfig struct {
	AppName     string // 업데이트할 애플리케이션 이름
	FromVersion string // 현재 버전
	ServerName  string // 서버 프로필 이름
	PID         int    // 접속기 애플리케이션의 PID
}

// NewUpdateConfig는 새로운 UpdateConfig 인스턴스를 생성합니다
func NewUpdateConfig(appName, fromVersion, serverName string, pid int) (*UpdateConfig, error) {
	if appName == "" {
		return nil, fmt.Errorf("앱 이름은 필수입니다")
	}
	if fromVersion == "" {
		return nil, fmt.Errorf("현재 버전은 필수입니다")
	}
	if serverName == "" {
		return nil, fmt.Errorf("서버 이름은 필수입니다")
	}
	if pid == 0 {
		return nil, fmt.Errorf("PID는 필수입니다")
	}

	return &UpdateConfig{
		AppName:     appName,
		FromVersion: fromVersion,
		ServerName:  serverName,
		PID:         pid,
	}, nil
}

// Validate는 설정값의 유효성을 검증합니다
func (c *UpdateConfig) Validate() error {
	if c.AppName == "" {
		return fmt.Errorf("앱 이름이 비어있습니다")
	}
	if c.FromVersion == "" {
		return fmt.Errorf("현재 버전이 비어있습니다")
	}
	if c.ServerName == "" {
		return fmt.Errorf("서버 이름이 비어있습니다")
	}
	if c.PID == 0 {
		return fmt.Errorf("PID가 설정되어 있지 않습니다.")
	}
	return nil
}
