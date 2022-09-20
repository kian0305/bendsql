package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDebugEnabled(t *testing.T) {
	type arg struct {
		value string
		debug bool
	}

	args := []arg{{
		value: "1",
		debug: true,
	},
		{
			value: "0",
			debug: false,
		},
		{
			value: "1",
			debug: true,
		},
		{
			value: "0",
			debug: false,
		},
	}

	for _, tt := range args {
		err := os.Setenv("DEBUG", tt.value)
		assert.NoError(t, err)
		assert.Equal(t, tt.debug, IsDebugEnabled())
	}
}
