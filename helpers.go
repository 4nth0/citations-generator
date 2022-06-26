package main

import (
	"fmt"
	"strings"

	"github.com/4nth0/citations-generator/internal/config"
)

func PagePathHelper(config *config.Config) func(page, index int) string {
	return func(page, index int) string {
		return fmt.Sprintf(config.Generator.Paths.Detail, page*config.Generator.CitationsPerPage+index)
	}
}

func RelatedPagePathHelper(config *config.Config) func(index int) string {
	return func(index int) string {
		return fmt.Sprintf(config.Generator.Paths.Detail, index)
	}
}

func PaginationHelper(config *config.Config) func(page, pages int) string {
	return func(page, pages int) string {
		links := []string{}

		for i := 0; i < pages+1; i++ {
			var path string
			if i == 0 {
				path = config.Generator.Paths.Index
			} else {
				path = fmt.Sprintf(config.Generator.Paths.Listing, i)
			}

			currentClass := ""
			if i == page {
				currentClass = "class='current'"
			}
			links = append(links, fmt.Sprintf(
				"<li %s><a href='%s'>%d</a></li>",
				currentClass,
				path,
				i+1,
			))
		}

		return fmt.Sprintf("<ul class='pagination'>%s</ul>", strings.Join(links, " "))
	}
}
