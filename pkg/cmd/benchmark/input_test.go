// Copyright 2022 Datafuse Labs.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package benchmark

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
