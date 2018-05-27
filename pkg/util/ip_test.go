package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetLocalIP(t *testing.T) {
	assert.NotEmpty(t, GetLocalIP())
}
