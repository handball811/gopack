package holder

import (
	"github.com/handball811/gobits"
	"github.com/handball811/goflag"
	"github.com/handball811/gopack/lib/packager"
)

/*
情報の登録
取得範囲かのチェック
取得対象かのチェック

外側にする


*/

const (
	InfoSize = 2048

	PackageSize = 512
	PackageBits = 1300 * 8
)

type Holder struct {
	Operator
	EventSetter
	Handler
	// 登録された情報のIDをもとにつないである
	flag  goflag.FlagOp
	infos []InfoSet

	pack    *packager.SizedPackage
	actions []packager.Action

	eventers []Eventer
}

func NewHolder() (*Holder, error) {
	flag, err := goflag.NewFlags(InfoSize)
	if err != nil {
		return nil, err
	}
	h := &Holder{
		flag:     flag,
		infos:    make([]InfoSet, InfoSize),
		actions:  make([]packager.Action, PackageSize),
		eventers: make([]Eventer, 0),
	}
	h.pack = packager.NewSizedPackage(
		PackageSize,
		PackageBits,
		h.check,
		h.wvGetter,
	)
	return h, nil
}

// 登録
func (h *Holder) Set(
	info Info,
) (uint32, error) {
	// 始めに保管できるすぺーがあるかを確認する
	index, err := h.flag.FindAndUp()
	if err != nil {
		return 0, err
	}
	var id uint32 = uint32(index)
	h.infos[index].info = info
	h.infos[index].id = id

	// Packageに登録する
	h.pack.Add(id)

	return id, nil
}

// 内容の更新を通知
/*
func (h *Holder) Update(
	id uint32,
) {
	if ok, err := h.flag.Check(int(id)); !ok || err != nil {
		// 未使用なら更新できないようにする
		return
	}
	// 更新するために一度削除し追加する
	h.pack.Remove(id)
	h.pack.Add(id)
}*/

// 削除
func (h *Holder) Unset(id uint32) {
	if ok, err := h.flag.Check(int(id)); !ok || err != nil {
		// 未使用なら更新できないようにする
		return
	}
	//Packageから削除する
	h.pack.Remove(id)

	// スペースから削除する
	h.flag.Down(int(id))
}

// イベントの登録
func (h *Holder) SetEvent(e Eventer) {
	h.eventers = append(h.eventers, e)
}

func (h *Holder) check(id uint32) uint32 {
	index := int(id)
	if ok, err := h.flag.Check(index); !ok || err != nil {
		// 未使用なら更新できないようにする
		return 0
	}
	return h.infos[index].info.State()
}

func (h *Holder) wvGetter(id uint32) (uint, uint, bool) {
	index := int(id)
	if ok, err := h.flag.Check(index); !ok || err != nil {
		// 未使用なら更新できないようにする
		return 0, 0, false
	}
	set := h.infos[index].info
	return set.Weight(), set.Value(), true
}

// 内容を取得する
// 届いたかを通知するのに利用する
func (h *Holder) Get(segment gobits.ReaderFrom) int {
	size, index := h.pack.Update(h.actions)
	h.onGet(index)
	for _, action := range h.actions[:size] {
		index := int(action.ID())
		if action.State() == 0 {
			segment.ReadFrom(h.infos[index].info.CancelReader())
		} else {
			segment.ReadFrom(h.infos[index].info.GetReader())
		}
	}
	return index
}

// 通知用の関数
func (h *Holder) onGet(index int) {
	for _, eventer := range h.eventers {
		eventer.OnGet(index)
	}
}

func (h *Holder) Complete(index int) {
	h.onComplete(index)
	h.pack.Complete(index)
}

// 通知用の関数
func (h *Holder) onComplete(index int) {
	for _, eventer := range h.eventers {
		eventer.OnComplete(index)
	}
}

func (h *Holder) Drop(index int) {
	h.onDrop(index)
	h.pack.Drop(index)
}

// 通知用の関数
func (h *Holder) onDrop(index int) {
	for _, eventer := range h.eventers {
		eventer.OnDrop(index)
	}
}
