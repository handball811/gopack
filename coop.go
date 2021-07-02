package gopack

// TODO: Unit Test

import "github.com/handball811/gobits"

var (
	IDSize = 1024
)

type Writer interface {
	Write(
		id int,
		segment gobits.WriterTo,
	) error
	Completes(ids []int) int
	Drops(ids []int) int
}

type WriterPack struct {
	w   Writer
	p   PackOp
	s   *gobits.Segment
	ids []int
}

func NewWriterPack(
	w Writer,
	p PackOp,
	bufferSize int,
) *WriterPack {
	pack := &WriterPack{
		w:   w,
		p:   p,
		s:   gobits.NewSegmentWithSize(bufferSize),
		ids: make([]int, IDSize),
	}
	return pack
}

// 書き込みとチェックを行う
func (p *WriterPack) Update() error {
	// Writerに書き込み
	p.s.Reset()
	id := p.p.Get(p.s)
	err := p.w.Write(id, p.s)
	if err != nil {
		p.p.Drop(id)
		return err
	}
	p.Check()
	return nil
}

func (p *WriterPack) Check() {
	{
		// Complete
		size := p.w.Completes(p.ids)
		for _, id := range p.ids[:size] {
			p.p.Complete(id)
		}
	}
	{
		// Drop
		size := p.w.Drops(p.ids)
		for _, id := range p.ids[:size] {
			p.p.Drop(id)
		}
	}
}
