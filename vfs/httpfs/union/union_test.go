package union_test

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/shurcooL/go/vfs/httpfs/union"
	"github.com/shurcooL/go/vfs/httpfs/vfsutil"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
)

func walk(path string, fi os.FileInfo, err error) error {
	if err != nil {
		log.Printf("can't stat file %s: %v\n", path, err)
		return nil
	}
	fmt.Println(path)
	return nil
}

func Example() {
	fs0 := httpfs.New(mapfs.New(map[string]string{
		"zzz-last-file.txt":   "It should be visited last.",
		"a-file.txt":          "It has stuff.",
		"another-file.txt":    "Also stuff.",
		"folderA/entry-A.txt": "Alpha.",
		"folderA/entry-B.txt": "Beta.",
	}))
	fs1 := httpfs.New(mapfs.New(map[string]string{
		"sample-file.txt":                "This file compresses well. Blaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaah!",
		"not-worth-compressing-file.txt": "Its normal contents are here.",
		"folderA/file1.txt":              "Stuff 1.",
		"folderA/file2.txt":              "Stuff 2.",
		"folderB/folderC/file3.txt":      "Stuff C-3.",
	}))

	fs := union.New(map[string]http.FileSystem{
		"/fs0": fs0,
		"/fs1": fs1,
	})

	err := vfsutil.Walk(fs, "/", walk)
	if err != nil {
		panic(err)
	}

	// Output:
	// /
	// /fs0
	// /fs0/a-file.txt
	// /fs0/another-file.txt
	// /fs0/folderA
	// /fs0/folderA/entry-A.txt
	// /fs0/folderA/entry-B.txt
	// /fs0/zzz-last-file.txt
	// /fs1
	// /fs1/folderA
	// /fs1/folderA/file1.txt
	// /fs1/folderA/file2.txt
	// /fs1/folderB
	// /fs1/folderB/folderC
	// /fs1/folderB/folderC/file3.txt
	// /fs1/not-worth-compressing-file.txt
	// /fs1/sample-file.txt
}

func ExampleEmpty() {
	empty := union.New(nil)

	err := vfsutil.Walk(empty, "/", walk)
	if err != nil {
		panic(err)
	}

	// Output:
	// /
}

func ExampleNotExist() {
	empty := union.New(nil)

	_, err := empty.Open("/does-not-exist")
	fmt.Println("os.IsNotExist:", os.IsNotExist(err))
	fmt.Println(err)

	// Output:
	// os.IsNotExist: true
	// open /does-not-exist: file does not exist
}
