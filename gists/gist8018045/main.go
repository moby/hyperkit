// Package gist8018045 provides funcs to get a list of all local Go packages in GOPATH workspaces and GOROOT.
package gist8018045

import (
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/shurcooL/go/gists/gist5504644"
	"github.com/shurcooL/go/gists/gist7480523"
)

// GetGoPackages gets all local Go packages (from GOROOT and all GOPATH workspaces).
func GetGoPackages(out chan<- *gist7480523.GoPackage) {
	for _, root := range build.Default.SrcDirs() {
		_ = filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				log.Printf("can't stat file %s: %v\n", path, err)
				return nil
			}
			switch {
			case !fi.IsDir():
				return nil
			case path == root:
				return nil
			case strings.HasPrefix(fi.Name(), ".") || strings.HasPrefix(fi.Name(), "_") || fi.Name() == "testdata":
				return filepath.SkipDir
			default:
				importPath, err := filepath.Rel(root, path)
				if err != nil {
					return nil
				}
				// Prune search if we encounter any of these import paths.
				switch importPath {
				case "builtin":
					return nil
				}
				if goPackage := gist7480523.GoPackageFromImportPath(importPath); goPackage != nil {
					out <- goPackage
				}
				return nil
			}
		})
	}
	close(out)
}

// GetGopathGoPackages gets Go packages in all GOPATH workspaces.
func GetGopathGoPackages(out chan<- *gist7480523.GoPackage) {
	gopathEntries := filepath.SplitList(build.Default.GOPATH)
	for _, gopathEntry := range gopathEntries {
		root := filepath.Join(gopathEntry, "src")
		if !isDir(root) {
			continue
		}
		_ = filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				log.Printf("can't stat file %s: %v\n", path, err)
				return nil
			}
			if !fi.IsDir() {
				return nil
			}
			if strings.HasPrefix(fi.Name(), ".") || strings.HasPrefix(fi.Name(), "_") || fi.Name() == "testdata" {
				return filepath.SkipDir
			}
			importPath, err := filepath.Rel(root, path)
			if err != nil {
				return nil
			}
			importPathFound := gist7480523.NewImportPathFound(importPath, gopathEntry)
			if goPackage := gist7480523.GoPackageFromImportPathFound(importPathFound); goPackage != nil {
				out <- goPackage
			}
			return nil
		})
	}
	close(out)
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

// rec is the recursive body of getGoPackagesA.
func rec(out chan<- gist7480523.ImportPathFound, importPathFound gist7480523.ImportPathFound) {
	if goPackage := gist7480523.GoPackageFromImportPathFound(importPathFound); goPackage != nil {
		out <- importPathFound
	}

	entries, err := ioutil.ReadDir(importPathFound.FullPath())
	if err == nil {
		for _, v := range entries {
			if v.IsDir() && !strings.HasPrefix(v.Name(), ".") && !strings.HasPrefix(v.Name(), "_") || v.Name() == "testdata" {
				rec(out, gist7480523.NewImportPathFound(filepath.Join(importPathFound.ImportPath(), v.Name()), importPathFound.GopathEntry()))
			}
		}
	}
}

func getGoPackagesA(out chan<- gist7480523.ImportPathFound) {
	gopathEntries := filepath.SplitList(build.Default.GOPATH)
	for _, gopathEntry := range gopathEntries {
		rec(out, gist7480523.NewImportPathFound(".", gopathEntry))
	}
	close(out)
}

func getGoPackagesB(out chan<- gist7480523.ImportPathFound) {
	gopathEntries := filepath.SplitList(build.Default.GOPATH)
	for _, gopathEntry := range gopathEntries {
		root := filepath.Join(gopathEntry, "src")
		if !isDir(root) {
			continue
		}
		_ = filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				log.Printf("can't stat file %s: %v\n", path, err)
				return nil
			}
			if !fi.IsDir() {
				return nil
			}
			if strings.HasPrefix(fi.Name(), ".") || strings.HasPrefix(fi.Name(), "_") || fi.Name() == "testdata" {
				return filepath.SkipDir
			}
			importPath, err := filepath.Rel(root, path)
			if err != nil {
				return nil
			}
			importPathFound := gist7480523.NewImportPathFound(importPath, gopathEntry)
			if goPackage := gist7480523.GoPackageFromImportPathFound(importPathFound); goPackage != nil {
				out <- importPathFound
			}
			return nil
		})
	}
	close(out)
}

func getGoPackagesC(out chan<- gist7480523.ImportPathFound) {
	gopathEntries := filepath.SplitList(build.Default.GOPATH)
	for _, gopathEntry := range gopathEntries {
		root := filepath.Join(gopathEntry, "src")
		if !isDir(root) {
			continue
		}
		_ = filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				log.Printf("can't stat file %s: %v\n", path, err)
				return nil
			}
			if !fi.IsDir() {
				return nil
			}
			if strings.HasPrefix(fi.Name(), ".") || strings.HasPrefix(fi.Name(), "_") || fi.Name() == "testdata" {
				return filepath.SkipDir
			}
			bpkg, err := gist5504644.BuildPackageFromSrcDir(path)
			if err != nil {
				return nil
			}
			/*if bpkg.Goroot {
				return nil
			}*/
			out <- gist7480523.NewImportPathFound(bpkg.ImportPath, bpkg.Root)
			return nil
		})
	}
	close(out)
}
