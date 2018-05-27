package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RemovePrefix(t *testing.T) {
	t.Run("without prefix", func(t *testing.T) {
		str, ok := RemovePrefix("prefix-", "aaa-foo")
		assert.False(t, ok)
		assert.Equal(t, "aaa-foo", str)
	})

	t.Run("with prefix", func(t *testing.T) {
		str, ok := RemovePrefix("prefix-", "prefix-foo")
		assert.True(t, ok)
		assert.Equal(t, "foo", str)
	})
}

func Test_ParseHelpFile(t *testing.T) {
	data := []byte("Summary line\nDescription goes here.\nAnd more lines...")
	summary, description := ParseHelpFile(data)

	assert.Equal(t, "Summary line", summary)
	assert.Contains(t, description, "Summary line")
	assert.Contains(t, description, "Description goes here")
}
