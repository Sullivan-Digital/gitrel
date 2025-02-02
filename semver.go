package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Function to validate semver format
func validateSemver(version string) bool {
	re := regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z-]+)?(\+[0-9A-Za-z-]+)?$`)
	return re.MatchString(version)
}

// Function to increment version
func incrementVersion(version, part string) string {
	parts := strings.Split(version, ".")
	major := parts[0]
	minor := parts[1]
	patch := parts[2]

	switch part {
	case "major":
		major = fmt.Sprintf("%d", atoi(major)+1)
		minor = "0"
		patch = "0"
	case "minor":
		minor = fmt.Sprintf("%d", atoi(minor)+1)
		patch = "0"
	case "patch":
		patch = fmt.Sprintf("%d", atoi(patch)+1)
	}

	return fmt.Sprintf("%s.%s.%s", major, minor, patch)
}

// Helper function to convert string to int
func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// Function to compare two semantic versions
func compareSemver(v1, v2 string) bool {
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
