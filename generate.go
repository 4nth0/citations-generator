package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/4nth0/citations-generator/internal/citations"
	"github.com/4nth0/citations-generator/internal/config"
	"github.com/4nth0/citations-generator/internal/generator"
	"github.com/4nth0/citations-generator/internal/template"
	cp "github.com/otiai10/copy"
)

const (
	exportPath = "./export"

	perpPage = 10
)

func main() {
	config, err := config.Load()
	if err != nil {
		panic(err)
	}

	citationsClient := citations.New(config.Generator.Source)
	citations, err := citationsClient.LoadCitations()
	if err != nil {
		fmt.Println(err)
		return
	}

	InitExportForlder()

	tpl, err := initTemplateEngine(config)
	if err != nil {
		panic(err)
	}

	gen := generator.New(
		config,
		tpl,
		citations,
		map[string]string{
			"detail":  config.Generator.Templates.Detail.Template,
			"listing": config.Generator.Templates.Listing.Template,
		},
		map[string]string{
			"index":   config.Generator.Templates.Index.Dest,
			"detail":  config.Generator.Templates.Detail.Dest,
			"listing": config.Generator.Templates.Listing.Dest,
		},
	)

	gen.GenerateDetailPages()
	gen.GenerateIndexPage(perpPage)
}

func initTemplateEngine(config *config.Config) (template.Client, error) {
	return template.New(
		template.WithPartials(config.Generator.Templates.Partials),
		template.WithHelper("pagePath", func(page, index int) string {
			return fmt.Sprintf(config.Generator.Paths.Detail, page*perpPage+index)
		}),
		template.WithHelper("relatedPagePath", func(index int) string {
			return fmt.Sprintf(config.Generator.Paths.Detail, index)
		}),
		template.WithHelper("pagination", func(page, pages int) string {
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
		}),
	)
}

func InitExportForlder() {
	os.RemoveAll(exportPath)
	os.Mkdir(exportPath, 0755)

	err := cp.Copy("./public", "./export")
	if err != nil {
		panic(err)
	}
}
