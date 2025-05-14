package cmd

import (
	"github.com/alex0ptr/toomani/pkg"
	"sync"
)

var shellWriter = sync.OnceValue(pkg.NewShellWriter)
var ManiWriter = sync.OnceValue(pkg.NewManiWriter)
