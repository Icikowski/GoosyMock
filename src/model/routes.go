package model

import (
	"sort"
	"strings"
)

// Routes represents the configuration consisting of paths
// and route configurations (sets of method-specific responses
// for requests)
type Routes map[string]Route

// GetOrderedPaths sorts the paths of Routes in order to add
// them in correct order in chi.Router
func (r Routes) GetOrderedPaths() []string {
	paths := []string{}
	for path := range r {
		paths = append(paths, path)
	}

	sort.SliceStable(paths, func(i, j int) bool {
		return paths[i] < paths[j]
	})
	sort.SliceStable(paths, func(i, j int) bool {
		return strings.Count(paths[i], "/") < strings.Count(paths[j], "/")
	})
	sort.SliceStable(paths, func(i, j int) bool {
		return paths[i] == "/*"
	})

	return paths
}
