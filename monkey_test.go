package monkey

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/For-ACGN/monkey/testdata"
)

func TestPatch(t *testing.T) {
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

func TestPatchMethod(t *testing.T) {
	w := new(testdata.Writer)

	t.Run("interface function", func(t *testing.T) {
		n, err := w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Equal(t, 7, n)

		patch := func(*testdata.Writer, []byte) (int, error) {
			return 0, nil
		}
		pg := PatchMethod(w, "Write", patch)

		n, err = w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Zero(t, n)

		pg.Unpatch()
		n, err = w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Equal(t, 7, n)
	})

	t.Run("common function", func(t *testing.T) {
		n, err := w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Equal(t, 7, n)
		n, err = w.Print()
		require.NoError(t, err)
		require.Equal(t, 7, n)

		patch := func(*testdata.Writer) (int, error) {
			return fmt.Println("oh!")
		}
		pg := PatchMethod(w, "Print", patch)

		n, err = w.Print()
		require.NoError(t, err)
		require.Equal(t, 4, n)

		pg.Unpatch()
		n, err = w.Print()
		require.NoError(t, err)
		require.Equal(t, 7, n)
	})
}
