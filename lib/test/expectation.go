package test

import (
	"fmt"
)

type Expectations struct {
	Backend        string
	Host           string
	Pathname       string
	HttpStatusCode int
	Body           string
}

func (e Expectations) String() string {
	return fmt.Sprintf(`
	backend			%v
	host			%v
	pathname		%v
	httpStatusCode	%v
	body		%v
`,
		e.Backend,
		e.Host,
		e.Pathname,
		e.HttpStatusCode,
		len(e.Body),
	)
}
