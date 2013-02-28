package gist4727543

import (
	. "gist.github.com/5052956.git"
	"strings"
)

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
	"gist.github.com/4670289.git":     "var _ = gist4670289.GoKeywords",
	"gist.github.com/5052956.git":     "var _ = gist5052956.GetGoFilePackageName",
}

func GetForcedUse(ImportPath string) string {
	return GetForcedUseRenamed(ImportPath, "")
}

func GetForcedUseRenamed(ImportPath, LocalPackageName string) string {
	if "" == LocalPackageName {
		return content[ImportPath]
	}

	filename := "./GoLand/src/" + ImportPath + "/gistfile1.go"
	packageName := GetGoFilePackageName(filename)
	if "." == LocalPackageName {
		return strings.Replace(content[ImportPath], packageName+".", "", 1)
	}
	return strings.Replace(content[ImportPath], packageName, LocalPackageName, 1)
}