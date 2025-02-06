package entity

// UpdateFile은 업데이트 파일 정보를 담는 엔티티입니다
type UpdateFile struct {
	Filename string // 업데이트 파일명
}

// FileInfo는 다운로드된 파일의 정보를 담는 값 객체입니다
type FileInfo struct {
	ContentLength int64  // 파일 크기
	ContentType   string // 파일 타입
}
