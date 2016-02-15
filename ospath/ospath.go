// Package ospath provides utilities to get OS-specific directories.
package ospath

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// CacheDir tries to acquire an OS-specific app cache directory for the given importPath.
// Cache directory contains the app's cached data that can be regenerated as needed.
// Apps should never rely on the existence of cache files.
//
// It's guaranteed to be a unique directory for the importPath.
// Before returning the directory's path, CacheDir creates the directory if it
// doesn't already exist, so it can be used right away.
func CacheDir(importPath string) (string, error) {
	var home string
	if u, err := user.Current(); err != nil {
		home = os.Getenv("HOME")
		if home == "" {
			return "", err
		}
	} else {
		home = u.HomeDir
	}
	// TODO: Support Linux, Windows in analogous ways. Also support mobile devices (iOS, Android).
	//       Think about web? While HTML5 Local Storage could be used, it's not going to be compatible
	//       with filepaths; so maybe consider returning a webdav.FileSystem or so instead? Needs consideration.
	switch {
	case runtime.GOOS == "darwin" && runtime.GOARCH == "amd64":
		dir := filepath.Join(home, "Library", "Caches", filepath.FromSlash(importPath))
		if err := os.MkdirAll(dir, 0700); err != nil {
			return "", err
		}
		return dir, nil
	default:
		return "", fmt.Errorf("ospath.CacheDir not implemented for %s/%s", runtime.GOOS, runtime.GOARCH)
	}
}
