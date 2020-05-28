package json

import "testing"

func BenchmarkJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generate()
	}
}