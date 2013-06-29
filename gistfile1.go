package main

import (
	. "gist.github.com/5504644.git"
	"strings"
	. "gist.github.com/5210270.git"
	"fmt"
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
	if 1 == len(ImportParts) {
		return GetForcedUse(TrimQuotes(ImportParts[0]))
	} else if 2 == len(ImportParts) {
		return GetForcedUseRenamed(TrimQuotes(ImportParts[1]), ImportParts[0])
	}
	panic("Invalid import string.")
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
	dpkg := GetDocPackage(ImportPath)

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

	if "" == LocalPackageName {
		return Prefix + dpkg.Name + "." + Usage
	} else if "." == LocalPackageName {
		return Prefix + Usage
	} else {
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