package gist4727543

var content = map[string]string{
	"fmt":       "var _ = fmt.Printf",
	"reflect":   "var _ = reflect.TypeOf",
	"io/ioutil": "var _ = ioutil.ReadFile",
	"os/exec":   "var _ = exec.Command",
	"net/http":  "var _ = http.Get",
	"go/ast":    "var _ ast.Ident",
	"github.com/davecgh/go-spew/spew": "var _ = spew.Dump",
}

func GetForcedUse(ImportPath string) string {
	return content[ImportPath]
}