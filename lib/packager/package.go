package packager

import (
	"github.com/handball811/gopack/lib/lockstep"
)

type Package struct {
	PackageOp
	PackageCtrl
	size       int
	lock       lockstep.LockStepOperator
	checker    Checker
	setchecker SetChecker
	actions    []lockstep.ActionSet
	ids        []uint32
}

func NewPackage(
	size int,
	checker Checker,
	setchecker SetChecker,
) *Package {
	return &Package{
		size: size,
		lock: lockstep.NewLockStep(
			func(current *lockstep.Element, action *lockstep.Action) {
				current.Value = action.Value
			},
			func(from interface{}) interface{} {
				return from.(uint32)
			}),
		checker:    checker,
		setchecker: setchecker,
		actions:    make([]lockstep.ActionSet, size),
		ids:        make([]uint32, size),
	}
}

func NewPackageWithLockStep(
	size int,
	lock lockstep.LockStepOperator,
	checker Checker,
	setchecker SetChecker,
) *Package {
	return &Package{
		size:       size,
		lock:       lock,
		checker:    checker,
		setchecker: setchecker,
		actions:    make([]lockstep.ActionSet, size),
		ids:        make([]uint32, size),
	}
}

func (p *Package) Add(id uint32) {
	p.lock.Add(id, uint32(0))
}

func (p *Package) Remove(id uint32) {
	p.lock.Remove(id)
}

// Retrieve Update Actions
// -> actionSize, index
// check
func (p *Package) Update(
	actions []Action,
) (int, int) {
	size := 0
	mp := p.actions
	maxsize := len(actions)
	if maxsize > p.size {
		maxsize = p.size
	}
	// 予測と違うところを抜き出す
	// サイズを超過したら終了する
	p.lock.Range(func(id uint32, base, expected *lockstep.Element) bool {
		// cur := expected.Value.(bool)
		if p.checker.Check(id) == expected.Value {
			return true
		}
		p.ids[size] = id
		size++
		return size < maxsize
	})
	size = p.setchecker.Check(p.ids[:size])
	for i, id := range p.ids[:size] {
		next := p.checker.Check(id)
		mp[i] = lockstep.ActionSet{
			ID: id,
			Action: &lockstep.Action{
				Value: next,
			},
		}
		actions[i] = Action{
			id:    id,
			state: next,
		}
	}
	// その前に回収する対象を厳選してもらう
	index := p.lock.AddActionSlice(mp[:size])
	return size, index
}

func (p *Package) Complete(index int) {
	p.lock.Complete(index)
}

func (p *Package) Drop(index int) {
	p.lock.Drop(index)
}
