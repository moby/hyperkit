package vfs_util

import (
	"os"
	"path/filepath"
	"sort"

	"code.google.com/p/go.tools/godoc/vfs"
)

func Walk(fs vfs.FileSystem, root string, walkFn filepath.WalkFunc) error {
	info, err := fs.Lstat(root)
	if err != nil {
		return walkFn(root, nil, err)
	}
	return walk(fs, root, info, walkFn)
}

// readDirNames reads the directory named by dirname and returns
// a sorted list of directory entries.
func readDirNames(fs vfs.FileSystem, dirname string) ([]string, error) {
	fis, err := fs.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, fi := range fis {
		names = append(names, fi.Name())
	}
	sort.Strings(names)
	return names, nil
}

// walk recursively descends path, calling w.
func walk(fs vfs.FileSystem, path string, info os.FileInfo, walkFn filepath.WalkFunc) error {
	err := walkFn(path, info, nil)
	if err != nil {
		if info.IsDir() && err == filepath.SkipDir {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return nil
	}

	names, err := readDirNames(fs, path)
	if err != nil {
		return walkFn(path, info, err)
	}

	for _, name := range names {
		filename := filepath.Join(path, name)
		fileInfo, err := fs.Lstat(filename)
		if err != nil {
			if err := walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir {
				return err
			}
		} else {
			err = walk(fs, filename, fileInfo, walkFn)
			if err != nil {
				if !fileInfo.IsDir() || err != filepath.SkipDir {
					return err
				}
			}
		}
	}
	return nil
}
