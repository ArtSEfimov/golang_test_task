package config

import (
	"os"
	"strconv"
)

type Config struct {
}

func NewConfig() *Config {
	return &Config{}
}

func (config *Config) GetDBDir() string {
	dbDir := os.Getenv("DATABASE_DIR")
	if dbDir == "" {
		return "database"
	}
	return dbDir
}

func (config *Config) GetMaxFileSegmentSize() uint64 {
	maxFileSegmentSize, parseErr := strconv.ParseUint(os.Getenv("MAX_FILE_SEGMENT_SIZE"), 10, 64)
	if parseErr != nil || maxFileSegmentSize == 0 {
		return 1024
	}
	return maxFileSegmentSize
}

func (config *Config) GetOriginalDBFileName() string {
	dbFileName := os.Getenv("ORIGINAL_DATABASE_FILE_NAME")
	if dbFileName == "" {
		return "data_001.db"
	}
	return dbFileName
}

func (config *Config) GetIndexFileName() string {
	indexFileName := os.Getenv("INDEX_FILE_NAME")
	if indexFileName == "" {
		return "indexes.db"
	}
	return indexFileName
}
