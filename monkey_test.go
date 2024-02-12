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
	t.Run("common", func(t *testing.T) {
		w := new(testdata.Writer)
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
		defer pg.Unpatch()

		n, err = w.Print()
		require.NoError(t, err)
		require.Equal(t, 4, n)

		pg.Unpatch()
		n, err = w.Print()
		require.NoError(t, err)
		require.Equal(t, 7, n)

		pg.Restore()
		n, err = w.Print()
		require.NoError(t, err)
		require.Equal(t, 4, n)
	})

	t.Run("implement interface", func(t *testing.T) {
		w := new(testdata.Writer)
		n, err := w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Equal(t, 7, n)

		patch := func(*testdata.Writer, []byte) (int, error) {
			return 0, nil
		}
		pg := PatchMethod(w, "Write", patch)
		defer pg.Unpatch()

		n, err = w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Zero(t, n)

		pg.Unpatch()
		n, err = w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Equal(t, 7, n)

		pg.Restore()
		n, err = w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Zero(t, n)
	})

	t.Run("not pointer receiver", func(t *testing.T) {
		w := testdata.Writer{}
		n, err := w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Equal(t, 7, n)
		n, err = w.Println()
		require.NoError(t, err)
		require.Equal(t, 8, n)

		patch := func(testdata.Writer) (int, error) {
			return 0, nil
		}
		pg := PatchMethod(w, "Println", patch)
		defer pg.Unpatch()

		n, err = w.Println()
		require.NoError(t, err)
		require.Zero(t, n)

		pg.Unpatch()
		n, err = w.Println()
		require.NoError(t, err)
		require.Equal(t, 8, n)

		pg.Restore()
		n, err = w.Println()
		require.NoError(t, err)
		require.Zero(t, n)
	})
}
