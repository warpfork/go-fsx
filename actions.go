package fsx

import (
	"fmt"
	"io"
	stdfs "io/fs"
	"os"
	"syscall"
)

// WalkDir is exactly as per stdlib io/fs.WalkDir:
// it reads and traverses directories recursively, visiting contents with a callback.
func WalkDir(fsys FS, root string, fn WalkDirFunc) error {
	return stdfs.WalkDir(fsys, root, fn)
}

// ReadDir is exactly as per stdlib io/fs.ReadDir:
// it returns a list of directory entries sorted by filename.
func ReadDir(fsys FS, name string) ([]DirEntry, error) {
	return stdfs.ReadDir(fsys, name)
}

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
	return &PathError{
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
		return nil, &PathError{
			Op:   "OpenFile",
			Path: name,
			Err:  fmt.Errorf("filesystem type %T does not support OpenFile", fsys),
		}
	}
}

// Mkdir creates a directory.
//
// This only works if the FS implements FSSupportingWrite; otherwise, an error will be returned.
func Mkdir(fsys FS, name string, perm FileMode) error {
	if fsys2, ok := fsys.(FSSupportingWrite); ok {
		return fsys2.Mkdir(name, perm)
	} else {
		return &PathError{
			Op:   "Mkdir",
			Path: name,
			Err:  fmt.Errorf("filesystem type %T does not support Mkdir", fsys),
		}
	}
}

// MkdirAll creates a directory, and any parent directories that do not yet exist.
//
// This only works if the FS implements FSSupportingWrite; otherwise, an error will be returned.
func MkdirAll(fsys FS, name string, perm FileMode) error {
	// The below code is considerably following from the function of the same name in the stdlib os package.
	// It is simplified in several places because it ignores the possibility of non-unix-style paths.

	fsys2, ok := fsys.(FSSupportingWrite)
	if !ok {
		return &PathError{
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
		return &PathError{Op: "Mkdir", Path: name, Err: syscall.ENOTDIR}
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
		dir, err1 := Lstat(fsys, name)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}
	return nil

}

// Readlink returns the target of a symlink, as a string.
//
// Note that while symlinks typically point at other files or directories, there is no guarantee of that, nor that the target actually exists.
// A symlink target is just a string.
//
// If the given FS implementation doesn't support FSSupportingReadlink, an error will be returned.
//
// Implementers of FSSupportingReadlink should also error if Readlink is called on a non-symlink.
func Readlink(fsys FS, name string) (string, error) {
	if fsys2, ok := fsys.(FSSupportingReadlink); ok {
		return fsys2.Readlink(name)
	} else {
		return "", &PathError{
			Op:   "Readlink",
			Path: name,
			Err:  fmt.Errorf("filesystem type %T does not support Readlink", fsys),
		}
	}
}

// Lstat returns FileInfo for a file, or for a symlink without traversing the link, if called on one.
//
// If the given FS implementation doesn't support FSSupportingReadlink, then this function assumes no symlinks,
// and falls back to regular Stat.
// If the given FS implementation also does not support FSSupportingStat, then an error is returned.
func Lstat(fsys FS, name string) (FileInfo, error) {
	if fsys2, ok := fsys.(FSSupportingReadlink); ok {
		return fsys2.Lstat(name)
	} else if fsys2, ok := fsys.(FSSupportingStat); ok {
		return fsys2.Stat(name)
	} else {
		return nil, &PathError{
			Op:   "Readlink",
			Path: name,
			Err:  fmt.Errorf("filesystem type %T does not support Lstat nor Stat", fsys),
		}
	}
}

// MkSymlink creates a symlink.
//
// It only works on filesystems that support FSSupportingMkSymlink, and errors otherwise.
func MkSymlink(fsys FS, name, target string) error {
	if fsys2, ok := fsys.(FSSupportingMkSymlink); ok {
		return fsys2.MkSymlink(name, target)
	} else {
		return &PathError{
			Op:   "MkSymlink",
			Path: name,
			Err:  fmt.Errorf("filesystem type %T does not support MkSymlink", fsys),
		}
	}
}

// IsPathFile peeks at the filesystem to see if the given path contains a regular file.
//
// If there is any error, false will returned.
//
// This is a convenience function.  Similar results can be obtained by use of Stat or other functions.
//
// Be mindful of TOCTOU if using this function; in many cases,
// and especially if security and concurrent filesystem access is a concern,
// it may be better to try open the file immediately, and then examine it.
func IsPathFile(fsys FS, name string) (bool, error) {
	fi, err := Stat(fsys, name)
	if err != nil {
		return false, err
	}
	return fi.Mode()&ModeType == 0, nil
}

// IsPathDir peeks at the filesystem to see if the given path contains a directory.
//
// If there is any error, false will returned.
//
// This is a convenience function.  Similar results can be obtained by use of Stat or other functions.
//
// Be mindful of TOCTOU if using this function; in many cases,
// and especially if security and concurrent filesystem access is a concern,
// it may be better to try to perform the intended operation immediately, and check the errors from that operation.
func IsPathDir(fsys FS, name string) (bool, error) {
	fi, err := Stat(fsys, name)
	if err != nil {
		return false, err
	}
	return fi.Mode()&ModeType == ModeDir, nil
}
