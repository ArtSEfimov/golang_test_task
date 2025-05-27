package db

import (
	"bufio"
	"fmt"
	"go_text_task/pkg/files"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func createPath(fileName string) string {
	return files.MakePath(os.Getenv("DATABASE_DIR"), fileName)
}

func updateFile(oldFileName string, manager *Manager) (updateFileErr error) {
	fileSegment := files.GetFileSegment(oldFileName)

	newFileName := "tmp.db"
	newFilePath := createPath(newFileName)
	oldFilePath := createPath(oldFileName)

	defer func() {
		renameErr := os.Rename(newFilePath, oldFilePath)

		if renameErr != nil {
			if updateFileErr != nil {
				updateFileErr = fmt.Errorf("file update error: %v, rename error: %v", updateFileErr, renameErr)
			} else {
				updateFileErr = renameErr
			}
		}
	}()

	files.CreateFile(os.Getenv("DATABASE_DIR"), newFileName)

	dbFile, openErr := os.OpenFile(newFilePath, os.O_APPEND, fs.ModeAppend)
	if openErr != nil {
		panic(openErr)
	}

	defer func(dbFile *os.File) {
		closeErr := dbFile.Close()
		if closeErr != nil {
			panic(closeErr)
		}
	}(dbFile)

	writer := bufio.NewWriter(dbFile)

	for dbIndex, dataLocation := range manager.IndexMap {
		if dataLocation.DBSegment != fileSegment {
			continue
		}

		data, readErr := Read(dbIndex, manager)
		if readErr != nil {
			return readErr
		}

		seek, _ := dbFile.Seek(0, io.SeekEnd)

		indexInstance := DataLocation{
			DBSegment: fileSegment,
			Size:      uint64(len(data)),
			Seek:      uint64(seek),
		}

		_, writeErr := writer.Write(data)
		if writeErr != nil {
			return writeErr
		}
		flushErr := writer.Flush()
		if flushErr != nil {
			return flushErr
		}

		manager.IndexMap[dbIndex] = indexInstance

	}

	removeErr := os.Remove(oldFilePath)
	if removeErr != nil {
		return removeErr
	}

	return nil

}

func getValidFileName(filePath string, dataSizePointer *uint64) (string, error) {

	if !files.IsFileExists(filePath) {
		files.CreateFile(os.Getenv("DATABASE_DIR"), filepath.Base(filePath))
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
