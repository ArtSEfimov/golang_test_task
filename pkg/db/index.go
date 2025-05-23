package db

const MaxFileSegments = 50
const OriginalFileName = "data_001.db"

var N uint64 = 0

type DataLocation struct {
	DBSegment uint16
	Seek      uint64
	Size      uint64
}

var DataIndexes = make(map[uint64]DataLocation)
