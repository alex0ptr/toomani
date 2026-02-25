package cmd

import (
	"github.com/alex0ptr/toomani/business"
	"github.com/alex0ptr/toomani/pkg"
	"sync"
)

var shellWriter = sync.OnceValue(pkg.NewShellWriter)
var maniWriter = sync.OnceValue(pkg.NewManiWriter)

func newWriter(output string) business.ConfigurationWriter {
	switch output {
	case "shell":
		return shellWriter()
	default:
		return maniWriter()
	}
}
