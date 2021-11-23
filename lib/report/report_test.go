package report

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Report(t *testing.T) {
	r := New()

	r.Add("mydomain.com.conf", "balance to new stack", "homepage", nil)
	r.Add("mydomain.com.conf", "balance to new stack", "product-list", nil)
	r.Add("mydomain.com.conf", "balance to new stack", "product", nil)
	r.Add("mydomain.com.conf", "balance to new stack", "faq", fmt.Errorf("invalid data"))

	assert.Equal(t, 4, r.tests, "should compute all tests")
	assert.Equal(t, 1, r.failed, "should compute all failed tests")
	assert.Equal(t, 3, r.success, "should compute all success tests")
	assert.Equal(t, 1, len(r.failures), "should store all failures")
	assert.Equal(
		t,
		"\n\t\tTests:\t\t4\n\t\tSkips: \t\t0\n\t\tFailed:\t\t1\n\t\tSuccess:\t3\n\t\t\n\t\tmydomain.com.conf > balance to new stack > faq\n\t\tinvalid data\n\t\n\n\t",
		r.String(),
		"should generate a valid report",
	)
}
