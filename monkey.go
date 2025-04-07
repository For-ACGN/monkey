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

// PatchMethod is used to patch structure methods, it supports unexported
// methods and unexported structure exported and unexported methods, usually
// the unexported structure is from interface.
func PatchMethod(target interface{}, method string, patch interface{}) *PatchGuard {
	tType := reflect.TypeOf(target)
	pType := reflect.TypeOf(patch)
	switch k := tType.Kind(); k {
	case reflect.Struct:
	case reflect.Pointer:
		if tType.Elem().Kind() == reflect.Struct {
			break
		}
		fallthrough
	default:
		panic("target is not a structure or pointer")
	}
	if pType.Kind() != reflect.Func {
		panic("patch is not a function")
	}
	tValue := reflect.ValueOf(target)
	pValue := reflect.ValueOf(patch)
	pg := new(PatchGuard)
	pg.patchMethod(tValue, method, pValue)
	return pg
}
