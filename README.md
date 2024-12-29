# lesiw.io/fill

A utility for filling Go values with random data.

## Example

``` go
package main

import (
    "crypto/tls"
    "math/rand/v2"
    "testing"

    "github.com/google/go-cmp/cmp"
    "github.com/google/go-cmp/cmp/cmpopts"
    "lesiw.io/fill"
)

func FuzzTlsConfigClone(f *testing.F) {
    opts := cmp.Options{cmpopts.IgnoreUnexported(tls.Config{})}
    f.Fuzz(func(t *testing.T, seed1 uint64, seed2 uint64) {
        cfg1 := &tls.Config{}
        fill.Do(&cfg1, rand.New(rand.NewPCG(seed1, seed2)))
        cfg2 := cfg1.Clone()
        if !cmp.Equal(cfg1, cfg2, opts) {
            t.Errorf("-original +cloned\n%s", cmp.Diff(cfg1, cfg2, opts))
        }
    })
}
```

[▶️ Run this example on the Go Playground](https://go.dev/play/p/_HX3J0f__AF)

To run locally, `go test -fuzz=Fuzz -fuzztime=10s`.
