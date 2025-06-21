package vdsk

// import (
// 	"encoding/binary"
// 	"vantadb/internal/fs"
// )

// // total pages = 2048
// // superblock = 1 page = 512B
// // inode table = 128 pages = 64KB
// // bitmpa = 1 page = 512B
// // data pages = 1918

// const (
// 	PAGE_SIZE        = 512 // bytes
// 	TOTAL_PAGES      = 2048
// 	INODE_TABLE_SIZE = 64 * 1024               // 64KB
// 	TOTAL_DISK_SIZE  = TOTAL_PAGES * PAGE_SIZE // 1MB
// )

// func CreateVDSKStorageData() error {

// 	// init diskstorage object
// 	diskStorage := make([]byte, TOTAL_DISK_SIZE)

// 	// put superblock data in diskstorage
// 	superblock := fs.NewSuperBlock()
// 	superblockData := serializeSuperblock(superblock)
// 	copy(diskStorage[0:PAGE_SIZE], superblockData)

// 	// inode table - 64KB zeroes - already zero due to make of byte array

// 	// bit map
// 	bitmap := fs.NewBitmap()
// 	bitmapData := serializeBitmap(bitmap)
// 	bitmapOffset := (PAGE_SIZE + INODE_TABLE_SIZE)                     // 1 page for super block + 128 pages for inode table
// 	copy(diskStorage[bitmapOffset:bitmapOffset+PAGE_SIZE], bitmapData) // take 1 page for bitmap

// 	// data pages - remaining space, no need to fill anything, already zeor due to make

// 	return writeToDisk(diskStorage)

// }

// func ReadVDSKStorage() ([]byte, error) {
// 	diskStorage := make([]byte, TOTAL_DISK_SIZE)

// 	file, err := fs.OpenVDSK()

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

// func serializeSuperblock(sb *fs.SuperBlock) []byte {
// 	data := make([]byte, PAGE_SIZE) // Full page for superblock

// 	copy(data[0:4], sb.Magic[:])
// 	copy(data[4:6], sb.Version[:])
// 	binary.LittleEndian.PutUint32(data[6:10], sb.Pagesize)
// 	binary.LittleEndian.PutUint32(data[10:14], sb.TotalPages)
// 	binary.LittleEndian.PutUint32(data[14:18], sb.InodeTableStartOffset)
// 	binary.LittleEndian.PutUint32(data[18:22], sb.BitmapStartOffset)
// 	binary.LittleEndian.PutUint32(data[22:26], sb.DataStartOffset)

// 	return data
// }

// func serializeBitmap(bm *fs.Bitmap) []byte {
// 	data := make([]byte, PAGE_SIZE)
// 	copy(data, bm.GetBits()) // You'll need to add this method to Bitmap
// 	return data
// }

// func writeToDisk(data []byte) error {
// 	file, err := fs.OpenVDSK()
// 	if err != nil {
// 		return err
// 	}

// 	defer file.Close()

// 	err = fs.WriteVDSK(file, data)

// 	return err
// }
