package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
)

// IsProcessRunning은 지정된 프로세스가 실행 중인지 확인합니다.
func IsProcessRunning(processName string) (bool, error) {
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s", processName))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("프로세스 확인 실패: %w", err)
	}

	return strings.Contains(string(output), processName), nil
}

// WaitForProcessToEnd는 프로세스가 종료될 때까지 대기합니다.
func WaitForProcessToEnd(processName string, timeout time.Duration, onWait func(attempt, maxAttempts int)) error {
	maxAttempts := int(timeout.Seconds())
	attempt := 0

	for {
		isRunning, err := IsProcessRunning(processName)
		if err != nil {
			return err
		}

		if !isRunning {
			return nil
		}

		attempt++
		if attempt >= maxAttempts {
			return fmt.Errorf("타임아웃: %s 프로세스가 %v 동안 종료되지 않음", processName, timeout)
		}

		if onWait != nil {
			onWait(attempt, maxAttempts)
		}

		time.Sleep(time.Second)
	}
}

// IsAdmin은 현재 프로세스가 관리자 권한으로 실행 중인지 확인합니다.
func IsAdmin() bool {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token := windows.Token(0)
	member, err := token.IsMember(sid)
	return err == nil && member
}

// CreateDirectoryIfNotExists는 디렉토리가 없으면 생성합니다.
func CreateDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

// LaunchProcessWithElevation은 프로세스를 관리자 권한으로 실행합니다.
func LaunchProcessWithElevation(execPath string, args string) error {
	verb := windows.StringToUTF16Ptr("runas") // verb도 StringToUTF16Ptr로 변환
	exe := windows.StringToUTF16Ptr(execPath)
	params := windows.StringToUTF16Ptr(args)
	dir := windows.StringToUTF16Ptr(filepath.Dir(execPath))

	err := windows.ShellExecute(0, verb, exe, params, dir, windows.SW_NORMAL)
	if err != nil {
		return fmt.Errorf("프로세스 실행 실패: %w", err)
	}

	return nil
}

// IsProcessRunningByPID는 지정된 PID의 프로세스가 실행 중인지 확인합니다.
func IsProcessRunningByPID(pid int) (bool, error) {
	if pid <= 0 {
		return false, fmt.Errorf("유효하지 않은 PID: %d", pid)
	}

	if runtime.GOOS == "windows" {
		// PROCESS_QUERY_LIMITED_INFORMATION 플래그를 사용하여 프로세스 핸들을 엽니다.
		h, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(pid))
		if err != nil {
			// 에러가 발생하면 프로세스가 존재하지 않거나 접근 권한이 없다는 의미입니다.
			return false, nil
		}
		defer syscall.CloseHandle(h)

		var exitCode uint32
		err = syscall.GetExitCodeProcess(h, &exitCode)
		if err != nil {
			return false, fmt.Errorf("프로세스 종료 코드 확인 실패: %w", err)
		}
		// STILL_ACTIVE는 Windows에서 프로세스가 여전히 실행 중임을 나타내는 값 (259)
		const STILL_ACTIVE = 259
		if exitCode == STILL_ACTIVE {
			return true, nil
		}
		return false, nil
	}

	// Windows 이외의 시스템에서는 os.FindProcess의 결과에 의존 (추가 검증 필요할 수 있음)
	process, err := os.FindProcess(pid)
	if err != nil {
		return false, fmt.Errorf("프로세스 조회 실패: %w", err)
	}
	// Unix 계열에서는 signal 0을 이용한 확인이 가능
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false, nil
	}
	return true, nil
}

// CheckApplicationRunningByPID는 지정된 PID의 애플리케이션이 실행 중인지 확인합니다
func CheckApplicationRunningByPID(pid int) (bool, error) {
	return IsProcessRunningByPID(pid)
}

// CheckApplicationRunning은 지정된 프로세스가 실행 중인지 확인합니다
func CheckApplicationRunning(appName string) (bool, error) {
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s", appName))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("프로세스 확인 중 오류: %v", err)
	}

	return strings.Contains(string(output), appName), nil
}

// WaitForApplicationToClose는 앱이 종료될 때까지 대기합니다
func WaitForApplicationToClose(appName string, onWait func(string)) error {
	const (
		maxAttempts = 30
		waitTime    = time.Second
	)

	for attempt := 0; attempt < maxAttempts; attempt++ {
		isRunning, err := CheckApplicationRunning(appName)
		if err != nil {
			return fmt.Errorf("프로세스 상태 확인 실패: %v", err)
		}

		if !isRunning {
			if onWait != nil {
				onWait("애플리케이션이 성공적으로 종료되었습니다.")
			}
			return nil
		}

		if onWait != nil {
			onWait(fmt.Sprintf("애플리케이션 종료 대기 중... (%d/%d)", attempt+1, maxAttempts))
		}

		time.Sleep(waitTime)
	}

	return fmt.Errorf("타임아웃: %d초 동안 애플리케이션이 종료되지 않았습니다", maxAttempts)
}

// LaunchApplication은 새 버전의 애플리케이션을 실행합니다
func LaunchApplication(execPath string, args []string) error {
	cmd := exec.Command(execPath, args...)
	// Windows 특정 설정 추가
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_NEW_CONSOLE,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("애플리케이션 실행 실패: %v", err)
	}

	return nil
}

// ElevateProcess는 프로세스를 관리자 권한으로 실행합니다
func ElevateProcess(execPath string, args string) error {
	verb := windows.StringToUTF16Ptr("runas")
	exe := windows.StringToUTF16Ptr(execPath)
	params := windows.StringToUTF16Ptr(args)
	dir := windows.StringToUTF16Ptr("")

	err := windows.ShellExecute(0, verb, exe, params, dir, windows.SW_NORMAL)
	if err != nil {
		return fmt.Errorf("관리자 권한으로 실행 실패: %v", err)
	}

	return nil
}
