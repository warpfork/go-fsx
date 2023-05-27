/*
Not yet implemented.
*/
package memfs

import (
	"testing/fstest"

	fs "github.com/warpfork/go-fsx"
)

func MemFS() fs.FS {
	return make(memFS)
}

type memFS fstest.MapFS

func (mem memFS) Open(name string) (fs.File, error) {
	panic("todo")
}

func (mem memFS) Stat(name string) (fs.FileInfo, error) {
	return ((fstest.MapFS)(mem)).Stat(name)
}

func (mem memFS) OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error) {
	panic("todo")
}

func (mem memFS) Mkdir(name string, perm fs.FileMode) error {
	panic("todo")
}

// ... much todo.  I initially hoped I could reuse `*fstest.MapFile` and `*fstest.openMapFile`, but on deeper look, it's quite unclear if it would be possible to make it support writing.
// the directory synthesis done by fstest also seems undesirable for our purposes here.
