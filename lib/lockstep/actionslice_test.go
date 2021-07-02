package lockstep_test

import (
	"testing"

	"github.com/handball811/gopack/lib/lockstep"
	. "github.com/handball811/gopack/lib/lockstep"
	"github.com/stretchr/testify/assert"
)

func TestActionSliceFind1(t *testing.T) {
	// setup
	target := ActionSlice([]ActionSet{
		{
			ID:     1,
			Action: &lockstep.Action{"1"},
		},
		{
			ID:     2,
			Action: &lockstep.Action{"2"},
		},
		{
			ID:     3,
			Action: &lockstep.Action{"3"},
		},
		{
			ID:     4,
			Action: &lockstep.Action{"4"},
		},
		{
			ID:     5,
			Action: &lockstep.Action{"5"},
		},
		{
			ID:     6,
			Action: &lockstep.Action{"6"},
		},
	})

	// when
	s, ok := target.FindAsIsSorted(4)

	// then
	assert.Equal(t, "4", s)
	assert.Equal(t, true, ok)
}

func TestActionSliceFind2(t *testing.T) {
	// setup
	target := ActionSlice([]ActionSet{
		{
			ID:     1,
			Action: &lockstep.Action{"1"},
		},
		{
			ID:     2,
			Action: &lockstep.Action{"2"},
		},
		{
			ID:     3,
			Action: &lockstep.Action{"3"},
		},
	})

	// when
	s, ok := target.FindAsIsSorted(1)

	// then
	assert.Equal(t, "1", s)
	assert.Equal(t, true, ok)
}

func TestActionSliceFind3(t *testing.T) {
	// setup
	target := ActionSlice([]ActionSet{
		{
			ID:     1,
			Action: &lockstep.Action{"1"},
		},
		{
			ID:     2,
			Action: &lockstep.Action{"2"},
		},
		{
			ID:     3,
			Action: &lockstep.Action{"3"},
		},
	})

	// when
	s, ok := target.FindAsIsSorted(3)

	// then
	assert.Equal(t, "3", s)
	assert.Equal(t, true, ok)
}

func TestActionSliceNotFound(t *testing.T) {
	// setup
	target := ActionSlice([]ActionSet{
		{
			ID:     1,
			Action: &lockstep.Action{"1"},
		},
		{
			ID:     2,
			Action: &lockstep.Action{"2"},
		},
		{
			ID:     3,
			Action: &lockstep.Action{"3"},
		},
		{
			ID:     5,
			Action: &lockstep.Action{"5"},
		},
		{
			ID:     6,
			Action: &lockstep.Action{"6"},
		},
	})

	// when
	s, ok := target.FindAsIsSorted(4)

	// then
	assert.Nil(t, s)
	assert.Equal(t, false, ok)
}
