package monkey

import (
	"fmt"
	"reflect"
)

// Patch is used to patch common function.
func Patch(target, patch interface{}) *PatchGuard {
	t := reflect.ValueOf(target)
	p := reflect.ValueOf(patch)
	pg := PatchGuard{
		target: t,
		patch:  p,
	}
	switch kind := t.Kind(); kind {
	case reflect.Func:
		pg.patchFunc(t, p)
	default:
		panic(fmt.Sprintf("invalid target kind: %s", kind))
	}
	return &pg
}
