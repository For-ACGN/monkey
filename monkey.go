package monkey

import (
	"reflect"
)

// Patch is used to patch common function.
func Patch(target, patch interface{}) *PatchGuard {
	pg := PatchGuard{}
	t := reflect.ValueOf(target)
	d := reflect.ValueOf(patch)
	pg.patchFunc(t, d)
	return &pg
}
