package report

import (
	"fmt"
	"path/filepath"
	"sort"
)

type report struct {
	tests    int
	failed   int
	success  int
	failures []failure
}

type failure struct {
	vhostName string
	suiteName string
	ruleName  string
	errorMsg  string
}

func New() report {
	return report{
		tests:    0,
		failed:   0,
		success:  0,
		failures: []failure{},
	}
}

func (r *report) Add(vhostName string, suiteName string, ruleName string, err error) failure {
	r.tests++
	if err == nil {
		r.success++
		return failure{}
	}
	r.failed++
	f := failure{
		vhostName: vhostName,
		suiteName: suiteName,
		ruleName:  ruleName,
		errorMsg:  fmt.Sprint(err),
	}
	r.failures = append(r.failures, f)
	return f
}

func (r report) String() string {
	sort.Slice(r.failures, func(i, j int) bool {
		return r.failures[i].vhostName < r.failures[j].vhostName
	})

	failures := ""
	for _, f := range r.failures {
		failures += fmt.Sprint(f) + "\n"
	}

	return fmt.Sprintf(`
		Tests:		%v
		Failed:		%v
		Success:	%v
		%v
	`, r.tests, r.failed, r.success, failures)
}

func (f failure) String() string {
	return fmt.Sprintf(`
		%v > %v > %v
		%v
	`, filepath.Base(f.vhostName), f.suiteName, f.ruleName, f.errorMsg)
}
