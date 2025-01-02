package fill

import (
	"fmt"
	"reflect"
	"strings"
	// math/rand/v2 should not be imported in this file.
)

var randFiller = Filler{MaxSize: 8, Runes: Base64}

// Rand fills a value with random data.
func Rand(a any) { randFiller.Fill(a) }

func (f *Filler) Fill(a any) {
	val := reflect.ValueOf(a)
	if val.Kind() != reflect.Pointer {
		panic("bad parameter: value to fill must be pointer")
	}
	f.fillValue(val)
}

func (f *Filler) fillValue(val reflect.Value) {
	k := val.Kind()
	switch {
	case k == reflect.Bool:
		val.SetBool(f.intN(2) == 0)
	case intKind(k):
		val.SetInt(f.int64())
	case uintKind(k):
		val.SetUint(f.uint64())
	case floatKind(k):
		val.SetFloat(f.float64())
	case complexKind(k):
		val.SetComplex(complex(f.float64(), f.float64()))
	case k == reflect.Array:
		f.fillArray(val)
	case k == reflect.Chan:
		f.fillChan(val)
	case k == reflect.Func || k == reflect.Interface:
		// Can't fill.
	case k == reflect.Map:
		f.fillMap(val)
	case k == reflect.Pointer:
		f.fillPointer(val)
	case k == reflect.Slice:
		f.fillSlice(val)
	case k == reflect.String:
		f.fillString(val)
	case k == reflect.Struct:
		f.fillStruct(val)
	default:
		panic(fmt.Sprintf("unhandled type: %s", val.Kind()))
	}
}

func (f *Filler) fillArray(val reflect.Value) {
	for i := range val.Len() {
		f.fillValue(val.Index(i))
	}
}

func (f *Filler) fillChan(val reflect.Value) {
	if sz := f.intN(f.MaxSize - f.MinSize); sz == 0 && !f.NeverNil {
		val.SetZero() // nil
	} else {
		val.Set(reflect.MakeChan(val.Type(),
			f.MinSize+f.intN(f.MaxSize-f.MinSize)))
	}
}

func (f *Filler) fillMap(val reflect.Value) {
	sz := f.intN(f.MaxSize - f.MinSize)
	if !f.NeverNil && sz == 0 {
		val.SetZero() // nil
		return
	}
	sz = f.MinSize + f.intN(f.MaxSize-f.MinSize)
	if val.IsZero() {
		val.Set(reflect.MakeMapWithSize(val.Type(), sz))
	} else {
		val.Clear()
	}
	for range sz {
		k := reflect.New(val.Type().Key()).Elem()
		f.fillValue(k)
		v := reflect.New(val.Type().Elem()).Elem()
		f.fillValue(v)
		val.SetMapIndex(k, v)
	}
}

func (f *Filler) fillPointer(val reflect.Value) {
	if !f.NeverNil && f.intN(f.MaxSize-f.MinSize) == 0 {
		val.SetZero() // nil
	} else {
		if val.IsZero() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		f.fillValue(val.Elem())
	}
}

func (f *Filler) fillSlice(val reflect.Value) {
	sz := f.intN(f.MaxSize - f.MinSize)
	if !f.NeverNil && sz == 0 {
		val.SetZero() // nil
		return
	}
	sz = f.MinSize + f.intN(f.MaxSize-f.MinSize)
	for range sz {
		e := reflect.New(val.Type().Elem()).Elem()
		f.fillValue(e)
		val.Set(reflect.Append(val, e))
	}
}

func (f *Filler) fillString(val reflect.Value) {
	sz := f.MinSize + f.intN(f.MaxSize-f.MinSize)
	var b strings.Builder
	if f.Runes == nil {
		for range sz {
			b.WriteByte(byte(f.intN(256)))
		}
	} else {
		for range sz {
			b.WriteRune(f.Runes[f.intN(len(f.Runes))])
		}
	}
	val.SetString(b.String())
}

func (f *Filler) fillStruct(val reflect.Value) {
	for i := range val.NumField() {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}
		f.fillValue(field)
	}
}
