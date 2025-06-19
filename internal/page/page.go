package page
type PageType int

// Each page has a size = 512B
// total pages assigned in 1MB = 2048

const (
    SUPERBLOCK_PAGE PageType = iota
    INODE_PAGE
    BITMAP_PAGE
    DATA_PAGE
)

type Page struct {
	pageType PageType
	data [512]byte
}

func NewPage() *Page{
	return &Page{
		pageType: DATA_PAGE,
		data: [512]byte{},
	}
}

func NewPageWithType(pType PageType) *Page {
    return &Page{
        pageType: pType,
        data:     [512]byte{},
    }
}

func ReadPageFromOffset(){

	

}