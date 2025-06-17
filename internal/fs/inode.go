package fs

import fs "vantadb/internal/page"

type Inode struct {
	Key           [32]byte
	Size          [4]byte // size of that value corresponding to this key in bytes - that number can be represented in 4 bytes because it won't be bigger than 2^32
	NumberofPages [1]byte // can be max 7
	PageList      []*fs.Page
}

func NewInode(key [32]byte, fileSize [4]byte) *Inode {
	return &Inode{
		Key:           key,
		Size:          fileSize,
		NumberofPages: [1]byte{0},
		PageList:      []*fs.Page{},
	}
}
