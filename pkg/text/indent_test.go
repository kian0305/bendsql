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

import "testing"

func Test_Indent(t *testing.T) {
	type args struct {
		s      string
		indent string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{
				s:      "",
				indent: "--",
			},
			want: "",
		},
		{
			name: "blank",
			args: args{
				s:      "\n",
				indent: "--",
			},
			want: "\n",
		},
		{
			name: "indent",
			args: args{
				s:      "one\ntwo\nthree",
				indent: "--",
			},
			want: "--one\n--two\n--three",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Indent(tt.args.s, tt.args.indent); got != tt.want {
				t.Errorf("indent() = %q, want %q", got, tt.want)
			}
		})
	}
}
