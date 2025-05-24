package db

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"regexp"
	"strconv"
)

func createFile(fileName string) error {
	file, creationErr := os.Create(fileName)
	if creationErr != nil {
		return creationErr
	}
	closeErr := file.Close()
	if closeErr != nil {
		return closeErr
	}

	return nil

}

func getFileSegment(fileName string) uint16 {
	re := regexp.MustCompile(`_(\d{3})`)
	match := re.FindStringSubmatch(fileName)
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

func isFileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func getFileSize(fileName string) uint64 {
	if isFileExists(fileName) {
		fileInfo, _ := os.Stat(fileName)
		return uint64(fileInfo.Size())
	}
	return 0
}

func updateFile(oldFileName string, manager *Manager) (updateFileErr error) {
	fileSegment := getFileSegment(oldFileName)

	newFileName := "tmp.db"
	defer func() {
		renameErr := os.Rename(newFileName, oldFileName)

		if renameErr != nil {
			if updateFileErr != nil {
				updateFileErr = fmt.Errorf("file update error: %w, rename error: %v", updateFileErr, renameErr)
			} else {
				updateFileErr = renameErr
			}
		}
	}()
	creationError := createFile(newFileName)
	if creationError != nil {
		return creationError
	}

	dbFile, openErr := os.OpenFile(newFileName, os.O_APPEND, fs.ModeAppend)
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

	removeErr := os.Remove(oldFileName)
	if removeErr != nil {
		return removeErr
	}

	return nil

}

func getValidFileName(fileName string, dataSizePointer *uint64) (string, error) {

	if !isFileExists(fileName) {
		creationErr := createFile(fileName)
		if creationErr != nil {
			return "", creationErr
		}
		return fileName, nil
	}

	if getFileSize(fileName)+(*dataSizePointer) > MaxFileSegmentSize {
		segment := getFileSegment(fileName)
		segment++
		return getValidFileName(getFileName(segment), dataSizePointer)
	}
	return fileName, nil
}
