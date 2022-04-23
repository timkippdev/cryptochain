package blockchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Hash(t *testing.T) {
	input := "test"
	expectedOutput := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	assert.Equal(t, expectedOutput, Hash(input))
}

func Test_Hash_GenerateSameHashWithMultipleArgumentsRegardlessOfOrder(t *testing.T) {
	h1 := Hash("foo", "bar")
	h2 := Hash("bar", "foo")
	assert.Equal(t, h1, h2)
}
