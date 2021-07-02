package packager

type Action struct {
	id    uint32
	state uint32
}

func (a *Action) ID() uint32 {
	return a.id
}

func (a *Action) State() uint32 {
	return a.state
}
