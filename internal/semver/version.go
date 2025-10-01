package semver

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrUnabledToParseVersion = errors.New("unable to parse version")
)

type Version struct {
	Major    int      `json:"major" yaml:"major"`
	Minor    int      `json:"minor" yaml:"minor"`
	Patch    int      `json:"patch" yaml:"patch"`
	Codename string   `json:"codename" yaml:"codename"`
	Tags     []string `json:"tags" yaml:"tags"`
}

func NewVersion(major, minor, patch int, codename string, tags []string) *Version {
	return &Version{
		Major:    major,
		Minor:    minor,
		Patch:    patch,
		Codename: codename,
		Tags:     tags,
	}
}

func VersionFromString(version string) (*Version, error) {
	if version == "" {
		return nil, ErrUnabledToParseVersion
	}
	// Split tags from version
	parts := strings.Split(version, " ")
	versionPart := parts[0]

	var tags []string
	// Iterate over tags and add to optional tag array
	for _, part := range parts[1:] {
		if strings.HasPrefix(part, "--") {
			tags = append(tags, strings.TrimPrefix(part, "--"))
		}
	}
	// Split Semantic Version from Codename
	versionComponents := strings.Split(versionPart, "-")

	// Codename is optional, but we must have at least a version
	// We return an error if we don't'
	if len(versionComponents) < 1 {
		return nil, ErrUnabledToParseVersion
	}

	// Split version into major, minor, and patch
	numbers := strings.Split(versionComponents[0], ".")

	// Try to get major, minor, and patch from version string
	major, err := strconv.Atoi(numbers[0])
	if err != nil {
		major = 0
	}

	minor, err := strconv.Atoi(numbers[1])
	if err != nil {
		minor = 0
	}

	patch, err := strconv.Atoi(numbers[2])
	if err != nil {
		patch = 0
	}

	// If major, minor, and patch are all 0, then we can't parse the version'
	if major == 0 && minor == 0 && patch == 0 {
		return nil, ErrUnabledToParseVersion
	}

	// Codename is optional, we must pass a value to the constructor
	codename := ""
	// If we have a codename, then we pass it to the constructor
	if len(versionComponents) > 1 {
		codename = versionComponents[1]
	}

	return NewVersion(major, minor, patch, codename, tags), nil
}

func (v *Version) String() string {

	tagString := ""
	// if we have tags, then we add them to the string
	if len(v.Tags) > 0 {
		for _, tag := range v.Tags {
			tagString += fmt.Sprintf(" --%s", tag)
		}
	}

	// if we have a codename, include it in the string
	if v.Codename == "" {
		return fmt.Sprintf("%d.%d.%d%s", v.Major, v.Minor, v.Patch, tagString)
	} else {
		return fmt.Sprintf("%d.%d.%d-%s%s", v.Major, v.Minor, v.Patch, v.Codename, tagString)
	}
}
