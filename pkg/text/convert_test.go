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

package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCamelToKebab(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "single lowercase word",
			in:   "test",
			out:  "test",
		},
		{
			name: "multiple mixed words",
			in:   "testTestTest",
			out:  "test-test-test",
		},
		{
			name: "multiple uppercase words",
			in:   "TestTest",
			out:  "test-test",
		},
		{
			name: "multiple lowercase words",
			in:   "testtest",
			out:  "testtest",
		},
		{
			name: "multiple mixed words with number",
			in:   "test2Test",
			out:  "test2-test",
		},
		{
			name: "multiple lowercase words with number",
			in:   "test2test",
			out:  "test2test",
		},
		{
			name: "multiple lowercase words with dash",
			in:   "test-test",
			out:  "test-test",
		},
		{
			name: "multiple uppercase words with dash",
			in:   "Test-Test",
			out:  "test--test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, CamelToKebab(tt.in))
		})
	}
}
