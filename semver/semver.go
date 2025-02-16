package semver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Function to validate semver format
func ValidateSemver(version string) bool {
	re := regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z-.]+)?(\+[0-9A-Za-z-.]+)?$`)
	return re.MatchString(version)
}

// Function to increment version
func IncrementVersion(version, part string) string {
	parts := strings.Split(version, "-")
	baseVersion := parts[0]
	preRelease := ""
	if len(parts) > 1 {
		preRelease = parts[1]
	}

	baseParts := strings.Split(baseVersion, ".")
	major := baseParts[0]
	minor := baseParts[1]
	patch := baseParts[2]

	switch part {
	case "major":
		major = fmt.Sprintf("%d", atoi(major)+1)
		minor = "0"
		patch = "0"
		preRelease = ""
	case "minor":
		minor = fmt.Sprintf("%d", atoi(minor)+1)
		patch = "0"
		preRelease = ""
	case "patch":
		patch = fmt.Sprintf("%d", atoi(patch)+1)
		preRelease = ""
	}

	if preRelease != "" {
		return fmt.Sprintf("%s.%s.%s-%s", major, minor, patch, preRelease)
	}
	return fmt.Sprintf("%s.%s.%s", major, minor, patch)
}

// Helper function to convert string to int
func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// true if v1 is less than v2
func CompareSemver(v1, v2 string) bool {
	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")

	for i := 0; i < len(v1Parts) && i < len(v2Parts); i++ {
		v1Num, _ := strconv.Atoi(v1Parts[i])
		v2Num, _ := strconv.Atoi(v2Parts[i])

		if v1Num != v2Num {
			return v1Num < v2Num
		}
	}

	return len(v1Parts) < len(v2Parts)
}
