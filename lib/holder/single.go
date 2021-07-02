package holder

// TODO: Unit Test

import (
	"github.com/handball811/gobits"
	"github.com/handball811/goflag"
	"github.com/handball811/goring"
)

/*

確実に一回の送信したいときに採用されるデータを管理するための機能

*/

type SingleData struct {
	Info
	index  uint32
	reader *gobits.FixedSegment
	value  uint
}

func (d *SingleData) GetReader() gobits.Reader {
	return d.reader
}

func (d *SingleData) CancelReader() gobits.Reader {
	//キャンセルされることはないが呼び出されたとしても無視する
	return EmptyReader
}

func (d *SingleData) State() uint32 {
	// 常に固定値で伝えるべき情報として保持する
	return 1
}

func (d *SingleData) Weight() uint {
	return uint(WeightOverhead) + uint(d.reader.Len())
}

func (d *SingleData) Value() uint {
	return d.value
}

type SingleHolder struct {
	Eventer
	op Operator

	// 一時保管データの管理用
	flag goflag.FlagOp
	data []SingleData
	ring goring.RingOp // 使用されているインデックスを保管する
}

func NewSingleHolder(
	op Operator,
	setter EventSetter,
) (*SingleHolder, error) {
	flag, err := goflag.NewFlags(uint(MaxSingleHolderSize))
	if err != nil {
		return nil, err
	}
	holder := &SingleHolder{
		op:   op,
		flag: flag,
		data: make([]SingleData, MaxSingleHolderSize),
		ring: goring.NewRingWithSize(MaxSingleHolderSize),
	}
	setter.SetEvent(holder)

	// data の初期化
	for i := 0; i < MaxSingleHolderSize; i++ {
		holder.data[i].reader = gobits.NewFixedSegment(gobits.NewSegmentWithSize(MaxByteSize))
	}
	return holder, nil
}

// セットした時点で伝えるべき情報として認識される
func (h *SingleHolder) Set(
	reader gobits.Reader, // 書き出したいデータの中身
	value uint,
) error {
	// 空きを探す
	id, err := h.flag.FindAndUp()
	if err != nil {
		return err
	}
	// 登録を開始する
	info := h.data[id]
	index, err := h.op.Set(&info)
	if err != nil {
		// フラグを解除する
		h.flag.Down(id)
		return err
	}

	info.index = index
	info.reader.Reset()
	info.reader.ReadFrom(reader)
	info.value = value

	h.ring.Push(id)
	return nil
}

func (h *SingleHolder) OnGet(index int) {
	// 特にやることはない
}

func (h *SingleHolder) OnComplete(index int) {
	// 特にやることはない
}

func (h *SingleHolder) OnDrop(index int) {
	// 特にやることはない
}
