package pkg

import (
	"github.com/charmbracelet/log"
	"gitlab-to-mani/business"
	"gopkg.in/yaml.v3"
)

type ManiWriter struct {
}

func NewManiWriter() ManiWriter {
	return ManiWriter{}
}

type ManiConfig struct {
	Projects map[string]ManiProject `yaml:"projects"`
	Tasks    map[string]Task        `yaml:"tasks"`
}
type ManiProject struct {
	Url  string `yaml:"url"`
	Path string `yaml:"path"`
}

type Task struct {
	Cmd         string `yaml:"cmd"`
	Description string `yaml:"description"`
}

var DefaultTasks = map[string]Task{
	"fetch-all": {
		Cmd:         "git fetch --all",
		Description: "Fetch all remotes.",
	},
	"pre-commit": {
		Cmd:         "pre-commit install --allow-missing-config",
		Description: "Install pre-commit hooks.",
	},
	"pull-main": {
		Cmd:         "git switch main && git pull",
		Description: "Switch to main and pull.",
	},
}

func (m ManiWriter) Write(repositories []business.Repository) string {
	maniConfig := ManiConfig{
		Projects: make(map[string]ManiProject),
	}
	for _, repo := range repositories {
		maniConfig.Projects[repo.Name.String()] = ManiProject{
			Url:  repo.SshUrl,
			Path: repo.SpacePath.String(),
		}
	}

	maniConfig.Tasks = DefaultTasks
	yamlBytes, err := yaml.Marshal(&maniConfig)
	if err != nil {
		log.Fatal(err)
	}

	return string(yamlBytes)
}
