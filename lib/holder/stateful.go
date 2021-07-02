package holder

import (
	"github.com/handball811/gobits"
	"github.com/handball811/goflag"
)

// TODO: Unit Test
/*

状態を持つようなデータに採用される管理機能
途中での更新に反応する

*/

type StatefulData struct {
	Info

	index  uint32
	reader *gobits.FixedSegment
	cancel *gobits.FixedSegment
	check  func() bool
	state  uint32
	value  uint
}

func (d *StatefulData) GetReader() gobits.Reader {
	return d.reader
}

func (d *StatefulData) CancelReader() gobits.Reader {
	return d.cancel
}

func (d *StatefulData) State() uint32 {
	if d.check() {
		return d.state
	}
	return 0
}

func (d *StatefulData) Weight() uint {
	return uint(WeightOverhead) + uint(d.reader.Len())
}

func (d *StatefulData) Value() uint {
	return d.value
}

type StatefulHolder struct {
	Eventer
	op Operator

	// 一時保管データの管理用
	flag goflag.FlagOp
	data []StatefulData
}

func NewStatefulHolder(
	op Operator,
	setter EventSetter,
) (*StatefulHolder, error) {
	flag, err := goflag.NewFlags(uint(MaxStatefulHolderSize))
	if err != nil {
		return nil, err
	}
	holder := &StatefulHolder{
		op:   op,
		flag: flag,
		data: make([]StatefulData, MaxStatefulHolderSize),
	}
	setter.SetEvent(holder)

	// data の初期化
	for i := 0; i < MaxStatefulHolderSize; i++ {
		holder.data[i].reader = gobits.NewFixedSegment(gobits.NewSegmentWithSize(MaxByteSize))
	}
	return holder, nil
}

func (h *StatefulHolder) Set(
	reader gobits.Reader,
	cancel gobits.Reader,
	value uint,
	check func() bool,
) (int, error) {
	id, err := h.flag.FindAndUp()
	if err != nil {
		return -1, err
	}
	info := &h.data[id]
	index, err := h.op.Set(info)
	if err != nil {
		h.flag.Down(id)
		return -1, err
	}

	info.index = index
	info.reader.Reset()
	info.reader.ReadFrom(reader)
	info.cancel.Reset()
	info.cancel.ReadFrom(cancel)
	info.check = check
	info.state = 1
	info.value = value

	return id, nil
}

func (h *StatefulHolder) Update(
	id int,
	reader gobits.Reader,
) error {
	if ok, err := h.flag.Check(id); !ok || err != nil {
		if !ok {
			return ErrUpdateNotExists
		}
		return err
	}

	info := &h.data[id]
	info.reader.Reset()
	info.reader.ReadFrom(reader)
	info.state++
	return nil
}

func (h *StatefulHolder) Unset(
	id int,
) {
	h.flag.Down(id)

	info := &h.data[id]
	h.op.Unset(info.index)
}

func (h *StatefulHolder) OnGet(index int) {
	// 特にやることはない
}

func (h *StatefulHolder) OnComplete(index int) {
	// 特にやることはない
}

func (h *StatefulHolder) OnDrop(index int) {
	// 特にやることはない
}
