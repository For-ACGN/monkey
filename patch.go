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
	jmp := buildJMPDirective(patchAddr)
	data := readMemory(targetAddr, len(jmp))
	original := make([]byte, len(data))
	copy(original, data)
	modifyBinary(targetAddr, jmp)
	pg.target = targetAddr
	pg.original = original
	pg.patch = jmp
}

func (pg *PatchGuard) patchMethod(target reflect.Value, method string, patch reflect.Value) {
	if method == "" {
		panic("empty method")
	}
	if unicode.IsLower([]rune(method)[0]) {
		creflect.MethodByName(target.Type(), method)
	} else {
		target.MethodByName(method)
	}

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

// Unpatch is used to recovery the original about target.
func (pg *PatchGuard) Unpatch() {
	modifyBinary(pg.target, pg.original)
}

// Restore is used to patch the target again.
func (pg *PatchGuard) Restore() {
	modifyBinary(pg.target, pg.patch)
}

func readMemory(p uintptr, l int) []byte {
	var b []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&b)) // #nosec
	sh.Data = p
	sh.Len = l
	sh.Cap = l
	return b
}

type reflectValue struct {
	_   uintptr
	ptr unsafe.Pointer
}

func getPointer(v reflect.Value) unsafe.Pointer {
	return (*reflectValue)(unsafe.Pointer(&v)).ptr // #nosec
}
