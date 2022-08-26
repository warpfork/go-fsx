package rwfs

import (
	"io/fs"
	"os"
)

// We embed the interfaces from the fs package into our interfaces here.
// Although this is not technically necessary,
// it also doesn't hurt, and provides strong hints that aren't often wrong.

// RWFS is an fs.FS which also supports writing to a filesystem,
// by adding at least three additional methods: OpenFile, Mkdir, and Remove.
//
// The RWFS interface is not often seen in method signatures:
// all the functions in this package take fs.FS as a parameter,
// and will attempt to feature-detect RWFS,
// erroring at runtime if the given filesystem does not support the required features.
type RWFS interface {
	fs.FS // RWFS is a superset of the read-only interface from the fs package.

	// OpenFile opens the named file, with some flags; it may be writable.
	//
	// The flags are the same as in os.OpenFile.
	OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error)

	Mkdir(name string, perm fs.FileMode) error

	// TODO Remove(name string) error
}

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
