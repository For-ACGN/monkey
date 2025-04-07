package monkey

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/For-ACGN/monkey/testpkg"
)

func TestPatchFunc(t *testing.T) {
	output := fmt.Sprintln("hello!")
	require.Equal(t, "hello!\n", output)

	patch := func(a ...interface{}) string {
		return "what!!!\n"
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

func TestPatchMethod_Exported(t *testing.T) {
	t.Run("common", func(t *testing.T) {
		w := new(testpkg.Writer)
		n, err := w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Equal(t, 7, n)
		n, err = w.Print()
		require.NoError(t, err)
		require.Equal(t, 8, n)

		patch := func(*testpkg.Writer) (int, error) {
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
		require.Equal(t, 8, n)

		pg.Restore()
		n, err = w.Print()
		require.NoError(t, err)
		require.Equal(t, 4, n)
	})

	t.Run("ignored receiver", func(t *testing.T) {
		var w io.Writer = new(testpkg.Writer)
		n, err := w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Equal(t, 7, n)

		patch := func([]byte) (int, error) {
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

	t.Run("implement interface", func(t *testing.T) {
		var w io.Writer = new(testpkg.Writer)
		n, err := w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Equal(t, 7, n)

		patch := func(*testpkg.Writer, []byte) (int, error) {
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
		w := testpkg.Writer{}
		n, err := w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Equal(t, 7, n)
		n, err = w.Println()
		require.NoError(t, err)
		require.Equal(t, 7, n)

		patch := func(testpkg.Writer) (int, error) {
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
		require.Equal(t, 7, n)

		pg.Restore()
		n, err = w.Println()
		require.NoError(t, err)
		require.Zero(t, n)
	})
}

func TestPatchMethod_Unexported(t *testing.T) {
	t.Run("common", func(t *testing.T) {
		w := new(testpkg.Writer)
		n, err := w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Equal(t, 7, n)
		n, err = w.Print()
		require.NoError(t, err)
		require.Equal(t, 8, n)

		patch := func(*testpkg.Writer) (int, error) {
			return fmt.Println("oh!")
		}
		pg := PatchMethod(w, "print", patch)
		defer pg.Unpatch()

		n, err = w.Print()
		require.NoError(t, err)
		require.Equal(t, 5, n)

		pg.Unpatch()
		n, err = w.Print()
		require.NoError(t, err)
		require.Equal(t, 8, n)

		pg.Restore()
		n, err = w.Print()
		require.NoError(t, err)
		require.Equal(t, 5, n)
	})

	t.Run("ignored receiver", func(t *testing.T) {
		w := new(testpkg.Writer)
		n, err := w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Equal(t, 7, n)
		n, err = w.Print()
		require.NoError(t, err)
		require.Equal(t, 8, n)

		patch := func() (int, error) {
			return fmt.Println("oh!")
		}
		pg := PatchMethod(w, "print", patch)
		defer pg.Unpatch()

		n, err = w.Print()
		require.NoError(t, err)
		require.Equal(t, 5, n)

		pg.Unpatch()
		n, err = w.Print()
		require.NoError(t, err)
		require.Equal(t, 8, n)

		pg.Restore()
		n, err = w.Print()
		require.NoError(t, err)
		require.Equal(t, 5, n)
	})

	t.Run("not pointer receiver", func(t *testing.T) {
		w := testpkg.Writer{}
		n, err := w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Equal(t, 7, n)
		n, err = w.Println()
		require.NoError(t, err)
		require.Equal(t, 7, n)

		patch := func(testpkg.Writer) (int, error) {
			return 0, nil
		}
		pg := PatchMethod(w, "println", patch)
		defer pg.Unpatch()

		n, err = w.Println()
		require.NoError(t, err)
		require.Zero(t, n)

		pg.Unpatch()
		n, err = w.Println()
		require.NoError(t, err)
		require.Equal(t, 7, n)

		pg.Restore()
		n, err = w.Println()
		require.NoError(t, err)
		require.Zero(t, n)
	})
}

func TestPatchMethod_Interface(t *testing.T) {
	w := testpkg.NewWriter()

	t.Run("exported method", func(t *testing.T) {
		n, err := w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Equal(t, 7, n)

		patch := func([]byte) (int, error) {
			return 0, nil
		}
		pg := PatchMethod(w, "Write", patch)
		defer pg.Unpatch()

		n, err = w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Zero(t, n)
	})

	t.Run("unexported method", func(t *testing.T) {
		n, err := w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Equal(t, 7, n)

		patch := func() (int, error) {
			return fmt.Println("oh!")
		}
		pg := PatchMethod(w, "print", patch)
		defer pg.Unpatch()

		n, err = w.Write([]byte("hello!"))
		require.NoError(t, err)
		require.Equal(t, 4, n)
	})
}
