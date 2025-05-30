package linked_list

import (
	"bufio"
	"encoding/json"
	"go_text_task/pkg/files"
	"os"
)

const LinkedListStorageName = "linked_list_storage"
const LinkedListFileName = "linked_list.db"

var OrderedMap map[uint64]*Node

func createPath(fileName string) string {
	return files.MakePath(os.Getenv("DATABASE_DIR"), LinkedListStorageName, fileName)
}

func storeOrderedMap(dl *DoubleLinkedList) chan struct{} {

	done := make(chan struct{})

	go func() {

		orderMapStorageFile, openErr := os.Create(createPath(LinkedListFileName))
		if openErr != nil {
			panic(openErr)
		}
		defer func(file *os.File) {
			closeErr := file.Close()
			if closeErr != nil {
				panic(closeErr)
			}
		}(orderMapStorageFile)
		writer := bufio.NewWriter(orderMapStorageFile)

		var flatLinkedList []uint64

		flatLinkedList = make([]uint64, dl.GetSize())
		getFlatLinkedList(dl, flatLinkedList)

		encoderErr := json.NewEncoder(writer).Encode(flatLinkedList)
		if encoderErr != nil {
			panic(encoderErr)
		}
		err := writer.Flush()
		if err != nil {
			panic(err)
		}

		done <- struct{}{}

	}()

	return done

}

func recoverOrderedMap(dll *DoubleLinkedList) {

	orderMapStorageFile, openErr := os.Open(createPath(LinkedListFileName))
	if openErr != nil {
		panic(openErr)
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}(orderMapStorageFile)

	fileReader := bufio.NewReader(orderMapStorageFile)

	var flatLinkedList []uint64
	decodeErr := json.NewDecoder(fileReader).Decode(&flatLinkedList)

	if decodeErr != nil {
		panic(decodeErr)
	}

	getLinkedList(flatLinkedList, dll)

}

func getFlatLinkedList(dl *DoubleLinkedList, dst []uint64) {
	for i, node := 0, dl.Head; node != nil; i, node = i+1, node.Next {
		dst[i] = node.Value
	}
}

func getLinkedList(src []uint64, dll *DoubleLinkedList) {
	if len(src) != 0 {
		for _, ID := range src {
			dll.Append(ID)
		}
	}
}
