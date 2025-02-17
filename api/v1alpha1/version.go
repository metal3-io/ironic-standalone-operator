package v1alpha1

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Version struct {
	Major, Minor int
}

func (v Version) IsLatest() bool {
	return v.Major == 0
}

func ParseVersion(version string) (Version, error) {
	if version == VersionLatest {
		return Version{}, nil
	}

	versionSplit := strings.SplitN(version, ".", 2)
	if len(versionSplit) != 2 {
		return Version{}, fmt.Errorf("invalid version %s, expected X.Y", version)
	}

	major, err := strconv.Atoi(versionSplit[0])
	if err != nil || major <= 0 {
		return Version{}, fmt.Errorf("invalid major version %s", versionSplit[0])
	}
	minor, err := strconv.Atoi(versionSplit[1])
	if err != nil {
		return Version{}, fmt.Errorf("invalid minor version %s", versionSplit[1])
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
	if SupportedVersions[version] == "" {
		var versions []string
		for ver := range SupportedVersions {
			versions = append(versions, ver)
		}
		slices.Sort(versions)
		return fmt.Errorf("version %s is not supported, supported versions are %s",
			version, strings.Join(versions, ", "))
	}

	return nil
}

func init() {
	for version := range SupportedVersions {
		MustParseVersion(version)
	}
}
