package parser

func Compile(re string) (*State, error) {
	post, err := re2post(transform(re))
	if err != nil {
		return nil, err
	}

	nfa, err := post2nfa(post)
	if err != nil {
		return nil, err
	}

	return nfa, nil
}
