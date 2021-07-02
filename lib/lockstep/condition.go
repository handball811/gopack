package lockstep

type ConditionLockStep struct {
	*LockStep
}

func NewConditionLockStep(
	clone func(from interface{}) interface{},
) *ConditionLockStep {
	return &ConditionLockStep{
		LockStep: NewLockStep(
			func(current *Element, action *Action) {
				current.Value = clone(action.Value)
			},
			clone,
		),
	}
}
