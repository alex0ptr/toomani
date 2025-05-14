package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"io"
	"net/http"
	url2 "net/url"
	"toomani/business"
)

type GitLab struct {
	BaseURL string
	token   string
}

func NewGitLab(baseURL, token string) *GitLab {
	return &GitLab{
		BaseURL: baseURL,
		token:   token,
	}
}

type ProjectDto struct {
	Path              string       `json:"Path"`
	PathWithNamespace string       `json:"path_with_namespace"` // GitLab-Path mit Namespace
	SSHURLToRepo      string       `json:"ssh_url_to_repo"`
	HTTPURLToRepo     string       `json:"http_url_to_repo"`
	Namespace         NamespaceDto `json:"namespace"`
}

type NamespaceDto struct {
	FullPath string `json:"full_path"`
}

func (p ProjectDto) ToRepository(space business.Path) business.Repository {
	path := business.NewPath(p.Path)
	fullPath := business.NewPath(p.Namespace.FullPath).Append(path)
	return business.Repository{
		Name:      path,
		FullPath:  fullPath,
		SpacePath: fullPath.TrimParent(space),
		SshUrl:    p.SSHURLToRepo,
		HttpUrl:   p.HTTPURLToRepo,
	}
}

func (g *GitLab) BySpace(space business.Path) ([]business.Repository, error) {
	url := fmt.Sprintf("%s/groups/%s/projects?per_page=100&archived=false&include_subgroups=true&with_shared=false", g.BaseURL, url2.QueryEscape(space.String()))
	var projects []business.Repository
	client := &http.Client{}

	for {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.token))
		log.Debugf("fetching %s...", url)
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("could not list projects: %s", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var pageProjects []ProjectDto
		if err := json.Unmarshal(body, &pageProjects); err != nil {
			return nil, err
		}
		for _, project := range pageProjects {
			projects = append(projects, project.ToRepository(space))
		}

		nextPage := resp.Header.Get("X-Next-Page")
		if nextPage != "" {
			url = fmt.Sprintf("%s/groups/%s/projects?per_page=100&archived=false&include_subgroups=true&with_shared=false&page=%s", g.BaseURL, url2.QueryEscape(string(space)), nextPage)
		} else {
			break
		}
	}

	return projects, nil
}
