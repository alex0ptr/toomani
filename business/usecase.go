package business

import (
	"github.com/charmbracelet/log"
	"slices"
	"strings"
)

// GenerateRepositoryListing is a usecase for generating a repository listing based on a repository provider
type GenerateRepositoryListing struct {
	Repositories        Repositories
	ConfigurationWriter ConfigurationWriter
}

// NewGenerateRepositoryListing creates a new instance of GenerateRepositoryListing.
func NewGenerateRepositoryListing(repositories Repositories, configurationWriter ConfigurationWriter) *GenerateRepositoryListing {
	return &GenerateRepositoryListing{
		Repositories:        repositories,
		ConfigurationWriter: configurationWriter,
	}
}

// RepositoriesBySpace retrieves repositories within a given namespace (spacePath),
// optionally filtering them by match and exclude prefixes.
func (g *GenerateRepositoryListing) RepositoriesBySpace(spacePath Path, matchPrefix []string, excludePrefix []string) ([]Repository, error) {
	log.Debugf("searching space %s for repositories...", spacePath)
	repositories, err := g.Repositories.BySpace(spacePath)
	if err != nil {
		return nil, err
	}
	log.Debugf("found %d repositories", len(repositories))

	if len(matchPrefix) > 0 {
		log.Debugf("keep repositories with prefix: %s", strings.Join(matchPrefix, ", "))
		repositories = filtered(repositories, matchPrefix, true)
	}

	if len(excludePrefix) > 0 {
		log.Debugf("exclude repositories with prefix: %s", strings.Join(excludePrefix, ", "))
		repositories = filtered(repositories, excludePrefix, false)
	}

	return repositories, nil
}

func filtered(repositories []Repository, matchPrefix []string, mode bool) (filtered []Repository) {
	filtered = slices.DeleteFunc(repositories, func(r Repository) bool {
		for _, prefix := range matchPrefix {
			if strings.HasPrefix(r.FullPath.String(), prefix) {
				return !mode
			}
		}
		return mode
	})
	log.Debugf("removed %d repositories", len(repositories)-len(filtered))
	return
}

// Write generates a configuration file content for the provided repositories.
func (g *GenerateRepositoryListing) Write(repositories []Repository) string {
	return g.ConfigurationWriter.Write(repositories)
}

// WriteManagementFile generates and writes a configuration file for repositories
// within a given namespace (spacePath), applying optional filters.
func (g *GenerateRepositoryListing) WriteManagementFile(spacePath Path, matchPrefix []string, excludePrefix []string) (string, error) {
	repositories, err := g.RepositoriesBySpace(spacePath, matchPrefix, excludePrefix)
	if err != nil {
		return "", err
	}
	return g.ConfigurationWriter.Write(repositories), nil
}
