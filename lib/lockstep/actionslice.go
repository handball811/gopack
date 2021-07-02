package lockstep

type ActionSet struct {
	ID     uint32
	Action *Action
}

type ActionSlice []ActionSet

func (s ActionSlice) FindAsIsSorted(id uint32) (*Action, bool) {
	var top, mid, bot int
	top = 0
	bot = len(s)
	for bot-top > 1 {
		mid = (top + bot) / 2
		if id == s[mid].ID {
			return s[mid].Action, true
		}
		if id > s[mid].ID {
			top = mid
		} else {
			bot = mid
		}
	}
	mid = (top + bot) / 2
	if id == s[mid].ID {
		return s[mid].Action, true
	}
	return nil, false
}

func (s ActionSlice) Len() int {
	return len(s)
}

func (s ActionSlice) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

func (s ActionSlice) Swap(i, j int) {
	s[i].ID, s[j].ID = s[j].ID, s[i].ID
	s[i].Action, s[j].Action = s[j].Action, s[i].Action
}
