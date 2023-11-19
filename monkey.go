package monkey

import (
	"fmt"
	"reflect"
	"syscall"
	"unsafe"
)

type PatchGuard struct {
	original []byte
	patch    reflect.Value
	value    reflect.Value
}

// Patch is a wrapper about monkey.Patch.
func Patch(target, patch interface{}) *PatchGuard {
	pg := PatchGuard{}
	t := reflect.ValueOf(target)
	d := reflect.ValueOf(patch)
	pg.apply(t, d)
	return &pg
}

func (pg *PatchGuard) apply(target, patch reflect.Value) {
	pg.check(target, patch)
	assTarget := *(*uintptr)(getPointer(target))
	pg.original = replace(assTarget, uintptr(getPointer(patch)))
	pg.patch = patch
}

func (pg *PatchGuard) check(target, patch reflect.Value) {
	if target.Kind() != reflect.Func {
		panic("target is not a func")
	}
	if patch.Kind() != reflect.Func {
		panic("patch is not a func")
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

type funcValue struct {
	_ uintptr
	p unsafe.Pointer
}

func getPointer(v reflect.Value) unsafe.Pointer {
	return (*funcValue)(unsafe.Pointer(&v)).p
}

func entryAddress(p uintptr, l int) []byte {
	var b []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&b)) // #nosec
	sh.Data = p
	sh.Len = l
	sh.Cap = l
	return b
}

func pageStart(ptr uintptr) uintptr {
	return ptr & ^(uintptr(syscall.Getpagesize() - 1))
}
