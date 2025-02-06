### **새로운 리팩토링 계획 (Updated plan.md)**  
(기존 계획에 **패키지별 순차적인 리팩토링 & 통합 테스트 과정**을 반영)

---

# **패키지별 순차적 리팩토링 및 통합 테스트 계획**

## **개요**  
본 프로젝트의 리팩토링은 **한 패키지씩 점진적으로 변경한 후 기존 코드와 통합하여 빌드 및 테스트를 진행하는 방식**으로 진행된다.  
이를 통해 **기능 안정성을 유지하면서 점진적으로 리팩토링을 수행**할 수 있다.

### **진행 방식**
1. 특정 패키지를 리팩토링  
2. 기존 코드와 통합 후 빌드 및 실행 테스트 진행  
3. 오류가 없는지 확인 후 다음 패키지 리팩토링으로 진행  

> **⚠️ 리팩토링 원칙:**  
> - 새로운 개념이 기존 코드와 자연스럽게 연결되도록 **점진적으로 변경**  
> - 기존 코드가 그대로 동작할 수 있도록 **점진적 도입 방식 (Backward Compatibility)** 활용  
> - 새롭게 정의된 개념이 기존 코드에서 오류를 유발하지 않도록 **초기에는 내부적으로만 사용 후 점진적 확장**  

---

## **Sprint 1: Logger 도메인 리팩토링 (2주)**
📌 **목표:**  
- Logger 객체를 생성하는 단일 지점을 `factory.go`로 이동  
- 기존 Logger 인터페이스와 호환성 유지  

📌 **변경 사항:**  
```plaintext
/internal
  /logger
    /domain
      /entity
        log_entry.go       # 로그 엔트리 구조체 정의
        log_level.go       # 로그 레벨 정의 
      /usecase
        logger_service.go  # 로깅 서비스 인터페이스
      /repository
        log_repository.go  # 로그 저장소 인터페이스
      /ports
        logger_port.go     # 현재 Logger 인터페이스와 호환되는 포트 정의
    /infrastructure
      file_logger.go      # 파일 기반 로거 구현
    factory.go           # 로거 생성을 담당하는 팩토리
```

📌 **검증 단계:**  
✅ 리팩토링 완료 후 기존 코드와 통합하여 빌드 및 실행 테스트 진행  
✅ 기존 Logger 기능이 정상적으로 작동하는지 확인  
✅ 이슈 발생 시 수정 후 다시 통합 테스트  

---

## **Sprint 2: 국제화(i18n) 도메인 리팩토링 (2주)**
📌 **목표:**  
- 기존 `LocaleManager`를 `factory.go` 기반으로 변경  

📌 **변경 사항:**  
```plaintext
/internal
  /i18n
    /domain
      /entity
        locale.go          # 로케일 정보 정의
        message.go         # 메시지 정의
      /usecase
        locale_service.go  # 로케일 서비스 인터페이스
      /repository  
        message_repo.go    # 메시지 저장소 인터페이스
      /ports
        locale_port.go     # 현재 LocaleManager와 호환되는 포트
    /infrastructure
      /providers
        en_provider.go     # 영어 메시지 제공자 
        ko_provider.go     # 한글 메시지 제공자
      locale_manager.go    # 새로운 로케일 관리자 구현
    factory.go           # i18n 생성을 담당하는 팩토리
```

📌 **검증 단계:**  
✅ 기존 `i18n` 모듈과 통합하여 다국어 지원이 정상적으로 작동하는지 테스트  
✅ `ko/en` 언어 변경이 정상적으로 동작하는지 확인  
✅ 빌드 후 UI 및 메시지 출력 테스트  

---

## **Sprint 3: 파일 및 네트워크 도메인 리팩토링 (3주)**
📌 **목표:**  
- `manager.go` 기반 구조를 **도메인 기반 구조로 변경**  
- 기존 기능과 호환성을 유지하면서 `BackupInfo` 같은 새로운 개념을 점진적으로 도입  

📌 **변경 사항:**  
```plaintext
/internal
  /file
    /domain
      /entity
        file_info.go        # 파일 정보 관련 도메인 엔티티
      /repository
        file_repository.go  # 파일 저장소 인터페이스
      /service
        file_service.go     # 파일 서비스 인터페이스
    
    /infrastructure
      /repository
        local_file_repository.go  # 로컬 파일 시스템 저장소 구현체
      file_manager.go      # 기존 매니저를 인프라 계층으로 이동
    factory.go      # 파일 매니저 생성 팩토리

  /network
    /domain
      /entity
        server_info.go    # 서버 정보 정의
        download_info.go  # 다운로드 정보 정의
      /usecase
        network_service.go # 네트워크 서비스 인터페이스 
      /repository
        network_repo.go   # 네트워크 저장소 인터페이스
      /ports
        network_port.go   # 현재 NetworkManager와 호환되는 포트
    /infrastructure  
      http_client.go     # HTTP 클라이언트 구현
    factory.go          # 네트워크 클라이언트 생성 팩토리

```

