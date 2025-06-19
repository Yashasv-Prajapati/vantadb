package kv

import (
	"fmt"
	"strings"
	"vantadb/internal/fs"

)

var disk *fs.Disk

func Init(d *fs.Disk) {
	disk = d
}

func Set(key string, value string) (bool, error) {

	keyBytes := [32]byte{}
	copy(keyBytes[:], key)

	valueBytes := []byte(value)
	valueSize := len(valueBytes)
	pagesNeeded := (valueSize + fs.PAGE_SIZE - 1) / fs.PAGE_SIZE // ceil division

	if pagesNeeded > 6 {
		return false, fmt.Errorf("value too large, can't be accomodated in 6 pages") // can only allot 6 pages
	}

	// first we have to search if this key exists or not
	for i := 0; i < len(disk.Inodes); i++ {
		if strings.TrimRight(string(disk.Inodes[i].Key[:]), "\x00") == key { // key already exists - update
			return updateExistingKey(i, keyBytes, valueBytes, valueSize, pagesNeeded)
		}
	}

	// does not exist, find empty place in array
	for i := 0; i < len(disk.Inodes); i++ {
		if disk.Inodes[i].InUse[0] == 0 { // not in use
			return createNewKey(i, keyBytes, valueBytes, valueSize, pagesNeeded)
		}
	}

	return false, fmt.Errorf("empty space not found to insert key")
}

func Get(key string) (string, error) {
	// first we have to search if this key exists or not
	for i := 0; i < len(disk.Inodes); i++ {
		if strings.TrimRight(string(disk.Inodes[i].Key[:]), "\x00") == key { // key already exists - update
			pageNumbers := disk.Inodes[i].PageNumbers

			value := make([]byte, fs.PAGE_SIZE*len(pageNumbers)) // stores the value for corresponding key
			offset := 0                                          // to read PAGE_SIZE chunk from each page

			for i := 0; i < len(pageNumbers); i++ {
				pageData, err := disk.ReadPageFromDisk(int(pageNumbers[i]))
				if err != nil {
					return "", fmt.Errorf("could not read page from disk")
				}
				copy(value[offset:offset+fs.PAGE_SIZE], pageData[:])
			}
			return strings.TrimRight(string(value), "\x00"), nil
		}
	}
	return "", fmt.Errorf("key-value not found")
}

// func Del(key string) error {

// }

func updateExistingKey(inodeIndex int, keyBytes [32]byte, valueBytes []byte, valueSize int, pagesNeeded int) (bool, error) {
	inode := disk.Inodes[inodeIndex]

	// first we will free the pages from the bitmap
	// basically we will set all those pages we have occupied free in the bitmap and search for new ones
	for i := 0; i < int(inode.NumberofPages[0]); i++ {
		if inode.PageNumbers[i] != 0 {
			disk.Bitmap.FreePage(int(inode.PageNumbers[i]))
		}
	}

	return allocatePagesAndWriteData(inodeIndex, keyBytes, valueBytes, valueSize, pagesNeeded)

}

func createNewKey(inodeIndex int, keyBytes [32]byte, valueBytes []byte, valueSize, pagesNeeded int) (bool, error) {

	disk.Inodes[inodeIndex].InUse[0] = 1
	return allocatePagesAndWriteData(inodeIndex, keyBytes, valueBytes, valueSize, pagesNeeded)

}

func allocatePagesAndWriteData(inodeIndex int, keyBytes [32]byte, valueBytes []byte, valueSize, pagesNeeded int) (bool, error) {
	inode := disk.Inodes[inodeIndex]

	// set inode metadata
	inode.Key = keyBytes
	sizeBytes := [4]byte{} // size of the value it is holding - value corresponding to key
	copy(sizeBytes[:], []byte{byte(valueSize), byte(valueSize >> 8), byte(valueSize >> 16), byte(valueSize >> 24)})
	inode.Size = sizeBytes
	inode.NumberofPages[0] = byte(pagesNeeded)

	// ------- Now, time to allocate pages and write data
	// find free pages
	freePageNumbers := disk.Bitmap.FindFreePages(pagesNeeded)

	if len(freePageNumbers) == 0 {
		return false, fmt.Errorf("no free pages available")
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

	disk.WriteBitmapToDisk()
	disk.WriteInodeToDisk(inodeIndex, inode)

	return true, nil

}
