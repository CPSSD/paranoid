package ignore

import (
	"testing"
)

var flagtests = []struct {
	arg1 string
	arg2 string
	out  bool
}{
	{"file.html", "file.html", true},
	{"hello*.html", "helloWorld.html", true},
	{"*.html", "file.html", true},
	{"bin", "bin/asdf", true},
	{"bin/*", "bin/abc", true},
	{"bin/*/", "bin/abc/asdf", true},
	{"bin/*/*/", "bin/abc/asdf/", true},
	{"bin/**/asdf", "bin/asdf", true},
	{"bin/**/asdf", "bin/some/random/path/asdf", true},
	{"bin/file.html", "bin/asdf.html", false},
	{"!bin/file.html", "bin/file.html", false},
	{"!bin", "bin", false},
	{"!bin", "asdf", false},
}

func TestGlobbing(t *testing.T) {
	for _, tt := range flagtests {
		output := Glob(tt.arg1, tt.arg2)
		if output != tt.out {
			t.Errorf("Testcase failure", tt.arg1, tt.arg2)
		}
	}
}
