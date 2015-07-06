// +build ignore

package gopherjs_http

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	pathpkg "path"
	"strings"
	"time"

	"sourcegraph.com/sourcegraph/go-vcs/vcs/util"
)

func NewTestFs(fs http.FileSystem) http.FileSystem {
	return &testFs{
		fs: fs,
	}
}

type testFs struct {
	fs http.FileSystem
}

func (v *testFs) Open(path string) (http.File, error) {
	// HACK.
	if path == "/" {
		fi := &util.FileInfo{
			//Name_:    pathpkg.Base("/script,edit.js"),
			Name_:    pathpkg.Base("/script.js"),
			Mode_:    os.FileMode(0),
			Size_:    int64(-1),
			ModTime_: time.Now(),
			Sys_:     nil,
		}

		return &httpDir{
			path: path,
			FileInfo: &util.FileInfo{
				Name_:    "/",
				Mode_:    os.FileMode(os.ModeDir),
				Size_:    int64(0),
				ModTime_: time.Time{}, //time.Now(),
				Sys_:     nil,
			},
			entries: []os.FileInfo{fi},
		}, nil
	}

	/*if path.Ext(name) == ".txt" {
		b, err := vfs.ReadFile(v.fs, name)
		if err != nil {
			return nil, err
		}
		return util.NopCloser{strings.NewReader(gist5423254.Reverse(string(b)))}, nil
	}*/

	if pathpkg.Ext(path) == ".js" {
		var f File

		name := pathpkg.Base(path)
		nameWithoutExt := name[:len(name)-len(".js")]
		sourcesWithoutExt := strings.Split(nameWithoutExt, ",")

		var names []string
		var goReaders []io.Reader
		var goClosers []io.Closer
		for _, sourceWithoutExt := range sourcesWithoutExt {
			file, err := v.fs.Open("/" + sourceWithoutExt + "/main.go") // TODO.
			if err != nil {
				return nil, err
			}
			names = append(names, sourceWithoutExt+".go")
			goReaders = append(goReaders, file)
			goClosers = append(goClosers, file)
			//f.dependencies = append(f.dependencies, "/assets/"+sourceWithoutExt+".go")
		}

		fmt.Println("REBUILDING SOURCE for:", name)
		//debug.PrintStack()
		content := []byte(handleJsError(goReadersToJs(names, goReaders)))
		f.Reader = bytes.NewReader(content)

		for _, closer := range goClosers {
			closer.Close()
		}

		f.path = name
		f.FileInfo = &util.FileInfo{
			Name_:    pathpkg.Base(name),
			Mode_:    os.FileMode(0),
			Size_:    int64(len(content)),
			ModTime_: time.Now(),
			Sys_:     nil,
		}
		return &f, nil
	}

	//return v.fs.Open(name)
	return nil, fmt.Errorf("no %q file", path)
}

type File struct {
	path string
	*util.FileInfo
	//content      []byte
	*bytes.Reader
	//dependencies []string
}

func (f *File) Stat() (os.FileInfo, error) {
	return f.FileInfo, nil
}

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	panic("Readdir in file")
}

func (_ *File) Close() error { return nil }

// httpDir implements http.File for a directory in a FileSystem.
type httpDir struct {
	path string
	*util.FileInfo
	entries []os.FileInfo
}

func (_ *httpDir) Close() error { return nil }

func (d *httpDir) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.path)
}

func (d *httpDir) Seek(offset int64, whence int) (int64, error) {
	return 0, fmt.Errorf("cannot Seek in directory %s", d.path)
}

func (d *httpDir) Stat() (os.FileInfo, error) {
	return d.FileInfo, nil
}

func (d *httpDir) Readdir(count int) ([]os.FileInfo, error) {
	if count != 0 {
		log.Panicln("httpDir.Readdir count unsupported value:", count)
	}

	return d.entries, nil
}

/*func (f *File) Stat(name string) (os.FileInfo, error) {
	/*if path.Ext(name) == ".txt" {
		return &util.FileInfo{
			Name_:    path.Base(name),
			Mode_:    os.FileMode(0),
			Size_:    3,
			ModTime_: time.Time{},
			Sys_:     nil,
		}, nil
	}* /

	f, ok := v.cache[name]

	/*if path.Ext(name) == ".txt" {
		f.content = []byte(name)
		f.FileInfo = &util.FileInfo{
			Name_:    path.Base(name),
			Mode_:    os.FileMode(0),
			Size_:    int64(len(f.content)),
			ModTime_: time.Now(),
			Sys_:     nil,
		}
		v.mu.Lock()
		v.cache[name] = f
		v.mu.Unlock()
		return f.FileInfo, nil
	}* /

	//return v.fs.Stat(name)
	return nil, fmt.Errorf("no %q file", name)
}*/

//func (v *testFs) ReadDir(path string) ([]os.FileInfo, error) { return nil, nil } /*return v.fs.ReadDir(path)*/
/*func (v *testFs) ReadDir(path string) ([]os.FileInfo, error) {
	if path == "/" {
		fi, err := v.Stat("/script,edit.js")
		if err != nil {
			return nil, err
		}

		return []os.FileInfo{
			fi,
		}, nil
	}
	return nil, nil
}*/

func (v *testFs) String() string { return "testfs" }

/*func NewTestHttpFs(fs vfs.FileSystem) vfs.FileSystem {
	return &testFs{
		fs:    fs,
		cache: make(map[string]File),
	}
}
type testHttpFs struct {
	fs vfs.FileSystem

	mu    sync.RWMutex
	cache map[string]File
}
type File2 struct {
	*util.FileInfo
	content []byte
}*/
