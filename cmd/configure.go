package cmd

import (
	"gitlab-to-mani/pkg"
	"sync"
)

var shellWriter = sync.OnceValue(pkg.NewShellWriter)
var ManiWriter = sync.OnceValue(pkg.NewManiWriter)
