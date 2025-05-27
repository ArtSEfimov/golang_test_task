package files

import (
	"os"
	"path/filepath"
)

func MakePath(path ...string) string {

	projectRootDir := GetProjectRootDir()
	path = append([]string{projectRootDir}, path...)

	filePath := filepath.Join(path...)
	dirPath := filepath.Dir(filePath)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		panic(err)
	}
	return filePath
}

func CreateFile(path ...string) {
	filePath := MakePath(path...)
	file, creationErr := os.Create(filePath)
	if creationErr != nil {
		panic(creationErr)
	}
	closeErr := file.Close()
	if closeErr != nil {
		panic(closeErr)
	}
}
