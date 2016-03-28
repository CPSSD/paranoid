package benchmark

import (
	. "github.com/cpssd/paranoid/pfsd/pfi/glob"
	"testing"
)

func BenchmarkBasic(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Glob("bin", "bin")
	}
}

func BenchmarkWildcard(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Glob("bin/*/index.html", "bin/test/index.html")
	}
}

func BenchmarkDualWildcard(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Glob("bin/**/index.html", "bin/some/large/nested/path/index.html")
	}
}

func BenchmarkNegation(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Glob("!bin", "bin")
	}
}

func BenchmarkComplex(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Glob("!bin/**/index.html", "bin/some/large/nested/path/index.html")
	}
}
