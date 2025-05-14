package cmd

import (
	"sync"
	"toomani/pkg"
)

var shellWriter = sync.OnceValue(pkg.NewShellWriter)
var ManiWriter = sync.OnceValue(pkg.NewManiWriter)
