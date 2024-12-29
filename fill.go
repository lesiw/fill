package fill

import (
	"fmt"
	"math/rand/v2"
	"reflect"
)

func Do(a any, rng *rand.Rand) {
	if rng == nil {
		panic("bad parameter: nil rand.Rand")
	}
	val := reflect.ValueOf(a)
	if val.Kind() != reflect.Ptr {
		panic("bad parameter: value to fill must be pointer")
	}
	val = val.Elem()
	fillValue(val, rng)
}

func fillValue(val reflect.Value, rng *rand.Rand) {
	k := val.Kind()
	switch {
	case intKind(k):
		val.SetInt(rng.Int64())
	case uintKind(k):
		val.SetUint(rng.Uint64())
	case floatKind(k):
		val.SetFloat(rng.Float64())
	case complexKind(k):
		val.SetComplex(complex(rng.Float64(), rng.Float64()))
	case k == reflect.String:
		fillString(val, rng)
	case k == reflect.Bool:
		val.SetBool(rng.IntN(2) == 0)
	case k == reflect.Struct:
		fillStruct(val, rng)
	case k == reflect.Slice:
		fillSlice(val, rng)
	case k == reflect.Map:
		fillMap(val, rng)
	case k == reflect.Array:
		fillArray(val, rng)
	case k == reflect.Chan:
		fillChan(val, rng)
	case k == reflect.Pointer:
		val.Set(reflect.New(val.Type().Elem()))
		fillValue(val.Elem(), rng)
	case k == reflect.Interface || k == reflect.Func:
		// Can't fill.
	default:
		panic(fmt.Sprintf("unhandled type: %s", val.Kind()))
	}
}

func fillStruct(val reflect.Value, rng *rand.Rand) {
	for i := range val.NumField() {
		f := val.Field(i)
		if !f.CanSet() {
			continue
		}
		fillValue(f, rng)
	}
}

func fillSlice(val reflect.Value, rng *rand.Rand) {
	for range rng.IntN(16) {
		e := reflect.New(val.Type().Elem()).Elem()
		fillValue(e, rng)
		val.Set(reflect.Append(val, e))
	}
}

func fillMap(val reflect.Value, rng *rand.Rand) {
	sz := rng.IntN(16)
	val.Set(reflect.MakeMapWithSize(val.Type(), sz))
	for range sz {
		k := reflect.New(val.Type().Key()).Elem()
		fillValue(k, rng)
		v := reflect.New(val.Type().Elem()).Elem()
		fillValue(v, rng)
		val.SetMapIndex(k, v)
	}
}

func fillChan(val reflect.Value, _ *rand.Rand) {
	reflect.MakeChan(val.Type(), 0)
}

func fillArray(val reflect.Value, rng *rand.Rand) {
	for i := range val.Len() {
		fillValue(val.Index(i), rng)
	}
}

func fillString(val reflect.Value, rng *rand.Rand) {
	sz := rng.IntN(16)
	b := make([]byte, sz)
	for i := range sz {
		b[i] = byte(rng.IntN(256))
	}
	val.SetString(string(b))
}

func intKind(k reflect.Kind) bool {
	switch k {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Int:
		return true
	default:
		return false
	}
}

func uintKind(k reflect.Kind) bool {
	switch k {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uint:
		return true
	default:
		return false
	}
}

func floatKind(k reflect.Kind) bool {
	switch k {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func complexKind(k reflect.Kind) bool {
	switch k {
	case reflect.Complex64, reflect.Complex128:
		return true
	default:
		return false
	}
}
