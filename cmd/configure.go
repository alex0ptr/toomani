package cmd

import (
	"sync"
	"github.com/alex0ptr/toomani/pkg"
)

var shellWriter = sync.OnceValue(pkg.NewShellWriter)
var ManiWriter = sync.OnceValue(pkg.NewManiWriter)
