// Package gist4727543 generates an anonymous usage of the package for given import path
// to avoid "imported and not used" errors.
//
// This is largely unneeded now that goimports exists.
package gist4727543

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shurcooL/go/gists/gist5504644"
)

// tryUnquote returns the unquoted string, or the original string if unquoting fails.
func tryUnquote(s string) string {
	if t, err := strconv.Unquote(s); err == nil {
		return t
	}
	return s
}

// GetForcedUseFromImport generates an anonymous usage for the given import spec to avoid "imported and not used" errors.
//
// E.g., `io/ioutil` -> `var _ = ioutil.NopCloser`,
//
// `renamed "io/ioutil"` -> `var _ = renamed.NopCloser`,
//
// `. "io/ioutil"` -> `var _ = NopCloser`.
func GetForcedUseFromImport(importSpec string) (out string) {
	switch parts := strings.Split(importSpec, " "); len(parts) {
	case 1:
		return GetForcedUse(tryUnquote(parts[0]))
	case 2:
		return GetForcedUseRenamed(tryUnquote(parts[1]), parts[0])
	default:
		return "Invalid import string."
	}
}

// GetForcedUse generates an anonymous usage of the package to avoid "imported and not used" errors
//
// E.g., `io/ioutil` -> `var _ = ioutil.NopCloser`.
func GetForcedUse(importPath string) string {
	return GetForcedUseRenamed(importPath, "")
}

// GetForcedUseRenamed generates an anonymous usage of a renamed imported package.
//
// E.g., `io/ioutil`, `RenamedPkg` -> `var _ = RenamedPkg.NopCloser`.
func GetForcedUseRenamed(importPath, localPackageName string) string {
	dpkg, err := gist5504644.GetDocPackage(gist5504644.BuildPackageFromImportPath(importPath))
	if err != nil {
		return fmt.Sprintf("Package %q not valid (doesn't exist or can't be built).", importPath)
	}

	// Uncomment only for testing purposes.
	//dpkg.Funcs = dpkg.Funcs[0:0]
	//dpkg.Vars = dpkg.Vars[0:0]
	//dpkg.Consts = dpkg.Consts[0:0]
	//dpkg.Types = dpkg.Types[0:0]

	prefix := "var _ = "
	var usage string
	if len(dpkg.Funcs) > 0 {
		usage = dpkg.Funcs[0].Name
	} else if len(dpkg.Vars) > 0 {
		usage = dpkg.Vars[0].Names[0]
	} else if len(dpkg.Consts) > 0 {
		usage = dpkg.Consts[0].Names[0]
	} else if len(dpkg.Types) > 0 {
		usage = dpkg.Types[0].Name
		prefix = "var _ "
	} else {
		return "Package doesn't have a single public func, var, const or type."
	}

	switch {
	case localPackageName == "":
		return prefix + dpkg.Name + "." + usage
	case localPackageName == ".":
		return prefix + usage
	default:
		return prefix + localPackageName + "." + usage
	}
}
