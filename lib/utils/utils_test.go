package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExtendStr(t *testing.T) {
	assert.Equal(t, "foo", ExtendStr("", "foo", "bar"), "should extend a string (case 1)")
	assert.Equal(t, "bar", ExtendStr("", "", "bar"), "should extend a string (case 2)")
	assert.Equal(t, "bar", ExtendStr("bar", "", "foo"), "should extend a string (case 3)")
}

func Test_ExtendInt(t *testing.T) {
	assert.Equal(t, 1, ExtendInt(0, 1, 2), "should extend an integer (case 1)")
	assert.Equal(t, 2, ExtendInt(0, 0, 2), "should extend an integer (case 2)")
	assert.Equal(t, 1, ExtendInt(1, 0, 2), "should extend an integer (case 3)")
}
