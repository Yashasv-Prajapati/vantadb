package kv

import (
	"encoding/binary"
	"fmt"
	"strings"
	"sync"
	"time"
	"github.com/Yashasv-Prajapati/vantadb/internal/fs"
	"github.com/Yashasv-Prajapati/vantadb/internal/wal"
)

var disk *fs.Disk
var batchMode bool
var batchMutex sync.RWMutex
var lastFlush time.Time

func Init(d *fs.Disk) {
	disk = d
}

// Enable batch mode - operations won't immediately flush to disk
func EnableBatchMode() {
	batchMutex.Lock()
	defer batchMutex.Unlock()
	batchMode = true
}

// Disable batch mode and flush all pending changes
func DisableBatchMode() error {
	batchMutex.Lock()
	defer batchMutex.Unlock()
	batchMode = false
	return FlushToDisk()
}

// Force flush all in-memory changes to disk
func FlushToDisk() error {
	// Write bitmap
	if err := disk.WriteBitmapToDisk(); err != nil {
		return err
	}

	// Write all inodes
	for i, inode := range disk.Inodes {
		if err := disk.WriteInodeToDisk(i, inode); err != nil {
			return err
		}
	}

	lastFlush = time.Now()
	return nil
}

// upserts key-value pair in db - key - max 32B value, max 6 pages = 3072B
func Set(key string, value string) (string, error){
	return setInternal(key, value, true)
}

func Get(key string) (string, error) {
	// first we have to search if this key exists or not
	idx := searchKeyInInodes(key)
	if idx == -1 { // key not found
		return "", fmt.Errorf("key not found")
	}

	// else found the key
	inode := disk.Inodes[idx]
	pageNumbers := inode.PageNumbers
	numPages := int(inode.NumberofPages[0])
	if numPages == 0 {
		return "", nil
	}

	actualSize := binary.LittleEndian.Uint32(inode.Size[:])
	value := make([]byte, actualSize) // stores the value for corresponding key

	offset := 0 // to read PAGE_SIZE chunk from each page
	for i := 0; i < numPages; i++ {
		if pageNumbers[i] == 0 {
			break // no more pages
		}
		pageData, err := disk.ReadPageFromDisk(int(pageNumbers[i]))
		if err != nil {
			return "", fmt.Errorf("could not read page from disk")
		}

		bytesToCopy := fs.PAGE_SIZE
		if offset+bytesToCopy > int(actualSize) {
			bytesToCopy = int(actualSize) - offset
		}

		copy(value[offset:offset+bytesToCopy], pageData[:bytesToCopy])
		offset += bytesToCopy
	}
	return strings.TrimRight(string(value), "\x00"), nil

}

func Del(key string) string {
	return delInternal(key, true)
}


func RecoverFromLogs() string {
	EnableBatchMode()
	defer DisableBatchMode()

	wals := wal.GetAllWALRecords()
	if len(wals) == 0 {
		return "no records in WAL file"
	}
	for i := 0; i < len(wals); i++ {
		record := wals[i]
		fmt.Println("RECOVERING ", string(record.Key[:]), string(record.Value))
		
		key := strings.TrimRight(string(record.Key[:]), "\x00")
		value := string(record.Value)

		switch record.EntryType[0] {
		case wal.SET_FLAT:
			setInternal(key, value, false)
			continue
		case wal.DELETE_FLAG:
			delInternal(key, false)
			continue
		default:
			continue
		}
	}

	return "OK"
}
