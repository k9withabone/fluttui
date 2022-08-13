package template

import "github.com/k9withabone/fluttui/tui/components/choose"

func New(choice string) (choose.Model, error) {
	items := []string{
		"app",
		"module",
		"package",
		"plugin",
		"plugin_ffi",
		"skeleton",
	}

	selected := 0
	for i, item := range items {
		if item == choice {
			selected = i
			break
		}
	}

	return choose.New(choose.Options{
		Title: "Flutter Template?",
		Items: items,
		Selected: []int{ selected },
		Limit: true,
	})
}