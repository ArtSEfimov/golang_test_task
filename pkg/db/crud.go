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

	manager.Tasks <- func() {
		manager.ID++
		manager.IndexMap[manager.ID] = indexInstance
		manager.DL.Append(manager.ID)
	}

	manager.Tasks <- manager.storeIndexes

	return nil
}

func Read(index uint64, manager *Manager) ([]byte, error) {
	manager.mtx.RLock()
	dbSegment := manager.IndexMap[index].DBSegment
	seek := manager.IndexMap[index].Seek
	size := manager.IndexMap[index].Size
	manager.mtx.RUnlock()

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
	manager.mtx.RLock()
	indexInstance := manager.IndexMap[index]
	manager.mtx.RUnlock()

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
			defer manager.wg.Done()
			updateErr = manager.updateFile(files.GetFileName(oldDBSegment))
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

		manager.mtx.Lock()
		manager.IndexMap[index] = indexInstance
		manager.mtx.Unlock()

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

	manager.Tasks <- manager.storeIndexes

	return nil
}

func Delete(index uint64, manager *Manager) error {
	manager.mtx.RLock()
	dataLocation, ok := manager.IndexMap[index]
	if !ok {
		manager.mtx.RUnlock()
		return fmt.Errorf("index %d not found", index)
	}
	manager.mtx.RUnlock()

	fileSegment := dataLocation.DBSegment

	manager.mtx.Lock()
	delete(manager.IndexMap, index)
	manager.mtx.Unlock()

	manager.DL.Remove(index)

	updateErr := manager.updateFile(files.GetFileName(fileSegment))
	if updateErr != nil {
		return updateErr
	}

	manager.Tasks <- manager.storeIndexes

	return nil
}
