package lockstep

import (
	"sort"
	"sync"

	"github.com/handball811/goring"
	"go.uber.org/atomic"
)

var (
	ActionSetSize = 1024
	ActionSetPool = sync.Pool{
		New: func() interface{} {
			return make([]ActionSet, ActionSetSize)
		},
	}
)

type LockStep struct {
	baseState     map[uint32]*Element     // uint32 -> interface{}
	expectedState map[uint32]*Element     // uint32 -> interface{}
	actions       goring.RingOp           // LockStepAction
	actionMap     map[int]*LockStepAction // int -> LockStepAction()

	indexer *atomic.Int64
	step    func(current *Element, action *Action)
	clone   func(from interface{}) interface{}
}

type LockStepAction struct {
	complete bool
	dropped  bool
	index    int
	actions  ActionSlice
}

func NewLockStep(
	step func(current *Element, action *Action), // next
	clone func(from interface{}) interface{},
) *LockStep {
	return &LockStep{
		baseState:     make(map[uint32]*Element),
		expectedState: make(map[uint32]*Element),
		actions:       goring.NewRing(),
		actionMap:     make(map[int]*LockStepAction),
		indexer:       atomic.NewInt64(1),
		step:          step,
		clone:         clone,
	}
}

func (s *LockStep) Add(id uint32, state interface{}) {
	s.baseState[id] = &Element{
		Value: s.clone(state),
	}
	s.expectedState[id] = &Element{
		Value: s.clone(state),
	}
	s.reset(id)
}

func (s *LockStep) Remove(id uint32) {
	delete(s.baseState, id)
	delete(s.expectedState, id)
}

func (s *LockStep) Range(c func(id uint32, base, expected *Element) bool) {
	for id, base := range s.baseState {
		expected, ok := s.expectedState[id]
		if !ok {
			continue
		}
		if !c(id, base, expected) {
			return
		}
	}
}

func (s *LockStep) LoadBase(id uint32) (interface{}, bool) {
	if ret, ok := s.baseState[id]; ok {
		return ret.Value, true
	}
	return nil, false
}

func (s *LockStep) LoadExpected(id uint32) (interface{}, bool) {
	if ret, ok := s.expectedState[id]; ok {
		return ret.Value, true
	}
	return nil, false
}

func (s *LockStep) AddActions(
	actions map[uint32]interface{},
) int {
	// create action
	index := int(s.indexer.Inc())
	action := &LockStepAction{
		complete: false,
		dropped:  false,
		index:    index,
		actions:  ActionSetPool.Get().([]ActionSet),
	}
	if len(actions) > cap(action.actions) {
		action.actions = make([]ActionSet, len(actions))
	}
	action.actions = action.actions[:len(actions)]
	mapCopy(action.actions, actions)

	s.addActions(index, action)
	return index
}

func (s *LockStep) AddActionSlice(
	actions []ActionSet,
) int {
	// create action
	index := int(s.indexer.Inc())
	action := &LockStepAction{
		complete: false,
		dropped:  false,
		index:    index,
		actions:  ActionSetPool.Get().([]ActionSet),
	}
	if len(actions) > cap(action.actions) {
		action.actions = make([]ActionSet, len(actions))
	}
	action.actions = action.actions[:len(actions)]
	mapSliceCopy(action.actions, actions)

	s.addActions(index, action)
	return index
}

func (s *LockStep) addActions(
	index int,
	action *LockStepAction,
) {
	// register action
	s.actions.Push(action)
	s.actionMap[index] = action

	// step up
	for _, action := range action.actions {
		if current, ok := s.expectedState[action.ID]; ok {
			s.step(current, action.Action)
		}
	}
	return
}

func (s *LockStep) Complete(index int) {
	actions, ok := s.actionMap[index]
	if !ok {
		return
	}
	actions.complete = true

	s.concat()

	// remove
	delete(s.actionMap, index)
	return
}

func (s *LockStep) Drop(index int) {
	actions, ok := s.actionMap[index]
	if !ok {
		return
	}
	actions.dropped = true

	for _, action := range actions.actions {
		s.reset(action.ID)
	}

	s.concat()

	delete(s.actionMap, index)
	return
}

func (s *LockStep) concat() {
	for s.actions.Len() > 0 {
		front, _ := s.actions.At(0)
		element := front.(*LockStepAction)
		if element.dropped {
			s.actions.Pop()
			ActionSetPool.Put([]ActionSet(element.actions))
			delete(s.actionMap, element.index)
			continue
		}
		if !element.complete {
			return
		}
		s.actions.Pop()
		// concat
		for _, action := range element.actions {
			if current, ok := s.baseState[action.ID]; ok {
				s.step(current, action.Action)
			}
		}
		ActionSetPool.Put([]ActionSet(element.actions))
		delete(s.actionMap, element.index)
	}
}

func (s *LockStep) reset(id uint32) {
	element, ok := s.expectedState[id]
	if !ok {
		return
	}
	element.Value = s.clone(s.baseState[id].Value)
	s.actions.Range(func(i int, t interface{}) bool {
		stepAction := t.(*LockStepAction)
		if stepAction.dropped {
			return true
		}
		action, ok := stepAction.actions.FindAsIsSorted(id)
		if !ok {
			return true
		}
		s.step(element, action)
		return true
	})
}

func mapCopy(dst []ActionSet, src map[uint32]interface{}) {
	i := 0
	for id, action := range src {
		dst[i].ID = id
		dst[i].Action = &Action{
			Value: action,
		}
		i++
	}
	sort.Sort(ActionSlice(dst))
}

func mapSliceCopy(dst []ActionSet, src []ActionSet) {
	copy(dst, src)
	sort.Sort(ActionSlice(dst))
}
