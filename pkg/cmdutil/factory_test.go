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

package cmdutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_executable(t *testing.T) {
	testExe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}

	testExeName := filepath.Base(testExe)

	// Create 3 extra PATH entries that each contain an executable with the same name as the running test
	// process. The first is a symlink, but to an unrelated executable, the second is a symlink to our test
	// process and thus represents the result we want, and the third one is an unrelated executable.
	dir := t.TempDir()
	bin1 := filepath.Join(dir, "bin1")
	bin1Exe := filepath.Join(bin1, testExeName)
	bin2 := filepath.Join(dir, "bin2")
	bin2Exe := filepath.Join(bin2, testExeName)
	bin3 := filepath.Join(dir, "bin3")
	bin3Exe := filepath.Join(bin3, testExeName)

	if err := os.MkdirAll(bin1, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(bin2, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(bin3, 0755); err != nil {
		t.Fatal(err)
	}
	if f, err := os.OpenFile(bin3Exe, os.O_CREATE, 0755); err == nil {
		f.Close()
	} else {
		t.Fatal(err)
	}
	if err := os.Symlink(testExe, bin2Exe); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(bin3Exe, bin1Exe); err != nil {
		t.Fatal(err)
	}

	oldPath := os.Getenv("PATH")
	t.Cleanup(func() {
		os.Setenv("PATH", oldPath)
	})
	os.Setenv("PATH", strings.Join([]string{bin1, bin2, bin3, oldPath}, string(os.PathListSeparator)))

	if got := executable(""); got != bin2Exe {
		t.Errorf("executable() = %q, want %q", got, bin2Exe)
	}
}
