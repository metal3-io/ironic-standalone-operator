package v1alpha1

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const (
	versionLatestString = "latest"
)

type Version struct {
	Major, Minor int
}

func (v Version) Compare(other Version) int {
	if v.IsLatest() {
		if other.IsLatest() {
			return 0
		}
		return 1
	}

	if v.Major != other.Major {
		return cmp.Compare(v.Major, other.Major)
	}

	return cmp.Compare(v.Minor, other.Minor)
}

func (v Version) IsLatest() bool {
	return v.Major == 0
}

func (v Version) String() string {
	if v.IsLatest() {
		return versionLatestString
	}

	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

func ParseVersion(version string) (Version, error) {
	if version == versionLatestString {
		return Version{}, nil
	}

	versionSplit := strings.SplitN(version, ".", 2)
	if len(versionSplit) != 2 {
		return Version{}, fmt.Errorf("invalid version %s, expected MAJOR.MINOR", version)
	}

	major, err := strconv.Atoi(versionSplit[0])
	if err != nil || major <= 0 {
		return Version{}, fmt.Errorf("invalid major version %s in %s", versionSplit[0], version)
	}
	minor, err := strconv.Atoi(versionSplit[1])
	if err != nil {
		return Version{}, fmt.Errorf("invalid minor version %s in %s", versionSplit[1], version)
	}

	return Version{Major: major, Minor: minor}, nil
}

func MustParseVersion(version string) Version {
	v, err := ParseVersion(version)
	if err != nil {
		panic(fmt.Sprintf("must parse version %s: %s", version, err))
	}

	return v
}

func ValidateVersion(version string) error {
	parsed, err := ParseVersion(version)
	if err != nil {
		return err
	}

	if SupportedVersions[parsed] == "" {
		var versions []string
		for ver := range SupportedVersions {
			versions = append(versions, ver.String())
		}
		slices.Sort(versions)
		return fmt.Errorf("version %s is not supported, supported versions are %s",
			version, strings.Join(versions, ", "))
	}

	return nil
}
