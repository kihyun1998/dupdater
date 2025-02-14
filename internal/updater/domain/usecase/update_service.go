package usecase

import (
	"fmt"
	"time"

	"github.com/kihyun1998/dupdater/internal/logger"
	"github.com/kihyun1998/dupdater/internal/ui"
	"github.com/kihyun1998/dupdater/internal/updater/domain/entity"
	"github.com/kihyun1998/dupdater/internal/updater/domain/repository"
)

// UpdateService는 업데이트 프로세스의 비즈니스 로직을 구현합니다
type UpdateService struct {
	repo   repository.UpdateRepository
	ui     ui.Manager
	logger logger.Logger
	status *entity.UpdateStatus
	config *entity.UpdateConfig
}

// NewUpdateService는 새로운 UpdateService 인스턴스를 생성합니다
func NewUpdateService(
	repo repository.UpdateRepository,
	ui ui.Manager,
	logger logger.Logger,
	status *entity.UpdateStatus,
	config *entity.UpdateConfig,
) *UpdateService {
	return &UpdateService{
		repo:   repo,
		ui:     ui,
		logger: logger,
		status: status,
		config: config,
	}
}

// Start는 업데이트 프로세스를 시작합니다
func (s *UpdateService) Start() {
	// UI 복구 핸들러 설정
	s.ui.SetRestoreHandler(func() {
		if err := s.restoreFiles(); err != nil {
			s.logger.Error("파일 복원 실패: %v", err)
			s.ui.ShowError(fmt.Errorf("복원 실패: %w", err))
		}
	})

	// 업데이트 프로세스 시작
	go func() {
		if err := s.processUpdate(); err != nil {
			s.logger.Error("업데이트 실패: %v", err)
			s.handleError("업데이트 실패", err)
		}
	}()

	// UI 실행
	s.ui.Run()
}

// processUpdate는 실제 업데이트 작업을 수행합니다
func (s *UpdateService) processUpdate() error {
	s.logger.Info("업데이트 프로세스 시작")

	// 1. 애플리케이션 실행 상태 확인
	if err := s.repo.CheckRunningApp(); err != nil {
		return s.handleError("애플리케이션 상태 확인 실패", err)
	}

	// 2. 서버 IP 가져오기
	serverIP, err := s.repo.GetServerIP(s.config.ServerName)
	if err != nil {
		return s.handleError("서버 IP 가져오기 실패", err)
	}
	s.status.SetServerIP(serverIP)

	// 3. 파일 백업
	if err := s.repo.BackupFiles(); err != nil {
		return s.handleError("파일 백업 실패", err)
	}
	s.status.MarkBackupCompleted()

	// 4. 업데이트 파일 다운로드
	updateFile, err := s.repo.DownloadUpdateFile(serverIP)
	if err != nil {
		return s.handleError("업데이트 파일 다운로드 실패", err)
	}
	s.status.SetUpdateFile(updateFile)

	// 5. 업데이트 파일 검증
	if err := s.repo.VerifyUpdateFile(updateFile); err != nil {
		return s.handleError("업데이트 파일 검증 실패", err)
	}

	// 6. 파일 압축해제
	if err := s.repo.ExtractUpdateFile(updateFile); err != nil {
		return s.handleError("파일 압축해제 실패", err)
	}

	// 7. 압축해제된 파일들 검증
	if err := s.repo.VerifyExtractedFiles(); err != nil {
		return s.handleError("압축해제된 파일 검증 실패", err)
	}

	// 8. 애플리케이션 재시작
	if err := s.repo.RestartApplication(s.config); err != nil {
		return s.handleError("애플리케이션 재시작 실패", err)
	}

	return nil
}

// restoreFiles는 백업된 파일들을 복원합니다
func (s *UpdateService) restoreFiles() error {
	if s.status.RestoreCompleted {
		s.logger.Info("이미 복원이 완료되었습니다")
		return nil
	}

	s.logger.Info("파일 복원 시작")
	if err := s.repo.RestoreFiles(); err != nil {
		return err
	}

	s.status.MarkRestoreCompleted()
	time.Sleep(2 * time.Second) // UI 메시지 표시를 위한 대기
	s.ui.ShowRestoreComplete()

	return nil
}

// handleError는 에러 상황을 처리합니다
func (s *UpdateService) handleError(message string, err error) error {
	s.logger.Error("%s: %v", message, err)

	if s.status.NeedsRestore() {
		s.ui.ShowRestoring()
		if restoreErr := s.restoreFiles(); restoreErr != nil {
			s.logger.Error("파일 복원 실패: %v", restoreErr)
			s.ui.ShowError(fmt.Errorf("복원 실패: %v", restoreErr))
			return fmt.Errorf("%s 및 복원 실패: %v", message, err)
		}

		// 복원 성공시 기존 버전으로 앱 재시작
		if err := s.repo.RestartAfterRestore(s.config); err != nil {
			s.logger.Error("복원 후 앱 재시작 실패: %v", err)
		}

		return fmt.Errorf("%s, 파일이 복원됨: %v", message, err)
	}

	s.ui.ShowError(fmt.Errorf("%s: %v", message, err))
	return fmt.Errorf("%s: %v", message, err)
}
