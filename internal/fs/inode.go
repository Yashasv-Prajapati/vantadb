package fs

import (
	"encoding/binary"
)

/*

TODO : Current code stores MAX_PAGES for each inode, while actually we should be only storing the need pages and not max_pages
change the implementation to use only the number of pages need

*/

// Each inode struct will always be less than 64B, so we will pad it to 64B(each inode)
// Thus, since inode table size = 64KB = 65536 B, we get 65536/64 = 1024 unique inodes in the table
const (
	MAX_PAGES = 6
)

type Inode struct {
	Key           [32]byte
	Size          [4]byte // size of that value corresponding to this key in bytes - that number can be represented in 4 bytes because it won't be bigger than 2^32
	NumberofPages [1]byte // can be max 6
	InUse         [1]byte
	PageNumbers   [6]uint32
}

func NewInode(key [32]byte, fileSize [4]byte) *Inode {
	return &Inode{
		Key:           key,
		Size:          fileSize,
		NumberofPages: [1]byte{0},
		InUse:         [1]byte{1},
		PageNumbers:   [6]uint32{},
	}
}

func (i *Inode) ToBytes() []byte {
	byteData := make([]byte, 64)

	copy(byteData[:32], i.Key[:])
	copy(byteData[32:36], i.Size[:])
	copy(byteData[36:37], i.NumberofPages[:])
	copy(byteData[37:38], i.InUse[:])

	// Store up to 6 page numbers (6 * 4 = 24 bytes, total = 61 bytes, fits in 64)
	for j := 0; j < MAX_PAGES && j < len(i.PageNumbers); j++ {
		offset := 38 + (j * 4)
		binary.LittleEndian.PutUint32(byteData[offset:offset+4], i.PageNumbers[j])
	}

	return byteData
}

func FromBytes(data []byte) *Inode {
	chunk := data[0:64]
	var key [32]byte
	var size [4]byte
	var numpages [1]byte
	var inuse [1]byte
	var pageNumbers [MAX_PAGES]uint32

	copy(key[:], chunk[:32])
	copy(size[:], chunk[32:36])
	copy(numpages[:], chunk[36:37])
	copy(inuse[:], chunk[37:38])

	// Read page numbers
	for i := 0; i < MAX_PAGES; i++ {
		offset := 38 + (i * 4)
		pageNumbers[i] = binary.LittleEndian.Uint32(chunk[offset : offset+4])
	}

	return &Inode{
		Key:           key,
		Size:          size,
		NumberofPages: numpages,
		InUse:         inuse,
		PageNumbers:   pageNumbers,
	}
}

func ReadInodes(dataBytes []byte, superblock *SuperBlock) []*Inode {
	var inodes []*Inode
	inodeTableOffset := superblock.InodeTableStartOffset

	// Read 64 bytes chunk by chunk, because each inode struct is padded to 64 bytes
	for i := 0; i < INODE_TABLE_SIZE; i += 64 {
		chunk := dataBytes[inodeTableOffset+uint32(i) : inodeTableOffset+uint32(i)+64]

		// Use FromBytes to properly deserialize the inode
		inode := FromBytes(chunk)
		inodes = append(inodes, inode)
	}

	return inodes
}
