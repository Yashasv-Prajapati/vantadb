package fs

// size of bitmap - 512B - 1page - 512 * 8 = 4096 bits
type Bitmap struct {
	bits [PAGE_SIZE]byte
}

func NewBitmap() *Bitmap {
	return &Bitmap{
		bits: [512]byte{0},
	}
}

func ReadBitmap(dataBytes []byte, superblock *SuperBlock) *Bitmap {

	offset := superblock.BitmapStartOffset

	bitmapdata := make([]byte, PAGE_SIZE)

	copy(bitmapdata[:], dataBytes[offset:offset+PAGE_SIZE])
	return &Bitmap{
		bits: [512]byte(bitmapdata),
	}

}

func (bm *Bitmap) AllocatePage(position int) {
	// Setbit Function
	/*
		Each byte has 8 bits, so for a position, we first need to get the position of which byte to take and in that byte, which bit to set
	*/
	byteIndex := position / 8
	bitIndex := position % 8

	// set bit as index -> position
	bm.bits[byteIndex] |= (1 << bitIndex)

}

func (bm *Bitmap) FreePage(position int) {
	// Unsetbit function
	/*
		Each byte has 8 bits, so for a position, we first need to get the position of which byte to take and in that byte, which bit to set
	*/
	byteIndex := position / 8
	bitIndex := position % 8

	// set bit as index -> position
	bm.bits[byteIndex] &= ^(1 << bitIndex)

}

func (bm *Bitmap) GetBits() []byte {
	return bm.bits[:]
}

// FindFreePages returns a slice of free page indices. If numberOfPages <= 0, returns all free pages. Else if numberOfPages > free pages, gives error
func (bm *Bitmap) FindFreePages(numberOfPages int) []int {
	freePages := []int{}
	for i := 0; i < len(bm.bits)*8; i++ {
		byteIndex := i / 8
		bitIndex := i % 8
		if (bm.bits[byteIndex] & (1 << bitIndex)) == 0 {
			freePages = append(freePages, i)
			if numberOfPages > 0 && len(freePages) == numberOfPages {
				break
			}
		}
	}

	if len(freePages) < numberOfPages {
		return []int{}
	}

	return freePages
}

func (bm *Bitmap) FindFreePage() int {
	for i := 0; i < len(bm.bits)*8; i++ {
		byteIndex := i / 8
		bitIndex := i % 8
		if (bm.bits[byteIndex] & (1 << bitIndex)) == 0 {
			return i
		}
	}
	return -1 // no free page found
}
