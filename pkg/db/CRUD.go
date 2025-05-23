package db

import (
	"bufio"
	"io"
	"io/fs"
	"os"
)

func Create(data []byte) error {

	dataSize := uint64(len(data))

	dbFileName, creationErr := getValidFileName(OriginalFileName, &dataSize)
	if creationErr != nil {
		panic(creationErr)
	}

	dbSegment := getFileSegment(dbFileName)

	dbFile, openErr := os.OpenFile(dbFileName, os.O_APPEND, fs.ModeAppend)
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

	N++
	DataIndexes[N] = indexInstance

	return nil
}

func Read(index uint64) ([]byte, error) {
	dbSegment := DataIndexes[index].DBSegment
	seek := DataIndexes[index].Seek
	size := DataIndexes[index].Size

	data := make([]byte, size)

	dbFile, openErr := os.OpenFile(getFileName(dbSegment), os.O_RDONLY, fs.ModePerm)

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
	indexInstance := DataIndexes[index]
	oldDBSegment := indexInstance.DBSegment

	dataSize := uint64(len(data))
	fileName := OriginalFileName

	var updateErr error

	for {
		dbFileName, creationErr := getValidFileName(fileName, &dataSize)
		if creationErr != nil {
			panic(creationErr)
		}
		newDBSegment := getFileSegment(dbFileName)

		if newDBSegment == oldDBSegment {
			newDBSegment++
			fileName = getFileName(newDBSegment)
			continue
		}

		manager.wg.Add(1)
		go func() {
			manager.mutex.Lock()
			defer manager.mutex.Unlock()
			defer manager.wg.Done()
			updateErr = updateFile(getFileName(oldDBSegment))
		}()

		dbFile, openErr := os.OpenFile(dbFileName, os.O_APPEND, fs.ModeAppend)
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

		DataIndexes[index] = indexInstance

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

	return nil
}

func Delete(index uint64) {
	delete(DataIndexes, index)
}
