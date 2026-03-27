package patch

import (
	"bufio"
	"io"
	"strings"
)

// FileStat holds the change statistics for a single file.
type FileStat struct {
	Path    string
	Added   int
	Deleted int
}

// Parse reads a unified diff from r and returns per-file change statistics.
func Parse(r io.Reader) ([]FileStat, error) {
	stats := make(map[string]*FileStat)
	var order []string
	var current *FileStat

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "diff --git ") {
			path := extractPath(line)
			if _, exists := stats[path]; !exists {
				stats[path] = &FileStat{Path: path}
				order = append(order, path)
			}
			current = stats[path]
			continue
		}

		if current == nil {
			continue
		}

		if strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---") {
			continue
		}

		if strings.HasPrefix(line, "+") {
			current.Added++
		} else if strings.HasPrefix(line, "-") {
			current.Deleted++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	result := make([]FileStat, len(order))
	for i, path := range order {
		result[i] = *stats[path]
	}
	return result, nil
}

// extractPath extracts the destination file path from a "diff --git a/... b/..." line.
func extractPath(line string) string {
	// Find the last occurrence of " b/" to handle paths with spaces
	idx := strings.LastIndex(line, " b/")
	if idx == -1 {
		return line
	}
	return line[idx+3:]
}
