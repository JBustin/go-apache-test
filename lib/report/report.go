package report

import (
	"fmt"
	"path/filepath"
	"sort"
)

type Report struct {
	tests    int
	failed   int
	success  int
	skips    int
	failures []failure
}

type failure struct {
	vhostName string
	suiteName string
	ruleName  string
	errorMsg  string
}

func New() Report {
	return Report{
		tests:    0,
		failed:   0,
		success:  0,
		skips:    0,
		failures: []failure{},
	}
}

func (r *Report) AddSkip() {
	r.tests++
	r.skips++
}

func (r *Report) Add(vhostName string, suiteName string, ruleName string, err error) failure {
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

func (r Report) String() string {
	sort.Slice(r.failures, func(i, j int) bool {
		return r.failures[i].vhostName < r.failures[j].vhostName
	})

	failures := ""
	for _, f := range r.failures {
		failures += fmt.Sprint(f) + "\n"
	}

	return fmt.Sprintf(`
		Tests:		%v
		Skips: 		%v
		Failed:		%v
		Success:	%v
		%v
	`, r.tests, r.skips, r.failed, r.success, failures)
}

func (f failure) String() string {
	return fmt.Sprintf(`
		%v > %v > %v
		%v
	`, filepath.Base(f.vhostName), f.suiteName, f.ruleName, f.errorMsg)
}
