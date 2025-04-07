package monkey

import (
	"bytes"
	"fmt"
)

func ExamplePatch() {
	patch := func(a ...interface{}) (int, error) {
		return fmt.Print("what?!")
	}
	pg := Patch(fmt.Println, patch)
	defer pg.Unpatch()

	// output: what?!
	fmt.Println("hello!")
}

func ExamplePatchMethod() {
	var r *bytes.Reader
	patch := func(b []byte) (int, error) {
		return 0, nil
	}
	pg := PatchMethod(r, "Read", patch)
	defer pg.Unpatch()

	reader := bytes.NewReader([]byte("hello"))
	buf := make([]byte, 1024)

	// output: 0 <nil>
	fmt.Println(reader.Read(buf))
}
