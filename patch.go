package monkey

import (
	"fmt"
	"reflect"
	"unsafe"
)

// PatchGuard contains original and patch data for unpatch.
type PatchGuard struct {
	original []byte
	target   reflect.Value
	patch    reflect.Value
}

func (pg *PatchGuard) patchFunc(target, patch reflect.Value) {
	pg.checkFunc(target, patch)
	tAddr := *(*uintptr)(getPointer(target))
	pAddr := uintptr(getPointer(patch))
	pg.original = replaceFunc(tAddr, pAddr)
}

func (pg *PatchGuard) checkFunc(target, patch reflect.Value) {
	if patch.Kind() != reflect.Func {
		panic("patch is not a function")
	}

	targetType := target.Type()
	patchType := patch.Type()

	// check the number of the function parameter and return value are equal
	invalidIn := targetType.NumIn() != patchType.NumIn()
	invalidOut := targetType.NumOut() != patchType.NumOut()
	invalidVar := targetType.NumIn() == patchType.NumIn() &&
		targetType.IsVariadic() != patchType.IsVariadic()
	if invalidIn || invalidOut || invalidVar {
		const format = "target type(%s) and patch type(%s) are different"
		panic(fmt.Sprintf(format, targetType, patchType))
	}

	// check the function parameters type are equal
	for i, size := 0, patchType.NumIn(); i < size; i++ {
		targetIn := targetType.In(i)
		patchIn := patchType.In(i)
		if targetIn.AssignableTo(patchIn) {
			continue
		}
		const format = "target type(%s) and patch type(%s) are different"
		panic(fmt.Sprintf(format, targetType, patchType))
	}

	// check the function return values type are equal
	for i, size := 0, patchType.NumOut(); i < size; i++ {
		targetOut := targetType.Out(i)
		patchOut := patchType.Out(i)
		if targetOut.AssignableTo(patchOut) {
			continue
		}
		const format = "target type(%s) and patch type(%s) are different"
		panic(fmt.Sprintf(format, targetType, patchType))
	}
}

func replaceFunc(target, patch uintptr) []byte {
	jmp := buildJMPDirective(patch)
	data := readMemory(target, len(jmp))
	original := make([]byte, len(data))
	copy(original, data)
	modifyBinary(target, jmp)
	return original
}

// Unpatch is used to recovery the original about target.
func (pg *PatchGuard) Unpatch() {
	address := *(*uintptr)(getPointer(pg.target))
	modifyBinary(address, pg.original)
}

// Restore is used to patch the target again.
func (pg *PatchGuard) Restore() {
	pg.patchFunc(pg.target, pg.patch)
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
