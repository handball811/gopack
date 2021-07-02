package packager

type PackageOp interface {
	Update(actions []Action) (int, int)
	Complete(index int)
	Drop(index int)
}

type PackageCtrl interface {
	Add(id uint32)
	Remove(id uint32)
}
