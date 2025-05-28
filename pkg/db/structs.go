package db

type DataLocation struct {
	DBSegment uint16 `json:"db_segment"`
	Seek      uint64 `json:"seek"`
	Size      uint64 `json:"size"`
}
type Storage struct {
	ID       uint64                  `json:"id"`
	IndexMap map[uint64]DataLocation `json:"index_map"`
}