📌 **검증 단계:**  
✅ 기존 `file_manager.go` 및 `network_manager.go`와 통합하여 테스트  
✅ `BackupInfo` 적용이 기존 기능을 깨지 않는지 확인  
✅ 파일 백업/복원 기능 정상 작동 확인  

---

## **Sprint 4: 해시 및 버전 도메인 리팩토링 (2주)**
📌 **목표:**  
- 기존 `hash_manager.go` 및 `version_manager.go`를 도메인 기반으로 재구성  
- `factory.go` 추가  

📌 **변경 사항:**  
```plaintext
/internal
  /hash
    /domain
      /entity
        hash_info.go     # 해시 정보 정의
      /usecase
        hash_service.go  # 해시 서비스 인터페이스
      /repository
        hash_repo.go     # 해시 저장소 인터페이스
      /ports
        hash_port.go     # 현재 HashManager와 호환되는 포트
    /infrastructure
      hash_manager.go    # 새로운 해시 관리자 구현
    factory.go          # 해시 관리자 생성 팩토리

  /version
    /domain
      /entity
        version_info.go   # 버전 정보 정의
      /usecase
        version_service.go # 버전 서비스 인터페이스
      /repository
        version_repo.go   # 버전 저장소 인터페이스
      /ports  
        version_port.go   # 현재 VersionManager와 호환되는 포트
    /infrastructure
      version_manager.go  # 새로운 버전 관리자 구현
    factory.go          # 버전 관리자 생성 팩토리
```

📌 **검증 단계:**  
✅ 기존 코드와 통합하여 버전 검증 및 해시 체크 기능 테스트  
✅ 해시 검증 및 업데이트 프로세스가 정상적으로 동작하는지 확인  

---

## **Sprint 5: UI 및 테마 도메인 리팩토링 (3주)**
📌 **목표:**  
- UI 레이어를 `factory.go` 기반으로 리팩토링  
- `status_card.go`, `progress_bar.go` 같은 UI 컴포넌트 분리  

```plaintext
/internal
  /ui
    /domain
      /entity
        theme_info.go     # 테마 정보 정의
        window_info.go    # 윈도우 정보 정의
      /usecase
        ui_service.go     # UI 서비스 인터페이스
      /repository
        ui_repo.go       # UI 저장소 인터페이스
      /ports
        ui_port.go       # 현재 UI 관련 인터페이스들과 호환되는 포트
    /infrastructure
      /theme
        light_theme.go   # 라이트 테마 구현
        dark_theme.go    # 다크 테마 구현
      /components
        status_card.go   # 상태 카드 컴포넌트
        progress_bar.go  # 진행 바 컴포넌트
      ui_manager.go     # 새로운 UI 관리자 구현
    factory.go         # UI 컴포넌트 생성 팩토리
```

📌 **검증 단계:**  
✅ UI 변경이 기존 업데이트 플로우를 깨지 않는지 확인  

---

## **Sprint 6: 핵심 업데이트 도메인 리팩토링 (2주)**
📌 **목표:**  
- `internal/core` 도메인 추가  
- 업데이트 서비스 `factory.go`로 생성  

```plaintext
/internal
  /core
    /domain
      /entity
        update_info.go     # 업데이트 정보 정의
      /usecase
        update_service.go  # 업데이트 서비스 인터페이스
      /repository
        update_repo.go     # 업데이트 저장소 인터페이스
    /infrastructure
      update_manager.go    # 새로운 업데이트 관리자 구현
    /di
      container.go        # 의존성 주입 컨테이너
    /config
      app_config.go      # 앱 설정
    factory.go          # 업데이트 서비스 생성 팩토리

  /cmd
    /dupdater
      main.go          # 메인 진입점
```

📌 **검증 단계:**  
✅ 새로운 업데이트 매니저가 기존 업데이트 프로세스와 충돌하지 않는지 확인  

---

## **결론**
- 기존 코드와의 **통합 테스트를 포함하는 방식으로 리팩토링**  
- 한 패키지씩 완료 후 **빌드 및 기능 테스트**  
- 새로운 개념은 **점진적으로 적용**  
- **기능 단위 리팩토링 & 점진적 도입 원칙 준수** 🚀