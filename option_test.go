package fuzzyfinder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindWithDefaultIndex(t *testing.T) {
	opt := defaultOption
	assert.Equal(t, opt.defaultIndex, -1)

	WithDefaultIndex(10)(&opt)
	assert.Equal(t, opt.defaultIndex, 10)
}