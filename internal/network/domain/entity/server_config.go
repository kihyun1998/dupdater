// Package entity는 네트워크 도메인의 핵심 개념을 정의합니다
package entity

import "fmt"

// ServerConfig는 서버 설정 정보를 담는 엔티티입니다
type ServerConfig struct {
	Profile Profile // 현재 사용중인 서버 프로필
}

// Profile은 서버 프로필 정보를 담는 값 객체입니다
type Profile struct {
	Name string // 서버 프로필명 (예: "server1")
	IP   string // 서버 IP 주소
}

// GetServerIP는 프로필의 서버 IP를 반환합니다
func (sc *ServerConfig) GetServerIP() (string, error) {
	if sc.Profile.IP == "" {
		return "", fmt.Errorf("서버 IP가 설정되지 않았습니다")
	}
	return sc.Profile.IP, nil
}
