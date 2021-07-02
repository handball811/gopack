package packager_test

import (
	"testing"

	. "github.com/handball811/gopack/lib/packager"
	"github.com/stretchr/testify/assert"
)

type setter struct {
	check  map[uint32]uint32
	weight map[uint32]uint
	value  map[uint32]uint
}

func (s *setter) Set(
	id uint32,
	check uint32,
	weight uint,
	value uint,
) {
	s.check[id] = check
	s.weight[id] = weight
	s.value[id] = value
}

func (s *setter) Check(id uint32) uint32 {
	if check, ok := s.check[id]; ok {
		return check
	}
	return 0
}
func (s *setter) WVGet(id uint32) (uint, uint, bool) {
	if _, ok := s.check[id]; ok {
		return s.weight[id], s.value[id], true
	}
	return 0, 0, false
}

func generatePackage() (*SizedPackage, *setter) {
	s := &setter{
		check:  make(map[uint32]uint32),
		weight: make(map[uint32]uint),
		value:  make(map[uint32]uint),
	}
	return NewSizedPackage(
		4,
		256,
		s.Check,
		s.WVGet,
	), s
}

func TestSizedPackage(t *testing.T) {
	// setup
	var actionSize int
	target, setter := generatePackage()
	actions := make([]Action, 8)

	// when
	var id0 uint32 = 10
	target.Add(id0)
	setter.Set(id0, 1, 10, 10)
	actionSize, _ = target.Update(actions)

	//then
	assert.Equal(t, 1, actionSize)
	assert.Equal(t, id0, actions[0].ID())
	assert.Equal(t, uint32(1), actions[0].State())
}
