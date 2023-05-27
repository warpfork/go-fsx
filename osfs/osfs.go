/*
osfs contains support for FS interfaces backed by the real filesystem,
powered by the standard library's familiar `os` package.

In contrast to `os.DirFS`, this package contains very similar features,
but also includes support for additional go-fsx features,
such as FSSupportingWrite, FSSupportingReadlink, and FSSupportingMkSymlink.
*/
package osfs

import (
	"os"

	fs "github.com/warpfork/go-fsx"
)

var (
	_ fs.FSSupportingWrite     = dirFS("")
	_ fs.FSSupportingReadlink  = dirFS("")
	_ fs.FSSupportingMkSymlink = dirFS("")
)

func DirFS(dir string) fs.FS {
	return dirFS(dir)
}

type dirFS string

func (dir dirFS) Open(name string) (fs.File, error) {
	return os.DirFS(string(dir)).Open(name)
}

func (dir dirFS) Stat(name string) (fs.FileInfo, error) {
	return os.DirFS(string(dir)).(fs.FSSupportingStat).Stat(name)
}

func (dir dirFS) OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error) {
	return os.OpenFile(string(dir)+"/"+name, flag, perm)
}

func (dir dirFS) Mkdir(name string, perm fs.FileMode) error {
	return os.Mkdir(string(dir)+"/"+name, perm)
}

func (dir dirFS) Readlink(name string) (string, error) {
	return os.Readlink(string(dir) + "/" + name)
}

func (dir dirFS) Lstat(name string) (fs.FileInfo, error) {
	return os.Lstat(string(dir) + "/" + name)
}

func (dir dirFS) MkSymlink(name, target string) error {
	return os.Symlink(target, string(dir)+"/"+name)
}
