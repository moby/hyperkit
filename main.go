package gist4727543

import (
	"fmt"
	"strings"

	. "gist.github.com/5210270.git"
	. "gist.github.com/5504644.git"
)

// Generates an anonymous usage for the given import statement to avoid "imported and not used" errors
//
// e.g. `. "io/ioutil"` -> `var _ = NopCloser`
func GetForcedUseFromImport(Import string) (out string) {
	defer func() {
		e := recover()
		if nil != e {
			out = fmt.Sprint(e)
		}
	}()
	ImportParts := strings.Split(Import, " ")
	switch len(ImportParts) {
	case 1:
		return GetForcedUse(TrimQuotes(ImportParts[0]))
	case 2:
		return GetForcedUseRenamed(TrimQuotes(ImportParts[1]), ImportParts[0])
	default:
		panic("Invalid import string.")
	}
}

// Generates an anonymous usage of the package to avoid "imported and not used" errors
//
// e.g. `io/ioutil` -> `var _ = ioutil.NopCloser`
func GetForcedUse(ImportPath string) string {
	return GetForcedUseRenamed(ImportPath, "")
}

// Generates an anonymous usage of a renamed imported package
//
// e.g. `io/ioutil`, `RenamedPkg` -> `var _ = RenamedPkg.NopCloser`
func GetForcedUseRenamed(ImportPath, LocalPackageName string) string {
	dpkg, err := GetDocPackage(BuildPackageFromImportPath(ImportPath))
	if err != nil {
		return "Package not valid (doesn't exist or can't be built)."
	}

	// Uncomment only for testing purposes
	//dpkg.Funcs = dpkg.Funcs[0:0]
	//dpkg.Vars = dpkg.Vars[0:0]
	//dpkg.Consts = dpkg.Consts[0:0]
	//dpkg.Types = dpkg.Types[0:0]

	Prefix := "var _ = "
	var Usage string
	if len(dpkg.Funcs) > 0 {
		Usage = dpkg.Funcs[0].Name
	} else if len(dpkg.Vars) > 0 {
		Usage = dpkg.Vars[0].Names[0]
	} else if len(dpkg.Consts) > 0 {
		Usage = dpkg.Consts[0].Names[0]
	} else if len(dpkg.Types) > 0 {
		Usage = dpkg.Types[0].Name
		Prefix = "var _ "
	} else {
		return "Package doesn't have a single public func, var, const or type."
	}

	switch {
	case LocalPackageName == "":
		return Prefix + dpkg.Name + "." + Usage
	case LocalPackageName == ".":
		return Prefix + Usage
	default:
		return Prefix + LocalPackageName + "." + Usage
	}
}

func main() {
	println(GetForcedUse("io/ioutil"))
	println(GetForcedUseRenamed("io/ioutil", ""))
	println(GetForcedUseRenamed("io/ioutil", "RenamedPkg"))
	println(GetForcedUseRenamed("io/ioutil", "."))
	println()
	println(GetForcedUseFromImport(`gist.github.com/5210270.git`))
	println(GetForcedUseFromImport(`"gist.github.com/5210270.git"`))
	println(GetForcedUseFromImport(`. "gist.github.com/5210270.git"`))
	println(GetForcedUseFromImport(`bad bad bad`))
	println(GetForcedUseFromImport(`bad`))
}
