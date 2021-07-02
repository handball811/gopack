package gopack

import (
	"github.com/handball811/gobits"
	"github.com/handball811/gopack/lib/holder"
)

type OnetimeHandler interface {
	Set(reader gobits.Reader, value uint) error
}

func GenerateOnetimeHandler(p PackOp) OnetimeHandler {
	h, _ := holder.NewOnetimeHolder(
		p.getOperator(),
		p.getEventSetter(),
	)
	return h
}

type SingleHandler interface {
	Set(reader gobits.Reader, value uint) error
}

func GenerateSingleHandler(p PackOp) SingleHandler {
	h, _ := holder.NewSingleHolder(
		p.getOperator(),
		p.getEventSetter(),
	)
	return h
}

type StatefulHandler interface {
	Set(
		reader gobits.Reader,
		cancel gobits.Reader,
		value uint,
		check func() bool,
	) (int, error)

	Update(
		id int,
		reader gobits.Reader,
	) error

	Unset(
		id int,
	)
}

func GenerateStatefulHandler(p PackOp) StatefulHandler {
	h, _ := holder.NewStatefulHolder(
		p.getOperator(),
		p.getEventSetter(),
	)
	return h
}
