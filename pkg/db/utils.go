package db

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func createPath(fileName string) string {
	filePath := filepath.Join(ProjectRootDir, DataBaseDir, fileName)
	dirPath := filepath.Dir(filePath)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		panic(err)
	}
	return filePath
}

func createFile(fileName string) error {
	filePath := createPath(fileName)
	file, creationErr := os.Create(filePath)
	if creationErr != nil {
		return creationErr
	}
	closeErr := file.Close()
	if closeErr != nil {
		return closeErr
	}

	return nil
}

func getFileSegment(filePath string) uint16 {
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

func getFileName(fileSegment uint16) string {
	return fmt.Sprintf("data_%03d.db", fileSegment)
}

func isFileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func getFileSize(filePath string) uint64 {
	if isFileExists(filePath) {
		fileInfo, _ := os.Stat(filePath)
		return uint64(fileInfo.Size())
	}
	return 0
}

func updateFile(oldFileName string, manager *Manager) (updateFileErr error) {
	fileSegment := getFileSegment(oldFileName)

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

	creationError := createFile(newFileName)
	if creationError != nil {
		return creationError
	}

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

	if !isFileExists(filePath) {
		creationErr := createFile(filepath.Base(filePath))
		if creationErr != nil {
			return "", creationErr
		}
		return filePath, nil
	}

	if getFileSize(filePath)+(*dataSizePointer) > MaxFileSegmentSize {
		segment := getFileSegment(filePath)
		segment++
		fileName := getFileName(segment)
		return getValidFileName(createPath(fileName), dataSizePointer)
	}
	return filePath, nil
}
