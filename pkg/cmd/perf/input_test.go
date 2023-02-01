package perf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InputQueryFile(t *testing.T) {
	sample := `
metadata:
  table: numbers

statements:
- name: Q1
  query: LOAD {{ "HOME" | env }}
`
	input := InputQueryFile{}
	err := input.Decode([]byte(sample))
	assert.NoError(t, err)
	assert.Equal(t, "numbers", input.MetaData.Table)
	assert.Equal(t, input.Statements[0].Query, "LOAD {{ \"HOME\" | env }}")
	got, err := RenderQueryStatment(input.Statements[0].Query)
	assert.NoError(t, err)
	assert.Contains(t, got, "/")
}
