package pkg

import (
	_ "embed"
	"github.com/alex0ptr/toomani/business"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed shell-fixture-1.sh
var fixture1 string

func TestShellWriter_Write(t *testing.T) {
	writer := NewShellWriter()
	repositories := []business.Repository{
		{
			Name:      "test",
			FullPath:  "test",
			SpacePath: "test-space-1",
			SshUrl:    "test-space-ssh-1",
			HttpUrl:   "test-space-http-1",
		},
		{
			Name:      "test",
			FullPath:  "test",
			SpacePath: "test-space-2",
			SshUrl:    "test-space-ssh-2",
			HttpUrl:   "test-space-http-2",
		},
	}

	assert.Equal(t, fixture1, writer.Write(repositories))
}
