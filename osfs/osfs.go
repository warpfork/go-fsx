package osfs

import (
	"os"

	fs "github.com/warpfork/go-fsx"
)

var (
	_ fs.FSSupportingWrite = dirFS("")
)

func DirFS(dir string) fs.FS {
	return dirFS(dir)
}

type dirFS string

func (dir dirFS) Open(name string) (fs.File, error) {
	return os.DirFS(string(dir)).Open(name)
}

func (dir dirFS) Stat(name string) (fs.FileInfo, error) {
	return os.DirFS(string(dir)).(fs.StatFS).Stat(name)
}

func (dir dirFS) OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error) {
	return os.OpenFile(string(dir)+"/"+name, flag, perm)
}

func (dir dirFS) Mkdir(name string, perm fs.FileMode) error {
	return os.Mkdir(string(dir)+"/"+name, perm)
}
