package fs

import "encoding/binary"

// Total allocated size for superblock = 1 page = 512B

// Superblock actual size = 26 B total
type SuperBlock struct {
	Magic                 [4]byte // 4B
	Version               [2]byte // 2B
	Pagesize              uint32  // 32 bits = 4 byte
	TotalPages            uint32  // 32 bits = 4 byte
	InodeTableStartOffset uint32  // 32 bits = 4 byte
	BitmapStartOffset     uint32  // 32 bits = 4 byte
	DataStartOffset       uint32  // 32 bits = 4 byte
}

func NewSuperBlock() *SuperBlock {

	// inode table offset - 512B (Superblock tages first page then inode)
	inodeTableOffset := 512

	// bitmap start offset - 512B + 64KB (inodes + super block)
	bitmapStartOffset := (512) + (64 * 1024)

	// data start offset - 512B + 64KB + 512B (inodes + superblock of 1 page + bitmap of 1 page)
	// 65KB
	dataStartOffset := 65 * 1024

	return &SuperBlock{
		Magic:                 [4]byte{'V', 'D', 'S', 'K'},
		Version:               [2]byte{'0', '1'},
		Pagesize:              uint32(512),
		TotalPages:            uint32(2048),
		InodeTableStartOffset: uint32(inodeTableOffset),
		BitmapStartOffset:     uint32(bitmapStartOffset),
		DataStartOffset:       uint32(dataStartOffset),
	}

}

func ReadSuperblock(blockData []byte) *SuperBlock{

	pagesize := blockData[6:10] // 32 bits = 8 bytes
	totalpages := blockData[10:14] // 32 bits = 8 bytes
	inodeTableStartOffset := blockData[14:18] // 32 bits = 8 bytes
	bitmapStartOffset := blockData[18:22] // 32 bits = 8 bytes
	dataStartOffset := blockData[22:26] // 32 bits = 8 bytes
	
	var magic [4]byte
	var version [2]byte

	copy(magic[:], blockData[:4])
	copy(version[:], blockData[4:6])

	return &SuperBlock{
		Magic: magic,
		Version: version,
		Pagesize:              binary.LittleEndian.Uint32(pagesize[:4]),
        TotalPages:            binary.LittleEndian.Uint32(totalpages[:4]),
        InodeTableStartOffset: binary.LittleEndian.Uint32(inodeTableStartOffset[:4]),
        BitmapStartOffset:     binary.LittleEndian.Uint32(bitmapStartOffset[:4]),
        DataStartOffset:       binary.LittleEndian.Uint32(dataStartOffset[:4]),
    }

}