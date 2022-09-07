package fsx

import (
	"fmt"
	"io"
	stdfs "io/fs"
	"os"
	"syscall"
)

// Stat returns a FileInfo describing the named file from the file system.
//
// It is exactly per fs.Stat, and is in fact merely a wrapper,
// which we include in this package for the convenience of having all features available in one place.
func Stat(fsys FS, name string) (FileInfo, error) {
	return stdfs.Stat(fsys, name)
}

// ReadFile reads the named file from the file system fs and returns its contents.
//
// It is exactly per fs.ReadFile, and is in fact merely a wrapper,
// which we include in this package for the convenience of having all features available in one place.
func ReadFile(fsys FS, name string) ([]byte, error) {
	return stdfs.ReadFile(fsys, name)
}

// WriteFile is a shorthand for opening a file in write mode,
// either truncating or creating it as necessary,
// attempting to write the entire body of bytes given, and closing the file.
func WriteFile(fsys FS, name string, perm FileMode, body []byte) (err error) {
	f, err := OpenFile(fsys, name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	if f2, ok := f.(io.Writer); ok {
		_, err = f2.Write(body)
		if err1 := f.Close(); err1 != nil && err == nil {
			err = err1
		}
		return
	}
	return &stdfs.PathError{
		Op:   "WriteFile",
		Path: name,
		Err:  fmt.Errorf("filesystem type %T did not correctly support OpenFile for writable files", fsys),
	}
}

// OpenFile opens a file, with the provided flags (which are as per os.OpenFile flags)
// and perm bits, returning a File interface to handle it.
//
// If wanting a writable file, check for the io.Writer interface on the returned File.
// Although this File interface does not itself guarantee any write methods,
// they can be expected to be present when the flags to OpenFile ask for a writable file.
func OpenFile(fsys FS, name string, flag int, perm FileMode) (File, error) {
	if fsys2, ok := fsys.(FSSupportingWrite); ok {
		return fsys2.OpenFile(name, flag, perm)
	} else if flag == os.O_RDONLY {
		return fsys.Open(name)
	} else {
		return nil, &stdfs.PathError{
			Op:   "OpenFile",
			Path: name,
			Err:  fmt.Errorf("filesystem type %T does not support OpenFile", fsys),
		}
	}
}

func Mkdir(fsys FS, name string, perm FileMode) error {
	if fsys2, ok := fsys.(FSSupportingWrite); ok {
		return fsys2.Mkdir(name, perm)
	} else {
		return &stdfs.PathError{
			Op:   "Mkdir",
			Path: name,
			Err:  fmt.Errorf("filesystem type %T does not support Mkdir", fsys),
		}
	}
}

func MkdirAll(fsys FS, name string, perm FileMode) error {
	// The below code is considerably following from the function of the same name in the stdlib os package.
	// It is simplified in several places because it ignores the possibility of non-unix-style paths.

	fsys2, ok := fsys.(FSSupportingWrite)
	if !ok {
		return &stdfs.PathError{
			Op:   "Mkdir",
			Path: name,
			Err:  fmt.Errorf("filesystem type %T does not support Mkdir", fsys),
		}
	}

	// Fast path: if we can tell whether path is a directory or file, stop with success or error.
	dir, err := Stat(fsys, name)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return &stdfs.PathError{Op: "Mkdir", Path: name, Err: syscall.ENOTDIR}
	}

	// Slow path: make sure parent exists and then call Mkdir for path.
	i := len(name)
	for i > 0 && name[i-1] == '/' { // Skip trailing path separator.
		i--
	}

	j := i
	for j > 0 && name[j-1] == '/' { // Scan backward over element.
		j--
	}

	if j > 1 {
		// Create parent.
		err = MkdirAll(fsys, name[:j-1], perm)
		if err != nil {
			return err
		}
	}

	// Parent now exists; invoke Mkdir and use its result.
	err = fsys2.Mkdir(name, perm)
	if err != nil {
		// Handle arguments like "foo/." by double-checking that directory doesn't exist.
		dir, err1 := Stat(fsys, name) // Note!  This is Lstat in the original os package code; however, this package (as yet) does not support symlinks and disregards the concept, so we have only Stat here.
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}
	return nil

}
