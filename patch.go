package monkey

import (
	"fmt"
	"reflect"
	"unsafe"
)

// PatchGuard contains original and patch data for unpatch.
type PatchGuard struct {
	original []byte
	patch    reflect.Value
	value    reflect.Value
}

func (pg *PatchGuard) patchFunc(target, patch reflect.Value) {
	pg.checkFunc(target, patch)
	assTarget := *(*uintptr)(getPointer(target))
	pg.original = replace(assTarget, uintptr(getPointer(patch)))
	pg.patch = patch
}

func (pg *PatchGuard) checkFunc(target, patch reflect.Value) {
	if target.Kind() != reflect.Func {
		panic("target is not a function")
	}
	if patch.Kind() != reflect.Func {
		panic("patch is not a function")
	}

	targetType := target.Type()
	patchType := patch.Type()

	// check the number of the function parameter and return value are equal
	invalidIn := targetType.NumIn() < patchType.NumIn()
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

func replace(target, patch uintptr) []byte {
	code := buildJMPDirective(patch)
	bytes := entryAddress(target, len(code))
	original := make([]byte, len(bytes))
	copy(original, bytes)
	modifyBinary(target, code)
	return original
}

func entryAddress(p uintptr, l int) []byte {
	var b []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&b)) // #nosec
	sh.Data = p
	sh.Len = l
	sh.Cap = l
	return b
}

type funcValue struct {
	_ uintptr
	p unsafe.Pointer
}

func getPointer(v reflect.Value) unsafe.Pointer {
	return (*funcValue)(unsafe.Pointer(&v)).p // #nosec
}
