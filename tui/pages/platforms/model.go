package platforms

import "github.com/k9withabone/fluttui/tui/components/choose"

func New(choices []string) (choose.Model, error) {
	items := []string{
		"ios",
		"android",
		"windows",
		"linux",
		"macos",
		"web",
	}

	selected := []int{}
	choicesLen := len(choices)
	for i, item := range items {
		if choicesLen > 0 {
			for _, choice := range choices {
				if choice == item {
					selected = append(selected, i)
					break
				}
			}
		} else {
			// default to all selected
			selected = append(selected, i)
		}
	}

	return choose.New(choose.Options{
		Title: "Flutter Platforms?",
		Items: items,
		Selected: selected,
		Limit: false,
	})
}