package db

import (
	"fmt"
	"os"
	"path/filepath"
)

const DataBaseDir = "database"
const MaxFileSegmentSize = 512
const OriginalFileName = "data_001.db"
const IndexFileName = "indexes.db"

var ProjectRootDir = findProjectRootDir()

type DataLocation struct {
	DBSegment uint16
	Seek      uint64
	Size      uint64
}

func findProjectRootDir() string {
	for {
		if info, infoErr := os.Stat("main.go"); infoErr == nil && !info.IsDir() {
			dir, wdErr := os.Getwd()
			if wdErr != nil {
				panic(wdErr)
			}
			return dir
		}

		dir, wdErr := os.Getwd()
		if wdErr != nil {
			panic(wdErr)
		}
		parent := filepath.Dir(dir)
		if dir == parent {
			panic(fmt.Errorf("cannot find root directory"))
		}

		if err := os.Chdir(".."); err != nil {
			panic(err)
		}
	}
}
