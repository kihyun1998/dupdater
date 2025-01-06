# 빌드하는 방법

```bash
fyne package -icon icon.png -name updater
```


## 커스텀 아이콘 사용하는 방법

```bash
fyne bundle -o bundled.go icon.png
```

```go
    icon := fyne.NewStaticResource("icon", resourceIconPng.StaticContent)

	ui.spinnerIcon = widget.NewIcon(icon)
	ui.spinnerIcon.Resize(fyne.NewSize(50, 50))
```

위처럼 사용할 수 있다.


## update flow


```mermaid
flowchart TB
    subgraph Main["main.go"]
        Start([시작]) --> CheckArgs{인자 확인}
        CheckArgs -->|실패| ShowError[에러 표시]
        CheckArgs -->|성공| InitLogger[로거 초기화]
        InitLogger --> ShowMainWindow[메인 윈도우 표시]
    end

    subgraph UpdateProcess["update.go"]
        ShowMainWindow --> CheckRunning{앱 실행중?}
        CheckRunning -->|Yes| Wait[대기]
        Wait --> CheckRunning
        CheckRunning -->|No| GetServerIP[서버 IP 가져오기]
        GetServerIP --> Backup[파일 백업]
        Backup --> GetFileName[업데이트 파일명 가져오기]
        GetFileName --> Download[파일 다운로드]
        Download --> Extract[파일 압축해제]
        Extract --> VerifyHash[해시 검증]
        VerifyHash -->|실패| Restore[백업 복원]
        VerifyHash -->|성공| LaunchApp[새 버전 실행]
    end

    subgraph FileManager["file_manager.go"]
        Backup --> |호출| MoveFiles[파일 이동]
        Extract --> |호출| UnzipFile[파일 압축해제]
        Restore --> |호출| RestoreFiles[파일 복원]
    end

    subgraph HashManager["hash_manager.go"]
        VerifyHash --> |호출| VerifyFileHash[파일 해시 검증]
        VerifyFileHash --> VerifyHashSum[해시섬 검증]
    end

    subgraph NetworkManager["network_manager.go"]
        GetServerIP --> |호출| ReadConfig[설정 파일 읽기]
        GetFileName --> |호출| ServerRequest[서버 요청]
    end

    ShowError --> End([종료])
    LaunchApp --> End
    Restore --> End

    style Main fill:#e1f5fe,stroke:#42a5f5
    style UpdateProcess fill:#fff3e0,stroke:#ff9800
    style FileManager fill:#e8f5e9,stroke:#66bb6a
    style HashManager fill:#f3e5f5,stroke:#ab47bc
    style NetworkManager fill:#fce4ec,stroke:#ec407a
```