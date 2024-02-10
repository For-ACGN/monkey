package monkey

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFunction(t *testing.T) {
	output := fmt.Sprintln("hello!")
	require.Equal(t, "hello!\n", output)

	patch := func(a ...interface{}) string {
		return fmt.Sprint("what!!!\n")
	}
	pg := Patch(fmt.Sprintln, patch)
	defer pg.Unpatch()

	output = fmt.Sprintln("hello?")
	require.Equal(t, "what!!!\n", output)
	output = fmt.Sprintln("??????")
	require.Equal(t, "what!!!\n", output)

	pg.Unpatch()
	output = fmt.Sprintln("hello!")
	require.Equal(t, "hello!\n", output)
	output = fmt.Sprintln("world!")
	require.Equal(t, "world!\n", output)

	pg.Restore()
	output = fmt.Sprintln("hello?")
	require.Equal(t, "what!!!\n", output)
	output = fmt.Sprintln("??????")
	require.Equal(t, "what!!!\n", output)
}
