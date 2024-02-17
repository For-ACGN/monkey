package monkey

import (
	"reflect"
)

// Patch is used to patch common function.
func Patch(target, patch interface{}) *PatchGuard {
	t := reflect.ValueOf(target)
	p := reflect.ValueOf(patch)
	if t.Kind() != reflect.Func {
		panic("target is not a function")
	}
	if p.Kind() != reflect.Func {
		panic("patch is not a function")
	}
	pg := new(PatchGuard)
	pg.patchFunc(t, p)
	return pg
}

// PatchMethod is used to patch structure methods, it supports private
// methods and private structure public and private methods, usually
// the private structure is from interface.
func PatchMethod(target interface{}, method string, patch interface{}) *PatchGuard {
	t := reflect.ValueOf(target)
	p := reflect.ValueOf(patch)
	switch k := t.Kind(); k {
	case reflect.Struct:
	case reflect.Pointer:
		if t.Elem().Kind() == reflect.Struct {
			break
		}
		fallthrough
	default:
		panic("target is not a structure or pointer")
	}
	if p.Kind() != reflect.Func {
		panic("patch is not a function")
	}
	pg := new(PatchGuard)
	pg.patchMethod(t, method, p)
	return pg
}
