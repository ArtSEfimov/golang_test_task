package db

import (
	"go_text_task/pkg/files"
	"path/filepath"
)

func createPath(fileName string) string {
	return files.MakePath(DataBaseDir, fileName)
}

func getValidFileName(filePath string, dataSizePointer *uint64) (string, error) {

	if !files.IsFileExists(filePath) {
		files.CreateFile(DataBaseDir, filepath.Base(filePath))
		return filePath, nil
	}

	if files.GetFileSize(filePath)+(*dataSizePointer) > MaxFileSegmentSize {
		segment := files.GetFileSegment(filePath)
		segment++
		fileName := files.GetFileName(segment)
		return getValidFileName(createPath(fileName), dataSizePointer)
	}
	return filePath, nil
}
