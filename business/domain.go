package business

// Repository represents a Git repository with its metadata.
type Repository struct {
	// Name is the URL-friendly name of the repository, unique within the namespace.
	Name Path
	// FullPath is the full path of the repository, including the top-level namespace.
	FullPath Path
	// SpacePath is the relative path of the repository within its namespace.
	SpacePath Path
	// SshUrl is the SSH URL for accessing the repository.
	SshUrl string
	// HttpUrl is the HTTP URL for accessing the repository.
	HttpUrl string
}

// Repositories defines an interface for interacting with Git server repositories.
type Repositories interface {
	// BySpace retrieves all repositories within a given namespace (space).
	BySpace(path Path) ([]Repository, error)
}

// ConfigurationWriter defines an interface for writing configuration files for repository managers.
type ConfigurationWriter interface {
	// Write generates a configuration file content for the provided repositories.
	Write(repositories []Repository) string
}
