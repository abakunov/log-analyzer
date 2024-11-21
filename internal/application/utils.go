package application

// IsURL checks if a given path is a URL.
func IsURL(path string) bool {
	return len(path) > 4 && (path[:4] == "http" || path[:5] == "https")
}
