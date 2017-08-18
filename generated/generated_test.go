package generated_test

import (
	"bytes"
	"io/ioutil"
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
		{"positive.5.src", true},
		{"positive.6.src", true},
		{"positive.7.src", true},
		{"positive.8.src", true},
		{"positive.9.src", true},
		{"positive.10.src", true},
		{"positive.11.src", true},
		{"positive.12.src", true},

		// Negative matches.
		{"negative.0.src", false},
		{"negative.1.src", false},
		{"negative.2.src", false},
		{"negative.3.src", false},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			hasGeneratedComment, err := generated.ParseFile(filepath.Join("testdata", tc.name))
			if err != nil {
				t.Error(err)
				return
			}
			if got, want := hasGeneratedComment, tc.want; got != want {
				t.Errorf("got hasGeneratedComment %v, want %v", got, want)
			}
		})

		// On Windows, a file that hasn't been gofmt'ed can have \r\n line endings.
		// Though rare and unusual, it's still a valid .go file and needs to be supported.
		t.Run(tc.name+` \r\n line ending version`, func(t *testing.T) {
			// Replace all "\n" line endings with "\r\n".
			b, err := ioutil.ReadFile(filepath.Join("testdata", tc.name))
			if err != nil {
				t.Error(err)
				return
			}
			b = bytes.Replace(b, []byte("\n"), []byte("\r\n"), -1)

			hasGeneratedComment, err := generated.Parse(bytes.NewReader(b))
			if err != nil {
				t.Error(err)
				return
			}
			if got, want := hasGeneratedComment, tc.want; got != want {
				t.Errorf("got hasGeneratedComment %v, want %v", got, want)
			}
		})
	}
}

func TestParseFileError(t *testing.T) {
	_, err := generated.ParseFile(filepath.Join("testdata", "doesnotexist"))
	if err == nil {
		t.Fatal("got nil error, want non-nil")
	}
}

func BenchmarkParseFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generated.ParseFile(filepath.Join("testdata", "positive.6.src"))
		generated.ParseFile(filepath.Join("testdata", "negative.3.src"))
	}
}
