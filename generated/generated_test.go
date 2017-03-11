package generated_test

import (
	"path/filepath"
	"testing"

	"github.com/shurcooL/go/generated"
)

func TestParseFile(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		// Positive matches.
		{"positive.0.src", true},
		{"positive.1.src", true},
		{"positive.2.src", true},
		{"positive.3.src", true},
		{"positive.4.src", true},

		// Negative matches.
		{"negative.0.src", false},
		{"negative.1.src", false},
		{"negative.2.src", false},
		{"negative.3.src", false},
		{"negative.4.src", false},
		{"negative.5.src", false},
		{"negative.6.src", false},
		{"negative.7.src", false},
		{"negative.8.src", false},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			hasGeneratedComment, err := generated.ParseFile(filepath.Join("testdata", tc.name))
			if err != nil {
				t.Fatal(err)
			}
			if got, want := hasGeneratedComment, tc.want; got != want {
				t.Errorf("got hasGeneratedComment %v, want %v", got, want)
			}
		})
	}
}
