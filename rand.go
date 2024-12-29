package fill

import (
	"fmt"
	"math/rand/v2"
	"reflect"
)

// Rand fills a value with random data.
func Rand(a any, rng *rand.Rand) {
	if rng == nil {
		panic("bad parameter: nil rand.Rand")
	}
	val := reflect.ValueOf(a)
	if val.Kind() != reflect.Ptr {
		panic("bad parameter: value to fill must be pointer")
	}
	fillValueRand(val.Elem(), rng)
}

func fillValueRand(val reflect.Value, rng *rand.Rand) {
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
		fillStringRand(val, rng)
	case k == reflect.Bool:
		val.SetBool(rng.IntN(2) == 0)
	case k == reflect.Struct:
		fillStructRand(val, rng)
	case k == reflect.Slice:
		fillSliceRand(val, rng)
	case k == reflect.Map:
		fillMapRand(val, rng)
	case k == reflect.Array:
		fillArrayRand(val, rng)
	case k == reflect.Chan:
		fillChanRand(val, rng)
	case k == reflect.Pointer:
		val.Set(reflect.New(val.Type().Elem()))
		fillValueRand(val.Elem(), rng)
	case k == reflect.Interface || k == reflect.Func:
		// Can't fill.
	default:
		panic(fmt.Sprintf("unhandled type: %s", val.Kind()))
	}
}

func fillStructRand(val reflect.Value, rng *rand.Rand) {
	for i := range val.NumField() {
		f := val.Field(i)
		if !f.CanSet() {
			continue
		}
		fillValueRand(f, rng)
	}
}

func fillSliceRand(val reflect.Value, rng *rand.Rand) {
	for range rng.IntN(16) {
		e := reflect.New(val.Type().Elem()).Elem()
		fillValueRand(e, rng)
		val.Set(reflect.Append(val, e))
	}
}

func fillMapRand(val reflect.Value, rng *rand.Rand) {
	sz := rng.IntN(16)
	val.Set(reflect.MakeMapWithSize(val.Type(), sz))
	for range sz {
		k := reflect.New(val.Type().Key()).Elem()
		fillValueRand(k, rng)
		v := reflect.New(val.Type().Elem()).Elem()
		fillValueRand(v, rng)
		val.SetMapIndex(k, v)
	}
}

func fillChanRand(val reflect.Value, _ *rand.Rand) {
	reflect.MakeChan(val.Type(), 0)
}

func fillArrayRand(val reflect.Value, rng *rand.Rand) {
	for i := range val.Len() {
		fillValueRand(val.Index(i), rng)
	}
}

func fillStringRand(val reflect.Value, rng *rand.Rand) {
	sz := rng.IntN(16)
	b := make([]byte, sz)
	for i := range sz {
		b[i] = byte(rng.IntN(256))
	}
	val.SetString(string(b))
}
