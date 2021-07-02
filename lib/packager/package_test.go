package packager_test

import (
	"testing"

	. "github.com/handball811/gopack/lib/packager"
	"github.com/stretchr/testify/assert"
)

var (
	id0 uint32 = 10
	id1 uint32 = 20
	id2 uint32 = 30
)

type MockChecker struct {
	Result map[uint32]uint32
}

func NewMockChecker() *MockChecker {
	return &MockChecker{
		Result: make(map[uint32]uint32),
	}
}

func (c *MockChecker) Check(id uint32) uint32 {
	return c.Result[id]
}

type MockSetChecker struct{}

func NewMockSetChecker() *MockSetChecker {
	return &MockSetChecker{}
}

func (c *MockSetChecker) Check(ids []uint32) int {
	return len(ids)
}

func TestUpdate(t *testing.T) {
	// setup
	actions := make([]Action, 1)
	checker := NewMockChecker()
	setchecker := NewMockSetChecker()
	target := NewPackage(1024, checker, setchecker)
	target.Add(id0)
	target.Add(id1)
	target.Add(id2)
	checker.Result[id0] = 1
	checker.Result[id1] = 0
	checker.Result[id2] = 0

	// when
	size, _ := target.Update(actions)

	// then
	assert.Equal(t, 1, size)
	assert.Equal(t, id0, actions[0].ID())
	assert.Equal(t, uint32(1), actions[0].State())
}

func TestUpdate2(t *testing.T) {
	// setup
	actions := make([]Action, 1)
	checker := NewMockChecker()
	setchecker := NewMockSetChecker()
	target := NewPackage(1024, checker, setchecker)
	target.Add(id0)
	target.Add(id1)
	target.Add(id2)
	checker.Result[id0] = 1
	checker.Result[id1] = 0
	checker.Result[id2] = 0

	// when
	target.Update(actions)
	size, _ := target.Update(actions)

	// then
	assert.Equal(t, 0, size)
}

func TestDrop(t *testing.T) {
	// setup
	actions := make([]Action, 1)
	checker := NewMockChecker()
	setchecker := NewMockSetChecker()
	target := NewPackage(1024, checker, setchecker)
	target.Add(id0)
	target.Add(id1)
	target.Add(id2)
	checker.Result[id0] = 1
	checker.Result[id1] = 0
	checker.Result[id2] = 0

	// when
	_, index := target.Update(actions)
	target.Drop(index)
	size, _ := target.Update(actions)

	// then
	assert.Equal(t, 1, size)
	assert.Equal(t, id0, actions[0].ID())
	assert.Equal(t, uint32(1), actions[0].State())
}

func TestDetectChange(t *testing.T) {
	// setup
	actions := make([]Action, 1)
	checker := NewMockChecker()
	setchecker := NewMockSetChecker()
	target := NewPackage(1024, checker, setchecker)
	target.Add(id0)
	target.Add(id1)
	target.Add(id2)
	checker.Result[id0] = 1
	checker.Result[id1] = 0
	checker.Result[id2] = 0

	// when
	target.Update(actions)
	checker.Result[id0] = 0
	size, _ := target.Update(actions)

	// then
	assert.Equal(t, 1, size)
	assert.Equal(t, id0, actions[0].ID())
	assert.Equal(t, uint32(0), actions[0].State())
}

func TestComplete(t *testing.T) {
	// setup
	actions := make([]Action, 1)
	checker := NewMockChecker()
	setchecker := NewMockSetChecker()
	target := NewPackage(1024, checker, setchecker)
	target.Add(id0)
	target.Add(id1)
	target.Add(id2)
	checker.Result[id0] = 1
	checker.Result[id1] = 0
	checker.Result[id2] = 0

	// when
	_, index := target.Update(actions)
	target.Complete(index)
	size, _ := target.Update(actions)

	// then
	assert.Equal(t, 0, size)
}

func TestDropAndUpdate(t *testing.T) {
	// setup
	actions := make([]Action, 1)
	checker := NewMockChecker()
	setchecker := NewMockSetChecker()
	target := NewPackage(1024, checker, setchecker)
	target.Add(id0)
	target.Add(id1)
	target.Add(id2)
	checker.Result[id0] = 1
	checker.Result[id1] = 0
	checker.Result[id2] = 0

	// when
	_, index := target.Update(actions)
	checker.Result[id0] = 0
	_, index2 := target.Update(actions)
	target.Drop(index2)
	target.Drop(index)
	size, _ := target.Update(actions)

	// then
	assert.Equal(t, 0, size)
}

func TestDropAndNotUpdate(t *testing.T) {
	// setup
	actions := make([]Action, 1)
	checker := NewMockChecker()
	setchecker := NewMockSetChecker()
	target := NewPackage(1024, checker, setchecker)
	target.Add(id0)
	target.Add(id1)
	target.Add(id2)
	checker.Result[id0] = 1
	checker.Result[id1] = 0
	checker.Result[id2] = 0

	// when
	_, index := target.Update(actions)
	checker.Result[id0] = 0
	_, index2 := target.Update(actions)
	target.Complete(index2)
	target.Drop(index)
	size, _ := target.Update(actions)

	// then
	assert.Equal(t, 0, size)
}

func TestCompleteAndUpdate(t *testing.T) {
	// setup
	actions := make([]Action, 1)
	checker := NewMockChecker()
	setchecker := NewMockSetChecker()
	target := NewPackage(1024, checker, setchecker)
	target.Add(id0)
	target.Add(id1)
	target.Add(id2)
	checker.Result[id0] = 1
	checker.Result[id1] = 0
	checker.Result[id2] = 0

	// when
	_, index := target.Update(actions)
	checker.Result[id0] = 0
	_, index2 := target.Update(actions)
	target.Complete(index2)
	target.Complete(index)
	size, _ := target.Update(actions)

	// then
	assert.Equal(t, 0, size)
}

func BenchmarkUpdateAndDrop(b *testing.B) {
	// setup
	var size uint32 = 32
	actions := make([]Action, size)
	checker := NewMockChecker()
	setchecker := NewMockSetChecker()
	target := NewPackage(1024, checker, setchecker)
	var i uint32
	for i = 0; i < size; i++ {
		target.Add(i)
		checker.Result[i] = 1
	}
	b.ResetTimer()
	for i = 0; i < uint32(b.N); i++ {
		_, index := target.Update(actions)
		_, index2 := target.Update(actions)
		_, index3 := target.Update(actions)
		_, index4 := target.Update(actions)
		target.Drop(index4)
		target.Drop(index3)
		target.Drop(index2)
		target.Drop(index)
	}
}
