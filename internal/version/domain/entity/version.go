package entity

import (
	"fmt"
	"regexp"
	"time"
)

// Version은 버전 정보를 나타내는 도메인 엔티티입니다
type Version struct {
	Major       int       // 주 버전
	Minor       int       // 부 버전
	Patch       int       // 패치 버전
	Date        time.Time // 배포 날짜
	FullVersion string    // 전체 버전 문자열
}

// NewVersion은 새로운 Version 엔티티를 생성합니다
func NewVersion(major, minor, patch int, date time.Time, fullVersion string) *Version {
	return &Version{
		Major:       major,
		Minor:       minor,
		Patch:       patch,
		Date:        date,
		FullVersion: fullVersion,
	}
}

// ParseVersion은 문자열을 Version 엔티티로 파싱합니다
// 예: V3.0.0(2024-01-01)
func ParseVersion(version string) (*Version, error) {
	pattern := `^V(\d+)\.(\d+)\.(\d+)\((\d{4}-\d{2}-\d{2})\)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(version)

	if matches == nil {
		return nil, fmt.Errorf("잘못된 버전 형식: %s (예: V3.0.0(2024-01-01))", version)
	}

	// 버전 번호 파싱
	var major, minor, patch int
	fmt.Sscanf(matches[1], "%d", &major)
	fmt.Sscanf(matches[2], "%d", &minor)
	fmt.Sscanf(matches[3], "%d", &patch)

	// 날짜 파싱
	date, err := time.Parse("2006-01-02", matches[4])
	if err != nil {
		return nil, fmt.Errorf("날짜 파싱 실패: %w", err)
	}

	return NewVersion(major, minor, patch, date, version), nil
}

// String은 Version을 문자열로 변환합니다
func (v *Version) String() string {
	return v.FullVersion
}

// Equal은 두 Version이 동일한지 비교합니다
func (v *Version) Equal(other *Version) bool {
	return v.Major == other.Major &&
		v.Minor == other.Minor &&
		v.Patch == other.Patch
}

// IsNewer는 현재 버전이 다른 버전보다 새로운지 확인합니다
func (v *Version) IsNewer(other *Version) bool {
	if v.Major != other.Major {
		return v.Major > other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor > other.Minor
	}
	if v.Patch != other.Patch {
		return v.Patch > other.Patch
	}
	if v.Date != other.Date {
		return v.Date.After(other.Date)
	}
	return false
}

// Validate는 버전 정보가 유효한지 검증합니다
func (v *Version) Validate() error {
	if v.Major < 0 || v.Minor < 0 || v.Patch < 0 {
		return fmt.Errorf("버전 번호는 음수일 수 없습니다")
	}
	if v.Date.IsZero() {
		return fmt.Errorf("날짜 정보가 없습니다")
	}
	if v.FullVersion == "" {
		return fmt.Errorf("전체 버전 문자열이 비어있습니다")
	}
	return nil
}
