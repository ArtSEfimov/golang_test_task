package db

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go_text_task/pkg/files"
	"io"
	"io/fs"
	"os"
)

func (m *Manager) storeIndexes() {

	indexFile, creationErr := os.Create(createPath(IndexFileName))
	if creationErr != nil {
		panic(creationErr)
	}
	defer func(indexFile *os.File) {
		closeErr := indexFile.Close()
		if closeErr != nil {
			panic(closeErr)
		}
	}(indexFile)
	writer := bufio.NewWriter(indexFile)
	encodingErr := json.NewEncoder(writer).Encode(m.Storage)
	if encodingErr != nil {
		panic(encodingErr)
	}

	flushErr := writer.Flush()
	if flushErr != nil {
		panic(flushErr)
	}

}

func (m *Manager) updateFile(oldFileName string) (updateFileErr error) {
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

	files.CreateFile(DataBaseDir, newFileName)

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

	m.mtx.Lock()
	for dbIndex, dataLocation := range m.IndexMap {
		if dataLocation.DBSegment != fileSegment {
			continue
		}

		data, readErr := m.Read(dbIndex)
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

		m.IndexMap[dbIndex] = indexInstance

	}
	m.mtx.Unlock()

	removeErr := os.Remove(oldFilePath)
	if removeErr != nil {
		return removeErr
	}

	return nil

}
