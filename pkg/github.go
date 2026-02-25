package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/alex0ptr/toomani/business"
	"github.com/charmbracelet/log"
)

type GitHub struct {
	BaseURL string
	token   string
}

func NewGitHub(token string) *GitHub {
	return &GitHub{
		BaseURL: "https://api.github.com",
		token:   token,
	}
}

type GitHubRepoDto struct {
	Name     string   `json:"name"`
	FullName string   `json:"full_name"`
	SSHURL   string   `json:"ssh_url"`
	CloneURL string   `json:"clone_url"`
	Owner    OwnerDto `json:"owner"`
}

type OwnerDto struct {
	Login string `json:"login"`
}

func (r GitHubRepoDto) ToRepository(space business.Path) business.Repository {
	name := business.NewPath(r.Name)
	fullPath := business.NewPath(r.FullName)
	return business.Repository{
		Name:      name,
		FullPath:  fullPath,
		SpacePath: fullPath.TrimParent(space),
		SshUrl:    r.SSHURL,
		HttpUrl:   r.CloneURL,
	}
}

func (g *GitHub) BySpace(space business.Path) ([]business.Repository, error) {
	var repositories []business.Repository
	var orgErr, userErr error

	orgRepos, orgErr := g.fetchOrgRepositories(space)
	if orgErr == nil {
		repositories = append(repositories, orgRepos...)
		log.Debugf("Found %d repositories from organization %s", len(orgRepos), space.String())
	} else {
		log.Debugf("Failed to fetch org repositories for %s: %v", space.String(), orgErr)
	}

	userRepos, userErr := g.fetchUserRepositories(space)
	if userErr == nil {
		repositories = append(repositories, userRepos...)
		log.Debugf("Found %d repositories from user %s", len(userRepos), space.String())
	} else {
		log.Debugf("Failed to fetch user repositories for %s: %v", space.String(), userErr)
	}

	if len(repositories) == 0 {
		if orgErr != nil && userErr != nil {
			return nil, fmt.Errorf("no repositories found for '%s': %w", space.String(), errors.Join(
				fmt.Errorf("organization endpoint: %w", orgErr),
				fmt.Errorf("user endpoint: %w", userErr),
			))
		}
		return nil, fmt.Errorf("no repositories found for '%s' (tried both org and user endpoints but found 0 repos)", space.String())
	}

	return repositories, nil
}

func (g *GitHub) fetchOrgRepositories(space business.Path) ([]business.Repository, error) {
	baseURL := fmt.Sprintf("%s/orgs/%s/repos?per_page=100&archived=false", g.BaseURL, url.QueryEscape(space.String()))
	repos, err := g.fetchRepositories(baseURL, space)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch organization repositories for '%s': %w", space.String(), err)
	}
	return repos, nil
}

func (g *GitHub) fetchUserRepositories(space business.Path) ([]business.Repository, error) {
	baseURL := fmt.Sprintf("%s/users/%s/repos?per_page=100&archived=false", g.BaseURL, url.QueryEscape(space.String()))
	repos, err := g.fetchRepositories(baseURL, space)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user repositories for '%s': %w", space.String(), err)
	}
	return repos, nil
}

func (g *GitHub) fetchRepositories(baseURL string, space business.Path) ([]business.Repository, error) {
	var repositories []business.Repository
	client := &http.Client{}

	page := 1
	for {
		pageURL := fmt.Sprintf("%s&page=%d", baseURL, page)
		req, err := http.NewRequest("GET", pageURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create HTTP request for %s: %w", pageURL, err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("token %s", g.token))
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		log.Debugf("fetching %s...", pageURL)
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("HTTP request failed for %s: %w (check network connectivity and GitHub API access)", pageURL, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			switch resp.StatusCode {
			case http.StatusUnauthorized:
				return nil, fmt.Errorf("authentication failed (HTTP %d): check your GitHub token - %s", resp.StatusCode, string(body))
			case http.StatusForbidden:
				return nil, fmt.Errorf("access forbidden (HTTP %d): token may lack required permissions or rate limit exceeded - %s", resp.StatusCode, string(body))
			case http.StatusNotFound:
				return nil, fmt.Errorf("resource not found (HTTP %d): organization/user '%s' does not exist or is not accessible - %s", resp.StatusCode, space.String(), string(body))
			default:
				return nil, fmt.Errorf("GitHub API request failed (HTTP %d): %s - Response: %s", resp.StatusCode, resp.Status, string(body))
			}
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body from %s: %w", pageURL, err)
		}

		var pageRepos []GitHubRepoDto
		if err := json.Unmarshal(body, &pageRepos); err != nil {
			return nil, fmt.Errorf("failed to parse JSON response from %s: %w - Response: %s", pageURL, err, string(body))
		}

		if len(pageRepos) == 0 {
			break
		}

		for _, repo := range pageRepos {
			repositories = append(repositories, repo.ToRepository(space))
		}

		page++
	}

	return repositories, nil
}
