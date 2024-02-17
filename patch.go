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
	checkFuncType(target.Type(), patch.Type())
	targetAddr := *(*uintptr)(getPointer(target))
	patchAddr := uintptr(getPointer(patch))
	pg.applyPatch(targetAddr, patchAddr)
}

func (pg *PatchGuard) patchMethod(target reflect.Value, method string, patch reflect.Value) {
	if method == "" {
		panic("empty method name")
	}
	if unicode.IsUpper([]rune(method)[0]) {
		pg.patchPublicMethod(target, method, patch)
	} else {
		pg.patchPrivateMethod(target, method, patch)
	}
}

func (pg *PatchGuard) patchPublicMethod(target reflect.Value, method string, patch reflect.Value) {
	m, ok := target.Type().MethodByName(method)
	if !ok {
		panic(fmt.Sprintf("failed to get method by name: %s\n", method))
	}
	// check the type of first argument in patch is the receiver type
	patchType := patch.Type()
	if patchType.NumIn() == 0 || patchType.In(0) != target.Type() {
		// build new patch function type
		numIn := patchType.NumIn()
		in := make([]reflect.Type, numIn+1)
		in[0] = target.Type()
		for i := 0; i < numIn; i++ {
			in[i+1] = patchType.In(i)
		}
		numOut := patchType.NumOut()
		out := make([]reflect.Type, numOut)
		for i := 0; i < numOut; i++ {
			out[i] = patchType.Out(i)
		}
		funcType := reflect.FuncOf(in, out, patchType.IsVariadic())
		// create new patch function
		rawPatch := patch
		patch = reflect.MakeFunc(funcType, func(args []reflect.Value) []reflect.Value {
			if rawPatch.Type().IsVariadic() {
				return rawPatch.CallSlice(args[1:])
			} else {
				return rawPatch.Call(args[1:])
			}
		})
		patchType = patch.Type()
	}
	// process when receiver is private structure
	checkFuncType(m.Type, patchType)
	targetAddr := *(*uintptr)(getPointer(m.Func))
	patchAddr := uintptr(getPointer(patch))
	pg.applyPatch(targetAddr, patchAddr)
}

func (pg *PatchGuard) patchPrivateMethod(target reflect.Value, method string, patch reflect.Value) {
	m, ok := creflect.MethodByName(target.Type(), method)
	if !ok {
		panic(fmt.Sprintf("failed to get method by name: %s\n", method))
	}
	// check the first argument in patch is the receiver
	patchType := patch.Type()
	if patchType.NumIn() == 0 || patchType.In(0) != target.Type() {

	}

	targetAddr := *(*uintptr)(m)
	patchAddr := uintptr(getPointer(patch))
	pg.applyPatch(targetAddr, patchAddr)
}

func checkFuncType(target, patch reflect.Type) {
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
	for i := 0; i < target.NumIn(); i++ {
		targetIn := target.In(i)
		patchIn := patch.In(i)
		if targetIn == patchIn {
			continue
		}
		// if targetIn.Size() == patchIn.Size() && targetIn.ConvertibleTo(patchIn) {
		// 	continue
		// }
		const format = "target type(%s) and patch type(%s) are different"
		panic(fmt.Sprintf(format, target, patch))
	}
	// check the function return values type are equal
	for i := 0; i < target.NumOut(); i++ {
		targetOut := target.Out(i)
		patchOut := patch.Out(i)
		if targetOut == patchOut {
			continue
		}
		// if targetOut.Size() == patchOut.Size() && targetOut.ConvertibleTo(patchOut) {
		// 	continue
		// }
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
