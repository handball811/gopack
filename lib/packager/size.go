package packager

/*
取得する変更対象の情報の個数を制限できるPackage
*/

type SizePackage struct {
	*Package
}

type SizeChecker struct {
	size int
}

func (c *SizeChecker) Check(ids []uint32) int {
	size := len(ids)
	if size > c.size {
		size = c.size
	}
	return size
}

func NewSizePackage(
	size int,
	checker Checker,
) *SizePackage {
	return &SizePackage{
		Package: NewPackage(
			size,
			checker,
			&SizeChecker{
				size: size,
			},
		),
	}
}
