package monkey

import (
	"fmt"
	"reflect"
	"unicode"
	"unsafe"

	"github.com/For-ACGN/monkey/creflect"
)

// PatchGuard contains original and patch data for unpatch.
type PatchGuard struct {
	target   uintptr
	original []byte
	patch    []byte
}

func (pg *PatchGuard) patchFunc(target, patch reflect.Value) {
	checkFunc(target.Type(), patch.Type())
	targetAddr := *(*uintptr)(getPointer(target))
	patchAddr := uintptr(getPointer(patch))
	pg.applyPatch(targetAddr, patchAddr)
}

func (pg *PatchGuard) patchMethod(target reflect.Value, method string, patch reflect.Value) {
	if method == "" {
		panic("empty method name ")
	}
	var (
		targetAddr uintptr
		patchAddr  uintptr
	)
	if unicode.IsLower([]rune(method)[0]) {
		m, ok := creflect.MethodByName(target.Type(), method)
		if !ok {
			panic(fmt.Sprintf("failed to get method by name: %s\n", method))
		}
		targetAddr = *(*uintptr)(m)
		patchAddr = uintptr(getPointer(patch))
	} else {
		m, ok := target.Type().MethodByName(method)
		if !ok {
			panic(fmt.Sprintf("failed to get method by name: %s\n", method))
		}
		// process when receiver is private structure
		patchType := patch.Type()
		checkFunc(m.Type, patchType)
		numArgs := patchType.NumIn()
		wrapper := reflect.MakeFunc(m.Type, func(args []reflect.Value) []reflect.Value {
			newArgs := make([]reflect.Value, numArgs)
			for i := 0; i < len(newArgs); i++ {
				newArgs[i] = args[i].Convert(patchType.In(i))
			}
			if patchType.IsVariadic() {
				return patch.CallSlice(newArgs)
			}
			return patch.Call(newArgs)
		})
		checkFunc(m.Type, wrapper.Type())
		targetAddr = *(*uintptr)(getPointer(m.Func))
		patchAddr = uintptr(getPointer(wrapper))
	}
	pg.applyPatch(targetAddr, patchAddr)
}

func checkFunc(target, patch reflect.Type) {
	if patch.Kind() != reflect.Func {
		panic("patch is not a function")
	}
	// check the number of the function parameter and return value are equal
	invalidIn := target.NumIn() != patch.NumIn()
	invalidOut := target.NumOut() != patch.NumOut()
	invalidVar := target.NumIn() == patch.NumIn() &&
		target.IsVariadic() != patch.IsVariadic()
	if invalidIn || invalidOut || invalidVar {
		const format = "target type(%s) and patch type(%s) are different"
		panic(fmt.Sprintf(format, target, patch))
	}
	// check the function parameters type are equal
	for i, size := 0, patch.NumIn(); i < size; i++ {
		targetIn := target.In(i)
		patchIn := patch.In(i)
		if targetIn.AssignableTo(patchIn) {
			continue
		}
		const format = "target type(%s) and patch type(%s) are different"
		panic(fmt.Sprintf(format, target, patch))
	}
	// check the function return values type are equal
	for i, size := 0, patch.NumOut(); i < size; i++ {
		targetOut := target.Out(i)
		patchOut := patch.Out(i)
		if targetOut.AssignableTo(patchOut) {
			continue
		}
		const format = "target type(%s) and patch type(%s) are different"
		panic(fmt.Sprintf(format, target, patch))
	}
}

func (pg *PatchGuard) applyPatch(target, patch uintptr) {
	jmp := buildJMPDirective(patch)
	mem := readMemory(target, len(jmp))
	original := make([]byte, len(mem))
	copy(original, mem)
	writeMemory(target, jmp)
	pg.target = target
	pg.original = original
	pg.patch = jmp
}

// Unpatch is used to recovery the original about target.
func (pg *PatchGuard) Unpatch() {
	writeMemory(pg.target, pg.original)
}

// Restore is used to patch the target again.
func (pg *PatchGuard) Restore() {
	writeMemory(pg.target, pg.patch)
}

func readMemory(address uintptr, size int) []byte {
	var b []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&b)) // #nosec
	sh.Data = address
	sh.Len = size
	sh.Cap = size
	return b
}

type reflectValue struct {
	_   uintptr
	ptr unsafe.Pointer
}

func getPointer(v reflect.Value) unsafe.Pointer {
	return (*reflectValue)(unsafe.Pointer(&v)).ptr // #nosec
}
