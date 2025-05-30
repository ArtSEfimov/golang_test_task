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
	//for {
	//	if info, infoErr := os.Stat("main.go"); infoErr == nil && !info.IsDir() {
	//		dir, wdErr := os.Getwd()
	//		if wdErr != nil {
	//			panic(wdErr)
	//		}
	//		return dir
	//	}
	//
	//	dir, wdErr := os.Getwd()
	//	if wdErr != nil {
	//		panic(wdErr)
	//	}
	//	parent := filepath.Dir(dir)
	//	if dir == parent {
	//		panic(fmt.Errorf("cannot find root directory"))
	//	}
	//
	//	if err := os.Chdir(".."); err != nil {
	//		panic(err)
	//	}
	//	//dir = parent
	//}
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
