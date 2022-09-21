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

package root

import (
	"github.com/databendcloud/bendsql/pkg/iostreams"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHelpTopic(t *testing.T) {
	tests := []struct {
		name     string
		topic    string
		args     []string
		flags    []string
		wantsErr bool
	}{
		{
			name:     "valid topic",
			topic:    "environment",
			args:     []string{},
			flags:    []string{},
			wantsErr: false,
		},
		{
			name:     "invalid topic",
			topic:    "invalid",
			args:     []string{},
			flags:    []string{},
			wantsErr: false,
		},
		{
			name:     "more than zero args",
			topic:    "environment",
			args:     []string{"invalid"},
			flags:    []string{},
			wantsErr: false,
		},
		{
			name:     "more than zero flags",
			topic:    "environment",
			args:     []string{},
			flags:    []string{"--invalid"},
			wantsErr: true,
		},
		{
			name:     "help arg",
			topic:    "environment",
			args:     []string{"help"},
			flags:    []string{},
			wantsErr: false,
		},
		{
			name:     "help flag",
			topic:    "environment",
			args:     []string{},
			flags:    []string{"--help"},
			wantsErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ios, _, _, stderr := iostreams.Test()

			cmd := NewHelpTopic(ios, tt.topic)
			cmd.SetArgs(append(tt.args, tt.flags...))
			cmd.SetOut(stderr)
			cmd.SetErr(stderr)

			_, err := cmd.ExecuteC()
			if tt.wantsErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
