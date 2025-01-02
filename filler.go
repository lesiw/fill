package fill

import "math/rand/v2"

// A Filler fills in random values.
type Filler struct {
	MinSize    int        // Arbitrary lower bound for value sizes.
	MaxSize    int        // Arbitrary upper bound for value sizes.
	Runes      []rune     // Runes to select from when filling strings.
	NeverNil   bool       // Never select nil for fillable values.
	RandSource *rand.Rand // Custom rand source.
}

var Base64 = []rune{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
	'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd',
	'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's',
	't', 'u', 'v', 'w', 'x', 'y', 'z', '0', '1', '2', '3', '4', '5', '6', '7',
	'8', '9', '+', '/',
}

func (f *Filler) intN(n int) int {
	if n == 0 {
		return 0
	} else if f.RandSource == nil {
		return rand.IntN(n)
	} else {
		return f.RandSource.IntN(n)
	}
}

func (f *Filler) int64() int64 {
	if f.RandSource == nil {
		return rand.Int64()
	} else {
		return f.RandSource.Int64()
	}
}

func (f *Filler) uint64() uint64 {
	if f.RandSource == nil {
		return rand.Uint64()
	} else {
		return f.RandSource.Uint64()
	}
}

func (f *Filler) float64() float64 {
	if f.RandSource == nil {
		return rand.Float64()
	} else {
		return f.RandSource.Float64()
	}
}
