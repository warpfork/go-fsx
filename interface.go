package fsx

import (
	"io/fs"
	"os"
)

type (
	FS        = fs.FS
	File      = fs.File
	FileInfo  = fs.FileInfo
	FileMode  = fs.FileMode
	DirEntry  = fs.DirEntry
	PathError = fs.PathError

	WalkDirFunc = fs.WalkDirFunc

	FSSupportingStat = fs.StatFS
)

// FSSupportingWrite extends FS to include functions which allow writing to the filesystem.
//
// The three critical additional methods are: OpenFile, Mkdir, and Remove.
// (For ability to write symlinks, see FSSupportingMkSymlink.)
//
// This interface is not often seen in method signatures:
// all the functions in this package take still take a plain FS as a parameter,
// and will attempt to feature-detect this interface if it is required,
// erroring at runtime if the given FS does not support the required features.
// (This keeps the amount of casting the programmer needs to do at a minimum.)
type FSSupportingWrite interface {
	fs.FS

	// OpenFile opens the named file, with some flags; it may be writable.
	//
	// The flags are the same as in os.OpenFile.
	OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error)

	Mkdir(name string, perm fs.FileMode) error

	// TODO Remove(name string) error
}

const (
	// Flags to OpenFile wrapping those of the underlying system. Not all
	// flags may be implemented on a given system.
	//
	// These are the same as the constants in the os package,
	// and are replicated here only for convenience.

	O_RDONLY int = os.O_RDONLY // open the file read-only.
	O_WRONLY int = os.O_WRONLY // open the file write-only.
	O_RDWR   int = os.O_RDWR   // open the file read-write.
	O_APPEND int = os.O_APPEND // append data to the file when writing.
	O_CREATE int = os.O_CREATE // create a new file if none exists.
	O_EXCL   int = os.O_EXCL   // used with O_CREATE, file must not exist.
	O_SYNC   int = os.O_SYNC   // open for synchronous I/O.
	O_TRUNC  int = os.O_TRUNC  // truncate regular writable file when opened.
)

const (
	// Bits used to indicate modes and permissions of files and directories.
	//
	// These are the same as the constants in the fs package,
	// and are replicated here only for convenience.

	ModeDir        = fs.ModeDir        // d: is a directory
	ModeAppend     = fs.ModeAppend     // a: append-only
	ModeExclusive  = fs.ModeExclusive  // l: exclusive use
	ModeTemporary  = fs.ModeTemporary  // T: temporary file; Plan 9 only
	ModeSymlink    = fs.ModeSymlink    // L: symbolic link
	ModeDevice     = fs.ModeDevice     // D: device file
	ModeNamedPipe  = fs.ModeNamedPipe  // p: named pipe (FIFO)
	ModeSocket     = fs.ModeSocket     // S: Unix domain socket
	ModeSetuid     = fs.ModeSetuid     // u: setuid
	ModeSetgid     = fs.ModeSetgid     // g: setgid
	ModeCharDevice = fs.ModeCharDevice // c: Unix character device, when ModeDevice is set
	ModeSticky     = fs.ModeSticky     // t: sticky
	ModeIrregular  = fs.ModeIrregular  // ?: non-regular file; nothing else is known about this file

	// Mask for the type bits. For regular files, none will be set.
	ModeType = fs.ModeType

	ModePerm = fs.ModePerm // Unix permission bits, 0o777
)

type FSSupportingReadlink interface {
	fs.FS

	Readlink(name string) (string, error)

	Lstat(name string) (fs.FileInfo, error)
}

type FSSupportingMkSymlink interface {
	// MkSymlink creates a symlink on the filesystem.
	// The name paramter is where the symlink will be created;
	// the target becomes the body of the symlink.
	//
	// Note that a symlink is truly a string.
	// Typically, the target refers to some other file or path,
	// but there is no guarantee of this.
	// The target string will also not be normalized in any way
	// (e.g. if you have a sub-FS, and the target starts with "/",
	// there is no expectation that the FS implementation do anything
	// to normalize or process that target string).
	MkSymlink(name, target string) error
}

// FUTURE: FSSupportingChown, etc.
