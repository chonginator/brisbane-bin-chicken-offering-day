package resource

import "strings"

type Resource struct {
	Name string
	Slug string
}

func FilterByName(resources []Resource, nameSubstring string) []Resource {
	filtered := make([]Resource, 0)
	for _, resource := range resources {
		if strings.Contains(strings.ToLower(resource.Name), strings.ToLower(nameSubstring)) {
			filtered = append(filtered, resource)
		}
	}

	return filtered
}
