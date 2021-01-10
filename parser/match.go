package parser

type List struct {
	states []*State
	listId int
}

func (l *List) add(s *State) {
	if s.ListId == l.listId {
		return
	}

	if s.char == CharSplit {
		l.add(s.out)
		l.add(s.out2)
	} else {
		s.ListId = l.listId
		l.states = append(l.states, s)
	}
}

func (l *List) step(c rune) *List {
	lst := &List{
		states: make([]*State, 0, len(l.states)),
		listId: l.listId + 1,
	}

	for _, s := range l.states {
		switch s.char {
		case CharSplit:
			lst.add(s)
		case CharMatch:
		default:
			if s.char == c || s.char == '.' {
				lst.add(s.out)
			}
		}
	}
	return lst
}

func (l *List) flow() *List {
	lst := &List{
		states: make([]*State, 0, len(l.states)),
		listId: l.listId + 1,
	}

	for _, s := range l.states {
		lst.add(s)
	}
	return lst
}

func (l *List) isMatch() bool {
	for _, s := range l.states {
		if s.char == CharMatch {
			return true
		}
	}
	return false
}

func (s *State) Match(text string) bool {
	l := &List{states: []*State{s}}

	l = l.flow()
	for _, c := range text {
		l = l.step(c)
	}
	l = l.flow()
	return l.isMatch()
}
