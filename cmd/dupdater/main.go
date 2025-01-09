package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kihyun1998/dupdater/internal/app"
	"github.com/kihyun1998/dupdater/internal/file"
	"github.com/kihyun1998/dupdater/internal/hash"
	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/network"
	"github.com/kihyun1998/dupdater/internal/ui"
)

const (
	AppName    = "dupdater"
	TotalSteps = 8 // 총 업데이트 단계 수
)

var (
	fromVersion = flag.String("fromVersion", "", "현재 앱 버전")
	serverName  = flag.String("server", "server1", "서버 프로필 이름")
)

func main() {
	// 1. 커맨드라인 플래그 파싱
	flag.Parse()

	// 필수 인자 체크
	if *fromVersion == "" {
		fmt.Println("Error: fromVersion is required")
		flag.Usage()
		os.Exit(1)
	}

	// 2. 로거 초기화
	logPath := getLogPath()
	logger, err := logger.New(logger.Config{
		LogPath:    logPath,
		LogLevel:   logger.INFO,
		MaxSize:    10 * 1024 * 1024, // 10MB
		MaxBackups: 5,
	})
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Info("Starting updater - version: %s, server: %s", *fromVersion, *serverName)

	// 3. UI 매니저 초기화
	uiManager := ui.New(ui.Config{
		AppName:    AppName,
		TotalSteps: TotalSteps,
		Logger:     logger,
	})

	// 4. 네트워크 매니저 초기화
	networkManager := network.New(network.Config{
		Logger: logger,
	})

	// 5. 파일 매니저 초기화
	fileManager, err := file.New(file.Config{
		Logger:     logger,
		BackupDir:  filepath.Join(os.TempDir(), "ACRABACK"),
		CurrentDir: getCurrentDir(),
	})
	if err != nil {
		logger.Error("Failed to initialize file manager: %v", err)
		uiManager.ShowError(fmt.Errorf("failed to initialize file manager: %w", err))
		os.Exit(1)
	}

	// 6. 해시 매니저 초기화
	hashManager, err := hash.New(hash.Config{
		Logger:     logger,
		CurrentDir: getCurrentDir(),
	})
	if err != nil {
		logger.Error("Failed to initialize hash manager: %v", err)
		uiManager.ShowError(fmt.Errorf("failed to initialize hash manager: %w", err))
		os.Exit(1)
	}

	// 7. Updater 생성 및 시작
	updater := app.New(app.Config{
		AppName:         AppName,
		FromVersion:     *fromVersion,
		ServerName:      *serverName,
		BackupCompleted: false,
		UIManager:       uiManager,
		Logger:          logger,
		NetworkManager:  networkManager,
		FileManager:     fileManager,
		HashManager:     hashManager,
	})

	// 8. 업데이트 프로세스 시작
	updater.Start()
}

// getLogPath는 로그 파일의 경로를 반환합니다.
func getLogPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Failed to get home directory: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(homeDir, ".testfolder", "logs", "updater.log")
}

// getCurrentDir는 현재 실행 파일의 디렉토리 경로를 반환합니다.
func getCurrentDir() string {
	dir, err := os.Executable()
	if err != nil {
		fmt.Printf("Failed to get executable path: %v\n", err)
		os.Exit(1)
	}
	// dupdater 폴더가 아닌 상위 디렉토리를 반환
	return filepath.Dir(filepath.Dir(dir))
}
