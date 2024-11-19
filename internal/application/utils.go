package application

// isURL checks if a given path is a URL.
func isURL(path string) bool {
	return len(path) > 4 && (path[:4] == "http" || path[:5] == "https")
}
