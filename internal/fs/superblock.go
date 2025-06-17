package fs

// Total allocated size for superblock = 1 page = 512B

// Superblock actual size = 168 B total
type SuperBlock struct {
	Magic                 [4]byte // 4B
	Version               [2]byte // 2B
	Pagesize              uint32  // 32B
	TotalPages            uint32  // 32B
	InodeTableStartOffset uint32  // 32B
	BitmapStartOffset     uint32  // 32B
	DataStartOffset       uint32  // 32B
}

func NewSuperBlock() *SuperBlock {

	// inode table offset - 512B - 512 * 8 bits (Superblock tages first page then inode)
	inodeTableOffset := 512 * 8

	// bitmap start offset - 512B + 64KB (inodes + super block)
	bitmapStartOffset := (512 * 8) + (64 * 1024 * 8)

	// data start offset - 512B + 64KB + 512B (inodes + superblock of 1 page + bitmap of 1 page)
	// 65KB
	dataStartOffset := 65 * 1024 * 8

	return &SuperBlock{
		Magic:                 [4]byte{'V', 'D', 'S', 'K'},
		Version:               [2]byte{'0', '1'},
		Pagesize:              512,
		TotalPages:            2048,
		InodeTableStartOffset: uint32(inodeTableOffset),
		BitmapStartOffset:     uint32(bitmapStartOffset),
		DataStartOffset:       uint32(dataStartOffset),
	}

}
