package holder

// TODO: Unit Test

import (
	"github.com/handball811/gobits"
	"github.com/handball811/goflag"
	"github.com/handball811/goring"
)

/*

一回の送信だけに採用されるデータを管理するための機能

*/

type OnetimeData struct {
	Info
	index  uint32
	reader *gobits.FixedSegment
	value  uint
}

func (d *OnetimeData) GetReader() gobits.Reader {
	return d.reader
}

func (d *OnetimeData) CancelReader() gobits.Reader {
	//キャンセルされることはないが呼び出されたとしても無視する
	return EmptyReader
}

func (d *OnetimeData) State() uint32 {
	// 常に固定値で伝えるべき情報として保持する
	return 1
}

func (d *OnetimeData) Weight() uint {
	return uint(WeightOverhead) + uint(d.reader.Len())
}

func (d *OnetimeData) Value() uint {
	return d.value
}

type OnetimeHolder struct {
	Eventer
	op Operator

	// 一時保管データの管理用
	flag goflag.FlagOp
	data []OnetimeData
	ring goring.RingOp // 使用されているインデックスを保管する
}

func NewOnetimeHolder(
	op Operator,
	setter EventSetter,
) (*OnetimeHolder, error) {
	flag, err := goflag.NewFlags(uint(MaxOnetimeHolderSize))
	if err != nil {
		return nil, err
	}
	holder := &OnetimeHolder{
		op:   op,
		flag: flag,
		data: make([]OnetimeData, MaxOnetimeHolderSize),
		ring: goring.NewRingWithSize(MaxOnetimeHolderSize),
	}
	setter.SetEvent(holder)

	// data の初期化
	for i := 0; i < MaxOnetimeHolderSize; i++ {
		holder.data[i].reader = gobits.NewFixedSegment(gobits.NewSegmentWithSize(MaxByteSize))
	}
	return holder, nil
}

// セットした時点で伝えるべき情報として認識される
func (h *OnetimeHolder) Set(
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

func (h *OnetimeHolder) OnGet(index int) {
	// 今存在するデータを全て削除する
	h.ring.Range(func(i int, t interface{}) bool {
		id := t.(int)
		h.flag.Down(id)

		info := &h.data[id]
		h.op.Unset(info.index)
		return true
	})
	h.ring.Clean()
}

func (h *OnetimeHolder) OnComplete(index int) {
	// 特にやることはない
}

func (h *OnetimeHolder) OnDrop(index int) {
	// 特にやることはない
}
