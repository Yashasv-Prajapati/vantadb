package fs

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"sync"
)

// total pages = 2048
// superblock = 1 page = 512B
// inode table = 128 pages = 64KB
// bitmpa = 1 page = 512B
// data pages = 1918
const (
	VDSK_PATH        = "/Users/yashasav_p/Developer/go-projects/vantadb/.vdsk"
	PAGE_SIZE        = 512 // bytes
	TOTAL_PAGES      = 2048
	INODE_TABLE_SIZE = 64 * 1024               // 64KB
	TOTAL_DISK_SIZE  = TOTAL_PAGES * PAGE_SIZE // 1MB
)

type Disk struct {
	File       *os.File
	SuperBlock *SuperBlock
	Inodes     []*Inode
	Bitmap     *Bitmap
	Mutex 	*sync.Mutex
}

func Mount(filePath string) (*Disk, error) {
	if len(filePath) == 0{
		filePath = VDSK_PATH
	}
	
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        return nil, err
    }

	// Get file info to check if it's empty (newly created)
    fileInfo, err := file.Stat()
    if err != nil {
        file.Close()
        return nil, err
    }

	// If file is empty or too small, initialize it
    if fileInfo.Size() < TOTAL_DISK_SIZE {
        // Initialize the disk storage
        err = CreateVDSKStorageData(filePath)
        if err != nil {
            file.Close()
            return nil, err
        }
        
        // Reopen the file to read the initialized data
        file.Close()
        file, err = os.OpenFile(filePath, os.O_RDWR, 0644)
        if err != nil {
            return nil, err
        }
    }

	// read the disk storage
	diskStorage := make([]byte, TOTAL_DISK_SIZE)
	_, err = file.Read(diskStorage)
	if err != nil {
		return nil, err
	}

	// Reset file pointer to beginning for future operations
    _, err = file.Seek(0, 0)
    if err != nil {
        file.Close()
        return nil, err
    }

	superblock := ReadSuperblock(diskStorage)
	inodes := ReadInodes(diskStorage, superblock)
	bitmap := ReadBitmap(diskStorage, superblock)

	disk := &Disk{
		File:       file,
		SuperBlock: superblock,
		Inodes:     inodes,
		Bitmap:     bitmap,
		Mutex:      &sync.Mutex{},
	}

	return disk, nil
}

func OpenVDSK(filePath string) (*os.File, error) {
	if len(filePath) == 0 {
		filePath = VDSK_PATH
	}

	file, err := os.Open(filePath)
	if errors.Is(err, os.ErrNotExist) {
		file, err = os.Create(VDSK_PATH)
		if err != nil {
			return nil, fmt.Errorf("file doesn't exist, and failed to create file: %s", err)
		}
	}

	return file, nil

}

func WriteVDSK(file *os.File, data []byte) error {
	_, err := file.Write(data)
	return err
}

func CloseVDSK(file *os.File) {
	file.Close()
}

func CreateVDSKStorageData(filePath string) error {

	// init diskstorage object
	diskStorage := make([]byte, TOTAL_DISK_SIZE)

	// put superblock data in diskstorage
	superblock := NewSuperBlock()
	superblockData := serializeSuperblock(superblock)
	copy(diskStorage[0:PAGE_SIZE], superblockData)

	// inode table - 64KB zeroes - already zero due to make of byte array

	// bit map
	bitmap := NewBitmap()
	bitmapData := serializeBitmap(bitmap)
	bitmapOffset := (PAGE_SIZE + INODE_TABLE_SIZE)                     // 1 page for super block + 128 pages for inode table
	copy(diskStorage[bitmapOffset:bitmapOffset+PAGE_SIZE], bitmapData) // take 1 page for bitmap

	// data pages - remaining space, no need to fill anything, already zeor due to make

	return writeToDisk(diskStorage, filePath)

}

// func ReadVDSKStorage() ([]byte, error) {
// 	diskStorage := make([]byte, TOTAL_DISK_SIZE)

// 	file, err := OpenVDSK()

// 	if err != nil {
// 		return nil, err
// 	}

// 	defer file.Close()

// 	_, err = file.Read(diskStorage)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return diskStorage, nil
// }

/*
Reads the disk file and returns all the data on the VDSK disk
*/
func ReadDisk(file *os.File) ([]byte, error) {

	var data []byte
	_, err := file.Read(data)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %s", err)

	}
	return data, nil
}



func (disk *Disk) WritePageToDisk(pageNumber int, data [PAGE_SIZE]byte) error {

	offset := disk.SuperBlock.DataStartOffset // offset from where the data starts
	offset += uint32(pageNumber) * PAGE_SIZE // add to offset the pageNumber where we want to go, each page is PAGE_SIZE bytes
	
	disk.Mutex.Lock()
	defer disk.Mutex.Unlock()

	_, err := disk.File.WriteAt(data[:], int64(offset))

	return err
}

func (disk *Disk) ReadPageFromDisk(pageNumber int) ([PAGE_SIZE]byte, error) {
    var data [PAGE_SIZE]byte
    offset := disk.SuperBlock.DataStartOffset + uint32(pageNumber) * PAGE_SIZE
    
    disk.Mutex.Lock()
    defer disk.Mutex.Unlock()
    
    _, err := disk.File.ReadAt(data[:], int64(offset))
    
    return data, err
}

func (disk *Disk) WriteInodeToDisk(inodeIndex int, inode *Inode) error {
    offset := disk.SuperBlock.InodeTableStartOffset + uint32(inodeIndex * 64) // Each inode is 64 bytes
    inodeData := inode.ToBytes()
    
    disk.Mutex.Lock()
    defer disk.Mutex.Unlock()
    
    _, err := disk.File.WriteAt(inodeData, int64(offset))
    
    return err
}

func (disk *Disk) WriteBitmapToDisk() error {
    offset := disk.SuperBlock.BitmapStartOffset
    bitmapData := serializeBitmap(disk.Bitmap)
    
    disk.Mutex.Lock()
    defer disk.Mutex.Unlock()
    
    _, err := disk.File.WriteAt(bitmapData, int64(offset))
    
    return err
}


func serializeSuperblock(sb *SuperBlock) []byte {
	data := make([]byte, PAGE_SIZE) // Full page for superblock

	copy(data[0:4], sb.Magic[:])
	copy(data[4:6], sb.Version[:])
	binary.LittleEndian.PutUint32(data[6:10], sb.Pagesize)
	binary.LittleEndian.PutUint32(data[10:14], sb.TotalPages)
	binary.LittleEndian.PutUint32(data[14:18], sb.InodeTableStartOffset)
	binary.LittleEndian.PutUint32(data[18:22], sb.BitmapStartOffset)
	binary.LittleEndian.PutUint32(data[22:26], sb.DataStartOffset)

	return data
}

func serializeBitmap(bm *Bitmap) []byte {
	data := make([]byte, PAGE_SIZE)
	copy(data, bm.GetBits()) // You'll need to add this method to Bitmap
	return data
}

func writeToDisk(data []byte, filePath string) error {
	file, err := OpenVDSK(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	err = WriteVDSK(file, data)

	return err
}

