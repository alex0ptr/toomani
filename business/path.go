package business

import "strings"

// Path represents a repository path in the form "path/to/repo".
// It does not contain leading or trailing slashes, and an empty string is a valid path.
type Path string

const emptyPath = Path("")

// String converts the Path to its string representation.
func (path Path) String() string {
	return string(path)
}

// Append appends another Path to the current Path, ensuring proper formatting.
func (path Path) Append(toAppend Path) Path {
	// Using NewPath ensures sanitizing.
	return NewPath(path.String() + "/" + toAppend.String())
}

// NewPath creates a new Path instance, by removing leading and trailing slashes.
func NewPath(path string) Path {
	return Path(strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/"))
}

// TrimParent removes the parent Path prefix from the current Path.
func (path Path) TrimParent(parent Path) Path {
	// Using NewPath removes leading slashes.
	return NewPath(strings.TrimPrefix(path.String(), parent.String()))
}
