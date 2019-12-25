package ttpr

import (
	"testing"
)

func BenchmarkWriteHtml(b *testing.B) {
	ranked := []string{"tuna in brine",
		"in soup, ginger crouton", "cover them in grease"}
	for i := 0; i < b.N; i++ {
		WriteHtml(ranked)
	}
}
