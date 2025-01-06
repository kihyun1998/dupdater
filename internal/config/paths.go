package config

import (
	"os"
	"path/filepath"
)

// Paths는 어플리케이션에서 사용하는 모든 경로 정보를 관리하는 구조체
type Paths struct {
	// 기본 디렉토리
	AppDir    string // 어플리케이션 실행 디렉토리
	ConfigDir string // 설정 파일 디렉토리
	BackupDir string // 백업 디렉토리
	UpdateDir string // 업데이트 파일 다운로드 디렉토리

	// 설정 파일들
	ConfigFile string // 메인 설정 파일 경로
	HashFile   string // 해시 검증 파일 경로
}

// NewPaths는 새로운 Paths 인스턴스를 생성하고 기본 경로를 설정
func NewPaths() *Paths {
	// 사용자의 홈 디렉토리 가져오기
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	// 임시 디렉토리 경로
	tempDir := os.TempDir()

	// 기본 설정 디렉토리 (.testfolder -> .acrapoint로 변경)
	configDir := filepath.Join(homeDir, ".acrapoint")

	p := &Paths{
		AppDir:     ".", // 현재 디렉토리
		ConfigDir:  configDir,
		BackupDir:  filepath.Join(tempDir, "ACRABACK"),
		UpdateDir:  filepath.Join(tempDir, "ACRAUPDATE"),
		ConfigFile: filepath.Join(configDir, "config.json"),
		HashFile:   "hash_sum.txt",
	}

	return p
}

// EnsureDirectories는 필요한 모든 디렉토리가 존재하는지 확인하고 생성
func (p *Paths) EnsureDirectories() error {
	dirs := []string{
		p.ConfigDir,
		p.BackupDir,
		p.UpdateDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

// CleanupTempDirs는 임시 디렉토리들을 정리
func (p *Paths) CleanupTempDirs() error {
	// UpdateDir은 항상 정리
	if err := os.RemoveAll(p.UpdateDir); err != nil {
		return err
	}

	// BackupDir은 조건부로 정리 (나중에 복구가 필요할 수 있음)
	if _, err := os.Stat(p.BackupDir); err == nil {
		if err := os.RemoveAll(p.BackupDir); err != nil {
			return err
		}
	}

	return nil
}
