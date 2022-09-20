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
	"strings"

	"github.com/muesli/reflow/ansi"
	"github.com/muesli/reflow/truncate"
)

const (
	ellipsis            = "..."
	minWidthForEllipsis = len(ellipsis) + 2
)

// DisplayWidth calculates what the rendered width of a string may be
func DisplayWidth(s string) int {
	return ansi.PrintableRuneWidth(s)
}

// Truncate shortens a string to fit the maximum display width
func Truncate(maxWidth int, s string) string {
	w := DisplayWidth(s)
	if w <= maxWidth {
		return s
	}

	tail := ""
	if maxWidth >= minWidthForEllipsis {
		tail = ellipsis
	}

	r := truncate.StringWithTail(s, uint(maxWidth), tail)
	if DisplayWidth(r) < maxWidth {
		r += " "
	}

	return r
}

// TruncateColumn replaces the first new line character with an ellipsis
// and shortens a string to fit the maximum display width
func TruncateColumn(maxWidth int, s string) string {
	if i := strings.IndexAny(s, "\r\n"); i >= 0 {
		s = s[:i] + ellipsis
	}
	return Truncate(maxWidth, s)
}
