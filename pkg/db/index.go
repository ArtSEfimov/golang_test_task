package db

const MaxFileSegmentSize = 1024
const OriginalFileName = "data_001.db"
const IndexFileName = "indexes.db"

type DataLocation struct {
	DBSegment uint16
	Seek      uint64
	Size      uint64
}
