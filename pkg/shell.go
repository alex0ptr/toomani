package pkg

import (
	_ "embed"
	"github.com/charmbracelet/log"
	"strings"
	"text/template"
	"toomani/business"
)

//go:embed shell.sh.tmpl
var shellTemplate string

type ShellWriter struct {
	parsedTemplate *template.Template
}

type templateInput struct {
	Repositories []shellRepository
}

type shellRepository struct {
	Path    business.Path
	SshUrl  string
	HttpUrl string
}

func NewShellWriter() ShellWriter {
	parsedTemplate, err := template.New("shell.sh").Parse(shellTemplate)
	if err != nil {
		panic(err)
	}

	return ShellWriter{
		parsedTemplate: parsedTemplate,
	}

}

func mapToShell(repository business.Repository) shellRepository {
	return shellRepository{
		Path:    repository.SpacePath,
		SshUrl:  repository.SshUrl,
		HttpUrl: repository.HttpUrl,
	}

}

func (s ShellWriter) Write(repositories []business.Repository) string {
	shellRepositories := make([]shellRepository, 0, len(repositories))
	for _, repository := range repositories {
		shellRepositories = append(shellRepositories, mapToShell(repository))
	}
	input := templateInput{
		Repositories: shellRepositories,
	}
	builder := strings.Builder{}
	err := s.parsedTemplate.Execute(&builder, input)
	if err != nil {
		log.Fatal(err)
	}
	return builder.String()
}
