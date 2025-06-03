package files

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func GetProjectRootDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for {
		candidate := filepath.Join(dir, "main.go")

		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			panic("cannot find project root: main.go not found")
		}
		dir = parent
	}
}

func IsFileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetFileSize(filePath string) uint64 {
	if IsFileExists(filePath) {
		fileInfo, _ := os.Stat(filePath)
		return uint64(fileInfo.Size())
	}
	return 0
}

func GetFileSegment(filePath string) uint16 {
	re := regexp.MustCompile(`_(\d{3})`)
	match := re.FindStringSubmatch(filepath.Base(filePath))
	if match == nil {
		return 0
	}
	numString := match[1]
	parseUint, err := strconv.ParseUint(numString, 10, 16)
	if err != nil {
		return 0
	}
	return uint16(parseUint)
}

func GetFileName(fileSegment uint16) string {
	return fmt.Sprintf("data_%03d.db", fileSegment)
}
