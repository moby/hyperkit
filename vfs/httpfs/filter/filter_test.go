package filter_test

import (
	"fmt"
	"io"
	"log"
	"os"
	pathpkg "path"

	"github.com/shurcooL/go/vfs/httpfs/filter"
	"github.com/shurcooL/go/vfs/httpfs/vfsutil"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
)

func Example() {
	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Printf("can't stat file %s: %v\n", path, err)
			return nil
		}
		fmt.Println(path)
		return nil
	}

	fs := httpfs.New(mapfs.New(map[string]string{
		"zzz-last-file.txt":                "It should be visited last.",
		"a-file.txt":                       "It has stuff.",
		"another-file.txt":                 "Also stuff.",
		"some-file.html":                   "<html>and stuff</html>",
		"folderA/entry-A.txt":              "Alpha.",
		"folderA/entry-B.txt":              "Beta.",
		"folderA/main.go":                  "package main\n",
		"folderA/folder-to-skip/many.txt":  "Entire folder can be skipped.",
		"folderA/folder-to-skip/files.txt": "Entire folder can be skipped.",
		"folder-to-skip":                   "This is a file, not a folder, and shouldn't be skipped.",
	}))

	ignore := func(fi os.FileInfo, _ string) bool {
		return pathpkg.Ext(fi.Name()) == ".go" || pathpkg.Ext(fi.Name()) == ".html" ||
			(fi.IsDir() && fi.Name() == "folder-to-skip")
	}

	fs = filter.NewIgnore(fs, ignore)

	err := vfsutil.Walk(fs, "/", walk)
	if err != nil {
		panic(err)
	}

	fmt.Println()

	// This file should be filtered out, even if accessed directly.
	_, err = fs.Open("/folderA/main.go")
	fmt.Println("os.IsNotExist(err):", os.IsNotExist(err))
	fmt.Println(err)

	fmt.Println()

	// This folder should be filtered out, even if accessed directly.
	_, err = fs.Open("/folderA/folder-to-skip")
	fmt.Println("os.IsNotExist(err):", os.IsNotExist(err))
	fmt.Println(err)

	fmt.Println()

	// This file should not be filtered out.
	f, err := fs.Open("/folder-to-skip")
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, f)
	f.Close()

	// Output:
	// /
	// /a-file.txt
	// /another-file.txt
	// /folder-to-skip
	// /folderA
	// /folderA/entry-A.txt
	// /folderA/entry-B.txt
	// /zzz-last-file.txt
	//
	// os.IsNotExist(err): true
	// open /folderA/main.go: file does not exist
	//
	// os.IsNotExist(err): true
	// open /folderA/folder-to-skip: file does not exist
	//
	// This is a file, not a folder, and shouldn't be skipped.
}
