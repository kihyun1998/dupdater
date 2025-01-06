package utils

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
)

// ProcessManager는 프로세스 관리를 위한 구조체
type ProcessManager struct {
	logger Logger
}

// Logger는 로깅을 위한 인터페이스
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
}

// NewProcessManager는 새로운 ProcessManager 인스턴스를 생성
func NewProcessManager(logger Logger) *ProcessManager {
	return &ProcessManager{
		logger: logger,
	}
}

// IsProcessRunning은 지정된 프로세스가 실행 중인지 확인
func (pm *ProcessManager) IsProcessRunning(processName string) (bool, error) {
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s", processName))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("프로세스 목록 획득 실패: %w", err)
	}

	return strings.Contains(string(output), processName), nil
}

// WaitForProcessToEnd는 프로세스가 종료될 때까지 대기
func (pm *ProcessManager) WaitForProcessToEnd(processName string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		running, err := pm.IsProcessRunning(processName)
		if err != nil {
			return fmt.Errorf("프로세스 상태 확인 실패: %w", err)
		}

		if !running {
			return nil
		}

		time.Sleep(time.Second)
	}

	return fmt.Errorf("timeout waiting for process to end: %s", processName)
}

// LaunchAppWithElevation은 관리자 권한으로 애플리케이션 실행
func (pm *ProcessManager) LaunchAppWithElevation(appPath string, args []string) error {
	lpFile := windows.StringToUTF16Ptr(appPath)
	lpParameters := windows.StringToUTF16Ptr(strings.Join(args, " "))
	lpDirectory := windows.StringToUTF16Ptr("")

	err := windows.ShellExecute(
		0,
		windows.StringToUTF16Ptr("runas"),
		lpFile,
		lpParameters,
		lpDirectory,
		windows.SW_NORMAL,
	)

	if err != nil {
		return fmt.Errorf("애플리케이션 실행 실패: %w", err)
	}

	return nil
}

// RunCommand는 명령을 실행하고 결과를 반환
func (pm *ProcessManager) RunCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("명령 실행 실패: %w", err)
	}

	return string(output), nil
}
