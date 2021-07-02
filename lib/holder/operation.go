package holder

import (
	"errors"

	"github.com/handball811/gobits"
)

var (
	// ビットのサイズに加えて伝送時のヘッダサイズを指定する
	WeightOverhead = 8
	// 各情報のビットサイズの上限
	MaxByteSize           = 128
	MaxOnetimeHolderSize  = 2048
	MaxSingleHolderSize   = 2048
	MaxStatefulHolderSize = 2048

	EmptyReader = &emptyReader{}

	ErrUpdateNotExists = errors.New("The id you set is not currently used")
)

type InfoSet struct {
	id   uint32
	info Info
}

// 読み取り項目を追加したり、更新を通知したりする
type Operator interface {
	Set(info Info) (uint32, error)
	Unset(id uint32)
}

type EventSetter interface {
	SetEvent(e Eventer)
}

// 項目をまとめて受け取り届いたかを通知する場所
type Handler interface {
	Get(segment gobits.ReaderFrom) int
	Complete(index int)
	Drop(index int)
}

type Eventer interface {
	// After Get Data
	OnGet(index int)
	// Before Complete
	OnComplete(index int)
	// Before Drop
	OnDrop(index int)
}

// 登録する際に保持しなくてはいけないデータ
type Info interface {
	GetReader() gobits.Reader
	CancelReader() gobits.Reader
	State() uint32
	Weight() uint
	Value() uint
}

type emptyReader struct {
	gobits.Reader
}

func (r *emptyReader) Read(b *gobits.Slice) (int, error) {
	return 0, nil
}
