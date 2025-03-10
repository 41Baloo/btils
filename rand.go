package btils

import (
	_ "unsafe"
)

// This isnt encouraged but it's extremely fast and somewhat widely used by other packages already

//go:linkname Fastrand runtime.cheaprand
func Fastrand() uint32
