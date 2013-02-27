package gist4727543

var content = map[string]string{
	"fmt":       "var _ = fmt.Printf",
	"reflect":   "var _ = reflect.TypeOf",
	"io/ioutil": "var _ = ioutil.ReadFile",
	"os/exec":   "var _ = exec.Command",
	"net/http":  "var _ = http.Get",
	"go/ast":    "var _ ast.Ident",
	"github.com/davecgh/go-spew/spew": "var _ = spew.Dump",
	"gist.github.com/4668739.git":     "var _ = gist4668739.HttpGet",
	"gist.github.com/4727543.git":     "var _ = gist4727543.GetForcedUse",
}

func GetForcedUse(ImportPath string) string {
	return content[ImportPath]
}