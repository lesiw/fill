# lesiw.io/fill

[![Go Reference](https://pkg.go.dev/badge/lesiw.io/fill.svg)](https://pkg.go.dev/lesiw.io/fill)

A utility for filling Go values.

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

func TestTlsConfigClone(f *testing.F) {
    opts := cmp.Options{cmpopts.IgnoreUnexported(tls.Config{})}
    for range 100 {
        var cfg tls.Config
        fill.Rand(&cfg)
        if want, got := cfg, cfg.Clone(); !cmp.Equal(want, got, opts) {
            t.Errorf("-original +cloned\n%s", cmp.Diff(want, got, opts))
        }
    }
}
```

[▶️ Run this example on the Go Playground](https://go.dev/play/p/PJCBSOH2VaO)
