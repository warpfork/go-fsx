package fsx_test

import (
	"fmt"

	"github.com/warpfork/go-fsx"
	"github.com/warpfork/go-fsx/osfs"
)

func ExampleHello() {
	fsys := osfs.DirFS("/tmp/")
	fsx.Mkdir(fsys, "hello", 0777)
	fsx.WriteFile(fsys, "hello/world.txt", 0666, []byte(`hello world!`))
	body, _ := fsx.ReadFile(fsys, "hello/world.txt")
	fmt.Printf("%s\n", body)

	// Output:
	// hello world!
}
