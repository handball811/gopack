package packager

import "go.uber.org/atomic"

/*

Control

OnOff
Weight
Value

*/

type checker struct {
	*SizedPackage
}

func (c *checker) Check(id uint32) uint32 {
	return c.checker(id)
}

type data struct {
	id uint32
	w  uint
	v  uint
}

type setChecker struct {
	*SizedPackage
	maxWeight uint
	datas     []*data
}

func (c *setChecker) Check(ids []uint32) int {
	// check if exists
	size := 0
	for _, id := range ids {
		if w, v, b := c.wvGetter(id); b {
			d := c.datas[size]
			d.id = id
			d.w = w
			d.v = v
			size++
		}
	}
	// napsack? <- cost merit...
	size = 0
	var total uint = 0
	for _, d := range c.datas {
		if total+d.w > c.maxWeight {
			total += d.w
			size++
			break
		}
	}
	return size
}

type SizedPackage struct {
	*Package
	idGenerator *atomic.Uint32
	checker     func(id uint32) uint32
	wvGetter    func(id uint32) (uint, uint, bool)
}

func NewSizedPackage(
	maxSize int, // 取得するデータ個数の最大量
	maxWeight int,
	check func(id uint32) uint32,
	wvGetter func(id uint32) (uint, uint, bool), // weight value error
) *SizedPackage {
	knapsack := &SizedPackage{
		idGenerator: atomic.NewUint32(1),
		checker:     check,
		wvGetter:    wvGetter,
	}
	sc := &setChecker{
		SizedPackage: knapsack,
		datas:        make([]*data, maxSize),
	}
	for i := 0; i < maxSize; i++ {
		sc.datas[i] = &data{}
	}
	pack := NewPackage(
		maxSize,
		&checker{
			SizedPackage: knapsack,
		},
		sc)

	knapsack.Package = pack
	return knapsack
}
