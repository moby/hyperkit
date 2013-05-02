package gist4727543

import (
	. "gist.github.com/5504644.git"
)

func GetForcedUse(ImportPath string) string {
	return GetForcedUseRenamed(ImportPath, "")
}

func GetForcedUseRenamed(ImportPath, LocalPackageName string) string {
	dpkg := GetDocPackage(ImportPath)

	// TODO: What if Funcs is empty? Fall back to Vars, Consts, Types

	if "" == LocalPackageName {
		return "var _ = " + dpkg.Name + "." + dpkg.Funcs[0].Name
	} else if "." == LocalPackageName {
		return "var _ = " + dpkg.Funcs[0].Name
	} else {
		return "var _ = " + LocalPackageName + "." + dpkg.Funcs[0].Name
	}
	return "" // TODO: Remove in go1.1
}