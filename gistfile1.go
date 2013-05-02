package main

import (
	. "gist.github.com/5504644.git"
)

func GetForcedUse(ImportPath string) string {
	return GetForcedUseRenamed(ImportPath, "")
}

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
	return "" // TODO: Remove in go1.1
}

func main() {
	println(GetForcedUse("io/ioutil"))
	println(GetForcedUseRenamed("io/ioutil", ""))
	println(GetForcedUseRenamed("io/ioutil", "RenamedPkg"))
	println(GetForcedUseRenamed("io/ioutil", "."))
}