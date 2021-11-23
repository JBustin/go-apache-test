package test

import "fmt"

type Suite struct {
	Name         string
	Rules        []Rule
	Skip         bool
	IsDocRoot    bool
	Headers      map[string]string
	Expectations Expectations
}

func (s Suite) String() string {
	return fmt.Sprintf(`
	name	%v
	rules	%v
	skip	%v
	isDocRoot	%v
	headers		%v
	expectations%v
`, s.Name, len(s.Rules), s.Skip, s.IsDocRoot, s.Headers, s.Expectations)
}
