package db

import (
	"bufio"
	"fmt"
	"go_text_task/pkg/files"
	"io"
	"io/fs"
	"os"
)

func Create(data []byte, manager *Manager) error {

	dataSize := uint64(len(data))

	filePath := createPath(OriginalDataBaseFileName)

	dbFilePath, creationErr := getValidFileName(filePath, &dataSize)
	if creationErr != nil {
		panic(creationErr)
	}

	dbSegment := files.GetFileSegment(dbFilePath)

	dbFile, openErr := os.OpenFile(dbFilePath, os.O_APPEND, fs.ModeAppend)
	if openErr != nil {
		panic(openErr)
	}

	defer func(dbFile *os.File) {
		closeErr := dbFile.Close()
		if closeErr != nil {
			panic(closeErr)
		}
	}(dbFile)

	seek, _ := dbFile.Seek(0, io.SeekEnd)

	indexInstance := DataLocation{
		DBSegment: dbSegment,
		Size:      uint64(len(data)),
		Seek:      uint64(seek),
	}

	writer := bufio.NewWriter(dbFile)

	_, writeErr := writer.Write(data)
	if writeErr != nil {
		return writeErr
	}

	flushErr := writer.Flush()
	if flushErr != nil {
		return flushErr
	}

	manager.ID++
	manager.IndexMap[manager.ID] = indexInstance
	manager.DL.Append(manager.ID)

	manager.tasks <- manager.storeIndexes

	return nil
}

func Read(index uint64, manager *Manager) ([]byte, error) {
	dbSegment := manager.IndexMap[index].DBSegment
	seek := manager.IndexMap[index].Seek
	size := manager.IndexMap[index].Size

	data := make([]byte, size)

	fileName := files.GetFileName(dbSegment)
	filePath := createPath(fileName)

	dbFile, openErr := os.OpenFile(filePath, os.O_RDONLY, fs.ModePerm)

	if openErr != nil {
		panic(openErr)
	}
	defer func(dbFile *os.File) {
		closeErr := dbFile.Close()
		if closeErr != nil {
			panic(closeErr)
		}
	}(dbFile)

	newSectionReader := io.NewSectionReader(dbFile, int64(seek), int64(size))
	n, readErr := newSectionReader.Read(data)

	if readErr != nil {
		return nil, readErr
	}

	return data[:n], nil

}

func Update(index uint64, data []byte, manager *Manager) error {
	indexInstance := manager.IndexMap[index]
	oldDBSegment := indexInstance.DBSegment

	dataSize := uint64(len(data))

	filePath := createPath(OriginalDataBaseFileName)

	var updateErr error

	for {
		dbFilePath, creationErr := getValidFileName(filePath, &dataSize)
		if creationErr != nil {
			panic(creationErr)
		}
		newDBSegment := files.GetFileSegment(dbFilePath)

		if newDBSegment == oldDBSegment {
			newDBSegment++
			fileName := files.GetFileName(newDBSegment)
			filePath = createPath(fileName)
			continue
		}

		manager.wg.Add(1)
		go func() {
			manager.mtx.Lock()
			defer manager.mtx.Unlock()
			defer manager.wg.Done()
			updateErr = updateFile(files.GetFileName(oldDBSegment), manager)
		}()

		dbFile, openErr := os.OpenFile(dbFilePath, os.O_APPEND, fs.ModeAppend)
		if openErr != nil {
			panic(openErr)
		}

		seek, _ := dbFile.Seek(0, io.SeekEnd)

		indexInstance = DataLocation{
			DBSegment: newDBSegment,
			Size:      uint64(len(data)),
			Seek:      uint64(seek),
		}

		writer := bufio.NewWriter(dbFile)

		_, writeErr := writer.Write(data)
		if writeErr != nil {
			return writeErr
		}

		manager.IndexMap[index] = indexInstance

		flushErr := writer.Flush()
		if flushErr != nil {
			return flushErr
		}

		defer func(dbFile *os.File) {
			closeErr := dbFile.Close()
			if closeErr != nil {
				panic(closeErr)
			}
		}(dbFile)

		break
	}

	manager.wg.Wait()

	if updateErr != nil {
		return updateErr
	}

	manager.tasks <- manager.storeIndexes

	return nil
}

func Delete(index uint64, manager *Manager) error {
	fileSegment := manager.IndexMap[index].DBSegment

	if _, ok := manager.IndexMap[index]; !ok {
		return fmt.Errorf("index %d not found", index)
	}

	delete(manager.IndexMap, index)
	manager.DL.Remove(index)

	manager.tasks <- manager.storeIndexes

	updateErr := updateFile(files.GetFileName(fileSegment), manager)
	if updateErr != nil {
		return updateErr
	}

	return nil
}
