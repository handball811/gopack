package lockstep

type LockStepOperator interface {
	Add(id uint32, state interface{})
	Remove(id uint32)
	Range(c func(id uint32, base, expected *Element) bool)
	LoadBase(id uint32) (interface{}, bool)
	LoadExpected(id uint32) (interface{}, bool)
	AddActions(actions map[uint32]interface{}) int
	AddActionSlice(actions []ActionSet) int
	Complete(index int)
	Drop(index int)
}

type Element struct {
	Value interface{}
}

type Action struct {
	Value interface{}
}
