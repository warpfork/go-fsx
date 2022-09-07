go-fsx -- golang filesystems interfaces, extended
=================================================

Golang introduced the [`io/fs.FS`](https://pkg.go.dev/io/fs) interface in Go 1.16.
[`go-fsx`](https://pkg.go.dev/github.com/warpfork/go-fsx) extends it.

While having the standard library's `fs.FS` interface was a very welcome improvement to the golang library ecosystem,
it doesn't (yet) go far enough:
the interface only covers read operations (no write support at all!),
and many features in common filesystems (such as understanding of symlinks) are absent.

This package, `fsx`, takes the style of the `io/fs` package, and extends it with more features:

- `fsx` provides the ability to write files (using `OpenFile`, which is much like the `os` package feature you're already familiar with)
- `fsx` provides the ability to create directories
- `fsx` provides the ability to delete files and directories
- `fsx` provides features for reading symlinks, and creating them (WIP)

Everything is done with the intention of feeling normal, and being a smooth extension of what we already know:

- `fsx` does everything it does in the functional idiom already used in `io/fs`, so it feels "natural" -- `fsx` has just got _more_ of it.
- `fsx` still uses `fs.FileInfo`, `fs.FileMode`, and `fs.File` -- no changes there.
- All of the `fsx` functions take an `fs.FS`, and do feature detection internally -- so you can keep using `fs.FS` in code that's already passing that interface around!
- As with `io/fs`, we attempt to add new convenient behaviors as package-scope functions... so that the `fsx.FS` interface doesn't grow over time, and bumping library versions forward is easy, even if you created your own unique implementations of the interface.

Additionally, we alias other `fs` package features, and relevant `os` constants, into this package --
so that you can have just one thing on your import path when working with filesystems.
Less to think about is better.


Example Usage
-------------

### Hello World

```go
import (
	"github.com/warpfork/go-fsx"
	"github.com/warpfork/go-fsx/osfs"
)

func ExampleHello() {
	fsys := osfs.DirFS("/tmp/")
	fsx.Mkdir(fsys, "hello", 0777)
	fsx.WriteFile(fsys, "hello/world.txt", 0666, []byte(`hello world!`))
	body, err := fsx.ReadFile(fsys, "hello/world.txt")
}
```

Note that the function for creating an FS is in the `osfs` subpackage.
Other implementations of the FS interface can be in other packages!
The `fsx` package is only interfaces.
(This is similar to how `os.DirFS` is used to get an `io/fs.FS`, in the standard library.)

### Importing as 'fs'

You may also choose to import fsx as just "fs", if you so choose:

```go
import (
	fs "github.com/warpfork/go-fsx"
)
```

This is a reasonably safe choice since we alias everything from the `io/fs` package,
so it should never be necessary to import both at the same time.


Future Work
-----------

- `chmod`, `chown`, and other similar operations are not yet present.  (PRs welcome!)
- Your ideas here?


Related Work
------------

Writable FS interfaces have been discussed before!
In particular, https://github.com/golang/go/issues/45757 contains a very rich discussion.

If anyone is interested in taking this code further upstream, either as reference material,
or verbatim copied, please, be my guest.


License
-------

You can have it as Apache-2.0 OR MIT OR BSD-3-Clause, or really anything you want.
