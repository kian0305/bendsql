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

package completion

import (
	"github.com/databendcloud/bendsql/pkg/iostreams"
	"strings"
	"testing"

	"github.com/google/shlex"
	"github.com/spf13/cobra"
)

func TestNewCmdCompletion(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		wantOut string
		wantErr string
	}{
		{
			name:    "no arguments",
			args:    "completion",
			wantOut: "complete -o default -F __start_bendsql bendsql",
		},
		{
			name:    "zsh completion",
			args:    "completion -s zsh",
			wantOut: "#compdef bendsql",
		},
		{
			name:    "fish completion",
			args:    "completion -s fish",
			wantOut: "complete -c bendsql ",
		},
		{
			name:    "PowerShell completion",
			args:    "completion -s powershell",
			wantOut: "Register-ArgumentCompleter",
		},
		{
			name:    "unsupported shell",
			args:    "completion -s csh",
			wantErr: "invalid argument \"csh\" for \"-s, --shell\" flag: valid values are {bash|zsh|fish|powershell}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ios, _, stdout, stderr := iostreams.Test()
			completeCmd := NewCmdCompletion(ios)
			rootCmd := &cobra.Command{Use: "bendsql"}
			rootCmd.AddCommand(completeCmd)

			argv, err := shlex.Split(tt.args)
			if err != nil {
				t.Fatalf("argument splitting error: %v", err)
			}
			rootCmd.SetArgs(argv)
			rootCmd.SetOut(stderr)
			rootCmd.SetErr(stderr)

			_, err = rootCmd.ExecuteC()
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %q", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("error executing command: %v", err)
			}

			if !strings.Contains(stdout.String(), tt.wantOut) {
				t.Errorf("completion output did not match:\n%s", stdout.String())
			}
			if len(stderr.String()) > 0 {
				t.Errorf("expected nothing on stderr, got %q", stderr.String())
			}
		})
	}
}
