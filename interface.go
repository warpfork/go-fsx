package fsx

import (
	"io/fs"
	"os"
)

// We embed the interfaces from the fs package into our interfaces here.
// Although this is not technically necessary,
// it also doesn't hurt, and provides strong hints that aren't often wrong.

// FS is an interface to a filesystem.
// It is like stdlib's fs.FS, and includes it, and adds many additional features.
//
// fsx.FS also supports writing to a filesystem,
// by adding at least three additional methods: OpenFile, Mkdir, and Remove.
// It also adds additional features such as Readlink and LStat.
//
// The fsx.FS interface is not often seen in method signatures:
// all the functions in this package take still take a stdlib fs.FS as a parameter,
// and will attempt to feature-detect fsx.FS,
// erroring at runtime if the given FS does not support the required features.
type FS interface {
	fs.FS // fsx.FS is a superset of the read-only interface from the fs package.

	// OpenFile opens the named file, with some flags; it may be writable.
	//
	// The flags are the same as in os.OpenFile.
	OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error)

	Mkdir(name string, perm fs.FileMode) error

	// TODO Remove(name string) error

	// TODO Readlink(name string) (string, error)

	// TODO Lstat(name string) (fs.FileInfo, error)
}

type (
	File     = fs.File
	FileInfo = fs.FileInfo
	FileMode = fs.FileMode

	StatFS = fs.StatFS
)

// Flags to OpenFile wrapping those of the underlying system. Not all
// flags may be implemented on a given system.
//
// These are the same as the constants in the os package,
// and are replicated here only for convenience.
const (
	O_RDONLY int = os.O_RDONLY // open the file read-only.
	O_WRONLY int = os.O_WRONLY // open the file write-only.
	O_RDWR   int = os.O_RDWR   // open the file read-write.
	O_APPEND int = os.O_APPEND // append data to the file when writing.
	O_CREATE int = os.O_CREATE // create a new file if none exists.
	O_EXCL   int = os.O_EXCL   // used with O_CREATE, file must not exist.
	O_SYNC   int = os.O_SYNC   // open for synchronous I/O.
	O_TRUNC  int = os.O_TRUNC  // truncate regular writable file when opened.
)
