package lib

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	major, minor, subminor int
}

func NewVersion(major, minor, subminor int) Version {
	return Version{major: major, minor: minor, subminor: subminor}
}

func NewVersionFromVersionString(v string) (Version, error) {
	if strings.TrimSpace(v) == "" {
		return NewVersion(0, 0, 0), nil
	}

	p := strings.Split(v, ".")
	if len(p) != 3 {
		return Version{}, errors.New("invalid version number, expected [major].[minor].[subminor]")
	}

	for i := range p {
		p[i] = strings.TrimSpace(p[i])
		if p[i] == "" {
			return Version{}, errors.New([]string{"major", "minor", "subminor"}[i] + " version is empty")
		}
	}

	major := 0
	minor := 0
	subminor := 0
	var err error

	if major, err = strconv.Atoi(strings.TrimSpace(p[0])); err != nil {
		return Version{}, fmt.Errorf("major is '%s' instead of a number", p[0])
	} else if minor, err = strconv.Atoi(strings.TrimSpace(p[1])); err != nil {
		return Version{}, fmt.Errorf("minor is '%s' instead of a number", p[1])
	} else if subminor, err = strconv.Atoi(strings.TrimSpace(p[2])); err != nil {
		return Version{}, fmt.Errorf("subminor is '%s' instead of a number", p[2])
	}

	return NewVersion(major, minor, subminor), nil
}

func (v Version) Fmt() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.subminor)
}

func (v Version) Compare(ver Version) int {
	if v.major > ver.major {
		return 1
	} else if v.major < ver.major {
		return -1
	}

	if v.minor > ver.minor {
		return 1
	} else if v.minor < ver.minor {
		return -1
	}

	if v.subminor > ver.subminor {
		return 1
	} else if v.subminor < ver.subminor {
		return -1
	}

	return 0
}
