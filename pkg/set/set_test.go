package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	s := NewSet()
	s.Add("foo")
	assert.True(t, s.Contains("foo"))
}

func TestAdd(t *testing.T) {
	s := NewSet()
	s.Add("foo")
	assert.False(t, s.Add("foo"))
}

func TestRemove(t *testing.T) {
	s := NewSet()
	s.Add("foo")
	assert.True(t, s.Contains("foo"))
	s.Remove("foo")
	assert.False(t, s.Contains("foo"))
}
