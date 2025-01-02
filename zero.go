package fill

var zeroFiller = Filler{NeverNil: true}

// Zero fills a value with zero-valued data.
func Zero(a any) { zeroFiller.Fill(a) }
