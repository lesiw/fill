package fill

import (
	"fmt"
	"math/rand/v2"
	"reflect"
)

// maxSize is an arbitrary upper bound on value length.
const maxSize = 16

// Rand fills a value with random data.
func Rand(a any, rng *rand.Rand) {
	if rng == nil {
		panic("bad parameter: nil rand.Rand")
	}
	val := reflect.ValueOf(a)
	if val.Kind() != reflect.Pointer {
		panic("bad parameter: value to fill must be pointer")
	}
	randValue(val.Elem(), rng)
}

func randValue(val reflect.Value, rng *rand.Rand) {
	k := val.Kind()
	switch {
	case k == reflect.Bool:
		val.SetBool(rng.IntN(2) == 0)
	case intKind(k):
		val.SetInt(rng.Int64())
	case uintKind(k):
		val.SetUint(rng.Uint64())
	case floatKind(k):
		val.SetFloat(rng.Float64())
	case complexKind(k):
		val.SetComplex(complex(rng.Float64(), rng.Float64()))
	case k == reflect.Array:
		randArray(val, rng)
	case k == reflect.Chan:
		randChan(val, rng)
	case k == reflect.Func || k == reflect.Interface:
		// Can't fill.
	case k == reflect.Map:
		randMap(val, rng)
	case k == reflect.Pointer:
		randPointer(val, rng)
	case k == reflect.Slice:
		randSlice(val, rng)
	case k == reflect.String:
		randString(val, rng)
	case k == reflect.Struct:
		randStruct(val, rng)
	default:
		panic(fmt.Sprintf("unhandled type: %s", val.Kind()))
	}
}

func randArray(val reflect.Value, rng *rand.Rand) {
	for i := range val.Len() {
		randValue(val.Index(i), rng)
	}
}

func randChan(val reflect.Value, rng *rand.Rand) {
	if sz := rng.IntN(maxSize); sz == 0 {
		val.SetZero() // nil
	} else {
		val.Set(reflect.MakeChan(val.Type(), sz-1))
	}
}

func randMap(val reflect.Value, rng *rand.Rand) {
	sz := rng.IntN(maxSize)
	if sz == 0 {
		val.SetZero() // nil
		return
	}
	sz-- // Allow zero size.
	if val.IsZero() {
		val.Set(reflect.MakeMapWithSize(val.Type(), sz))
	} else {
		val.Clear()
	}
	for range sz {
		k := reflect.New(val.Type().Key()).Elem()
		randValue(k, rng)
		v := reflect.New(val.Type().Elem()).Elem()
		randValue(v, rng)
		val.SetMapIndex(k, v)
	}
}

func randPointer(val reflect.Value, rng *rand.Rand) {
	if rng.IntN(maxSize) == 0 {
		val.SetZero() // nil
	} else {
		if val.IsZero() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		randValue(val.Elem(), rng)
	}
}

func randSlice(val reflect.Value, rng *rand.Rand) {
	sz := rng.IntN(maxSize)
	if sz == 0 {
		val.SetZero() // nil
		return
	}
	sz-- // Allow zero size.
	for range sz {
		e := reflect.New(val.Type().Elem()).Elem()
		randValue(e, rng)
		val.Set(reflect.Append(val, e))
	}
}

func randString(val reflect.Value, rng *rand.Rand) {
	sz := rng.IntN(maxSize)
	b := make([]byte, sz)
	for i := range sz {
		b[i] = byte(rng.IntN(256))
	}
	val.SetString(string(b))
}

func randStruct(val reflect.Value, rng *rand.Rand) {
	for i := range val.NumField() {
		f := val.Field(i)
		if !f.CanSet() {
			continue
		}
		randValue(f, rng)
	}
}
