package file

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// ReadFile reads the file at the given path and returns its contents.
func ReadFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// ShouldIgnorePath returns true if the given path should be ignored based on the ignore list.
func ShouldIgnorePath(path string, ignore string) bool {
	if ignore == "" {
		return false
	}

	ignoreList := strings.Split(ignore, ",")

	for _, item := range ignoreList {
		if filepath.Base(path) == item || strings.HasPrefix(path, item) {
			return true
		}
	}
	return false
}

// ShouldIncludeFile returns true if the given file should be included based on the filter.
func ShouldIncludeFile(path string, filter string) bool {
	if filter == "" {
		return true
	}

	filterList := strings.Split(filter, ",")

	for _, item := range filterList {
		if strings.HasSuffix(path, "."+item) {
			return true
		}
	}

	return false
}
