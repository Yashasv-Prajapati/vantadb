package kv

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"
	"github.com/Yashasv-Prajapati/vantadb/internal/fs"
	"github.com/Yashasv-Prajapati/vantadb/internal/wal"
)

// ----------------------------------- internal functions to this package -----------------------------------

func autoFlush() {
	if batchMode && time.Since(lastFlush) > 5*time.Second {
		FlushToDisk()
	}
}

func setInternal(key string, value string, writeWAL bool) (string, error) {

	keyBytes := [32]byte{}
	copy(keyBytes[:], key)

	valueBytes := []byte(value)
	valueSize := len(valueBytes)
	pagesNeeded := (valueSize + fs.PAGE_SIZE - 1) / fs.PAGE_SIZE // ceil division

	if pagesNeeded > 6 {
		return "value too large, can't be accomodated in 6 pages", fmt.Errorf("value too large, can't be accomodated in 6 pages") // can only allot 6 pages
	}

	// first we have to search if this key exists or not
	idx := searchKeyInInodes(key)
	if idx >= 0 { // key found

		check, err := updateExistingKey(idx, keyBytes, valueBytes, valueSize, pagesNeeded, key, value, writeWAL)
		if check && (err == nil) {
			autoFlush()
			return "UPDATE OK", nil
		}
	}

	// does not exist, find empty place in array
	for i := 0; i < len(disk.Inodes); i++ {
		if disk.Inodes[i].InUse[0] == 0 { // not in use

			check, err := createNewKey(i, keyBytes, valueBytes, valueSize, pagesNeeded, key, value, writeWAL)
			if check && (err == nil) {
				autoFlush()
				return "SET OK", nil
			}
		}
	}

	return "empty space not found to insert key", fmt.Errorf("empty space not found to insert key")
}

func delInternal(key string, writeWAL bool) string {

	idx := searchKeyInInodes(key) // idx of the inode
	if idx == -1 {                // key not found - does not exist
		return "key not found"
	}

	// everything is fine, first write this command to wal for safety
	if writeWAL{
		wr := wal.NewWALRecord("delete", key, "")
		wr.WriteWALRecordToFile(0)
	}

	// free its pages from bitmap
	inode := disk.Inodes[idx]
	for i := 0; i < len(inode.PageNumbers); i++ {
		disk.Bitmap.FreePage(int(inode.PageNumbers[i]))
	}

	// free the inode space
	disk.Inodes[idx].InUse[0] = 0

	// Flush to disk if not in batch mode
    batchMutex.RLock()
    shouldFlush := !batchMode
    batchMutex.RUnlock()
    
    if shouldFlush {
        disk.WriteBitmapToDisk()
        disk.WriteInodeToDisk(idx, inode)
    } else {
        autoFlush()
    }

	return "OK"
}

func searchKeyInInodes(key string) int {
	for i := 0; i < len(disk.Inodes); i++ {
		inode := disk.Inodes[i]
		if (inode.InUse[0] == 1) && strings.TrimRight(string(inode.Key[:]), "\x00") == key { // the inode is in use and in that inode we have found the key
			return i
		}
	}
	return -1
}

func updateExistingKey(
	inodeIndex int,
	keyBytes [32]byte,
	valueBytes []byte,
	valueSize int,
	pagesNeeded int,
	key string,
	value string,
	writeWAL bool) (bool, error) {
	inode := disk.Inodes[inodeIndex]

	// first we will free the pages from the bitmap
	// basically we will set all those pages we have occupied free in the bitmap and search for new ones
	for i := 0; i < int(inode.NumberofPages[0]); i++ {
		if inode.PageNumbers[i] != 0 {
			disk.Bitmap.FreePage(int(inode.PageNumbers[i]))
		}
	}
	return allocatePagesAndWriteData(inodeIndex, keyBytes, valueBytes, valueSize, pagesNeeded, key, value, writeWAL)
}

func createNewKey(
	inodeIndex int,
	keyBytes [32]byte,
	valueBytes []byte,
	valueSize,
	pagesNeeded int,
	key string,
	value string,
	writeWAL bool) (bool, error) {
	disk.Inodes[inodeIndex].InUse[0] = 1
	return allocatePagesAndWriteData(inodeIndex, keyBytes, valueBytes, valueSize, pagesNeeded, key, value, writeWAL)
}

func allocatePagesAndWriteData(
	inodeIndex int,
	keyBytes [32]byte,
	valueBytes []byte,
	valueSize,
	pagesNeeded int,
	key string,
	value string,
	writeWAL bool) (bool, error) {

	inode := disk.Inodes[inodeIndex]

	// set inode metadata
	inode.Key = keyBytes
	sizeBytes := [4]byte{} // size of the value it is holding - value corresponding to key

	binary.LittleEndian.PutUint32(sizeBytes[:], uint32(valueSize))

	inode.Size = sizeBytes
	inode.NumberofPages[0] = byte(pagesNeeded)

	// ------- Now, time to allocate pages and write data
	// find free pages
	freePageNumbers := disk.Bitmap.FindFreePages(pagesNeeded)

	if len(freePageNumbers) == 0 {
		return false, fmt.Errorf("no free pages available")
	}

	// everything is alright, we can write to WAL then to disk
	// first write to WAL - this function is only used in set, so we can fix that
	// fmt.Println("WRITING TO WAL FILE")
	if writeWAL {
		wr := wal.NewWALRecord("set", key, value)
		wr.WriteWALRecordToFile(0)
	}

	dataOffset := 0
	for i := 0; i < pagesNeeded; i++ {
		// mark this page as allocated in bitmap - pageNumber
		disk.Bitmap.AllocatePage(freePageNumbers[i])

		inode.PageNumbers[i] = uint32(freePageNumbers[i])

		// now fill the pages with data
		pageData := [512]byte{}
		bytesToCopy := fs.PAGE_SIZE
		if dataOffset+bytesToCopy > len(valueBytes) {
			bytesToCopy = len(valueBytes) - dataOffset
		}

		copy(pageData[:fs.PAGE_SIZE], valueBytes[dataOffset:dataOffset+bytesToCopy])
		dataOffset += bytesToCopy
		if err := disk.WritePageToDisk(freePageNumbers[i], pageData); err != nil {
			return false, fmt.Errorf("failed to write to disk: %v", err)
		}
	}

	// flush to disk if not in batch mode
	batchMutex.RLock()
	shouldFlush := !batchMode
	batchMutex.RUnlock()

	if shouldFlush {
		disk.WriteBitmapToDisk()
		disk.WriteInodeToDisk(inodeIndex, inode)
	}


	return true, nil

}
