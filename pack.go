package gopack

import (
	"github.com/handball811/gobits"
	"github.com/handball811/gopack/lib/holder"
)

type PackOp interface {
	Handler
	getEventSetter() holder.EventSetter
	getOperator() holder.Operator
}

type Handler interface {
	Get(segment gobits.ReaderFrom) int
	Complete(index int)
	Drop(index int)
}

type Pack struct {
	Handler

	holder *holder.Holder
}

func NewPack() (*Pack, error) {
	holder, err := holder.NewHolder()
	if err != nil {
		return nil, err
	}
	pack := &Pack{
		Handler: holder,
		holder:  holder,
	}
	return pack, nil
}

func (p *Pack) getEventSetter() holder.EventSetter {
	return p.holder
}

func (p *Pack) getOperator() holder.Operator {
	return p.holder
}
