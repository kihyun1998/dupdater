package components

import "fyne.io/fyne/v2"

var (
	// 체크 아이콘 (완료 상태)
	resourceCompletedIconSvg = &fyne.StaticResource{
		StaticName: "completed.svg",
		StaticContent: []byte(`
<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
    <circle cx="12" cy="12" r="10" stroke="#4ADE80" stroke-width="2"/>
    <path d="M7 13l3 3 7-7" stroke="#4ADE80" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
</svg>
`),
	}

	// 진행중 아이콘 (다운로드/설치 중)
	resourceProgressIconSvg = &fyne.StaticResource{
		StaticName: "progress.svg",
		StaticContent: []byte(`
<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path d="M12 4v4M12 16v4M8 12H4m16 0h-4" stroke="#60A5FA" stroke-width="2" stroke-linecap="round"/>
    <circle cx="12" cy="12" r="10" stroke="#60A5FA" stroke-width="2"/>
</svg>
`),
	}

	// 실패 아이콘 (에러 상태)
	resourceFailedIconSvg = &fyne.StaticResource{
		StaticName: "failed.svg",
		StaticContent: []byte(`
<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
    <circle cx="12" cy="12" r="10" stroke="#F87171" stroke-width="2"/>
    <path d="M8 8l8 8M16 8l-8 8" stroke="#F87171" stroke-width="2" stroke-linecap="round"/>
</svg>
`),
	}

	// 복원 아이콘 (백업 복원 중)
	resourceRestoringIconSvg = &fyne.StaticResource{
		StaticName: "restoring.svg",
		StaticContent: []byte(`
<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
    <circle cx="12" cy="12" r="10" stroke="#FBBF24" stroke-width="2"/>
    <path d="M16 12l-4-4m0 0l-4 4m4-4v8" stroke="#FBBF24" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
</svg>
`),
	}

	// 경고 아이콘 (업데이트 취소됨)
	resourceWarningIconSvg = &fyne.StaticResource{
		StaticName: "warning.svg",
		StaticContent: []byte(`
<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
    <circle cx="12" cy="12" r="10" stroke="#9CA3AF" stroke-width="2"/>
    <path d="M12 8v5m0 3v.01" stroke="#9CA3AF" stroke-width="2" stroke-linecap="round"/>
</svg>
`),
	}
)
