package version

import (
	"fmt"
	"regexp"
	"time"
)

// Version은 버전 정보를 나타내는 구조체입니다
type Version struct {
	Major       int
	Minor       int
	Patch       int
	Date        time.Time
	FullVersion string // 원본 버전 문자열 저장
}

// ParseVersion은 문자열을 Version 구조체로 파싱합니다
// 예: V3.0.0(2024-01-01)
func ParseVersion(version string) (*Version, error) {
	pattern := `^V(\d+)\.(\d+)\.(\d+)\((\d{4}-\d{2}-\d{2})\)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(version)

	if matches == nil {
		return nil, fmt.Errorf("잘못된 버전 형식: %s (예: V3.0.0(2024-01-01))", version)
	}

	// Parse version numbers
	var major, minor, patch int
	fmt.Sscanf(matches[1], "%d", &major)
	fmt.Sscanf(matches[2], "%d", &minor)
	fmt.Sscanf(matches[3], "%d", &patch)

	// Parse date
	date, err := time.Parse("2006-01-02", matches[4])
	if err != nil {
		return nil, fmt.Errorf("날짜 파싱 실패: %s", matches[4])
	}

	return &Version{
		Major:       major,
		Minor:       minor,
		Patch:       patch,
		Date:        date,
		FullVersion: version,
	}, nil
}

// String은 Version을 문자열로 변환합니다
func (v *Version) String() string {
	return v.FullVersion
}
