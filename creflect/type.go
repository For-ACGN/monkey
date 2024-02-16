// customized reflect package for monkey, copy most code from go/src/reflect/type.go

package creflect

import (
	"reflect"
	"unsafe"
)

type tflag uint8
type nameOff int32 // offset to a name
type typeOff int32 // offset to an *rtype
type textOff int32 // offset from top of text section

// rtype is the common implementation of most values.
// rtype must be kept in sync with ../runtime/type.go:/^type._type.
type rtype struct {
	size       uintptr
	ptrdata    uintptr // number of bytes in the type that can contain pointers
	hash       uint32  // hash of type; avoids computation in hash tables
	tflag      tflag   // extra type information flags
	align      uint8   // alignment of variable with this type
	fieldAlign uint8   // alignment of struct field with this type
	kind       uint8   // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal     func(unsafe.Pointer, unsafe.Pointer) bool
	gcdata    *byte   // garbage collection data
	str       nameOff // string form
	ptrToThis typeOff // type for pointer to this type, may be zero
}

func newType(t reflect.Type) *rtype {
	i := *(*funcValue)(unsafe.Pointer(&t)) // #nosec
	r := (*rtype)(i.p)
	return r
}

type funcValue struct {
	_ uintptr
	p unsafe.Pointer
}

func funcPointer(v reflect.Method, ok bool) (unsafe.Pointer, bool) {
	return (*funcValue)(unsafe.Pointer(&v.Func)).p, ok // #nosec
}

// MethodByName returns the method with that name in the type's
// method set and a boolean indicating if the method was found.
//
// For a non-interface type T or *T, the returned Method's Type and Func
// fields describe a function whose first argument is the receiver.
//
// For an interface type, the returned Method's Type field gives the
// method signature, without a receiver, and the Func field is nil.
func MethodByName(r reflect.Type, name string) (fn unsafe.Pointer, ok bool) {
	t := newType(r)
	if r.Kind() == reflect.Interface {
		return funcPointer(r.MethodByName(name))
	}
	ut := t.uncommon(r)
	if ut == nil {
		return nil, false
	}
	for _, p := range ut.methods() {
		if t.nameOff(p.name).name() == name {
			return t.Method(p), true
		}
	}
	return nil, false
}

func (t *rtype) Method(p method) unsafe.Pointer {
	tfn := t.textOff(p.tfn)
	return unsafe.Pointer(&tfn) // #nosec
}

//go:linkname resolveTextOff reflect.resolveTextOff
func resolveTextOff(rtype unsafe.Pointer, off int32) unsafe.Pointer

//go:linkname resolveNameOff reflect.resolveNameOff
func resolveNameOff(ptrInModule unsafe.Pointer, off int32) unsafe.Pointer

func (t *rtype) textOff(off textOff) unsafe.Pointer {
	return resolveTextOff(unsafe.Pointer(t), int32(off)) // #nosec
}

func (t *rtype) nameOff(off nameOff) name {
	return name{(*byte)(resolveNameOff(unsafe.Pointer(t), int32(off)))} // #nosec
}

const (
	tflagUncommon tflag = 1 << 0
)

// uncommonType is present only for defined types or types with methods
type uncommonType struct {
	pkgPath nameOff // import path; empty for built-in types like int, string
	mcount  uint16  // number of methods
	xcount  uint16  // number of exported methods
	moff    uint32  // offset from this uncommontype to [mcount]method
	_       uint32  // unused
}

// ptrType represents a pointer type.
type ptrType struct {
	rtype
	elem *rtype // pointer element (pointed at) type
}

// funcType represents a function type.
type funcType struct {
	rtype
	inCount  uint16
	outCount uint16 // top bit is set if last input parameter is ...
}

func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x) // #nosec
}

// interfaceType represents an interface type.
type interfaceType struct {
	rtype
	pkgPath name      // import path
	methods []imethod // sorted by hash
}

type imethod struct {
	name nameOff // name of method
	typ  typeOff // .(*FuncType) underneath
}

type stringHeader struct {
	data unsafe.Pointer
	len  int
}

func (t *rtype) uncommon(r reflect.Type) *uncommonType {
	if t.tflag&tflagUncommon == 0 {
		return nil
	}
	switch r.Kind() {
	case reflect.Ptr:
		type u struct {
			ptrType
			u uncommonType
		}
		return &(*u)(unsafe.Pointer(t)).u // #nosec
	case reflect.Func:
		type u struct {
			funcType
			u uncommonType
		}
		return &(*u)(unsafe.Pointer(t)).u // #nosec
	case reflect.Interface:
		type u struct {
			interfaceType
			u uncommonType
		}
		return &(*u)(unsafe.Pointer(t)).u // #nosec
	case reflect.Struct:
		type u struct {
			interfaceType
			u uncommonType
		}
		return &(*u)(unsafe.Pointer(t)).u // #nosec
	default:
		return nil
	}
}

// Method on non-interface type
type method struct {
	name nameOff // name of method
	mtyp typeOff // method type (without receiver)
	ifn  textOff // fn used in interface call (one-word receiver)
	tfn  textOff // fn used for normal method call
}

func (t *uncommonType) methods() []method {
	if t.mcount == 0 {
		return nil
	}
	return (*[1 << 16]method)(add(unsafe.Pointer(t), uintptr(t.moff)))[:t.mcount:t.mcount] // #nosec
}
