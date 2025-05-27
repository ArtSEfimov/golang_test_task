package db

const MaxFileSegmentSize = 512
const OriginalDataBaseFileName = "data_001.db"
const IndexFileName = "indexes.db"

type DataLocation struct {
	DBSegment uint16
	Seek      uint64
	Size      uint64
}
