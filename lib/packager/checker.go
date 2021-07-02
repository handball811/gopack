package packager

type Checker interface {
	// Should tell or not
	// 0: it should be cancelled
	// 1~: state of this information
	Check(id uint32) uint32
}

type SetChecker interface {
	Check(ids []uint32) int
}
