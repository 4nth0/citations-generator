package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/4nth0/citations-generator/internal/citations"
	"github.com/4nth0/citations-generator/internal/generator"
	"github.com/4nth0/citations-generator/internal/template"
	cp "github.com/otiai10/copy"
)

const (
	exportPath = "./export"

	perpPage = 10
)

func main() {
	citationsClient := citations.New("./citations.json")
	citations, err := citationsClient.LoadCitations()
	if err != nil {
		fmt.Println(err)
		return
	}

	InitExportForlder()

	tpl, err := initTemplateEngine()
	if err != nil {
		panic(err)
	}

	gen := generator.New(
		tpl,
		citations,
		map[string]string{
			"detail":  "./layouts/detail.hbs",
			"listing": "./layouts/listing.hbs",
		},
		map[string]string{
			"index":   exportPath + "/index.html",
			"detail":  exportPath + "/citation-%d.html",
			"listing": exportPath + "/index-%d.html",
		},
	)

	gen.GenerateDetailPages()
	gen.GenerateIndexPage(perpPage)
}

func initTemplateEngine() (*template.Client, error) {
	return template.New(
		template.WithPartial("header", "./layouts/partials/header.hbs"),
		template.WithPartial("footer", "./layouts/partials/footer.hbs"),
		template.WithHelper("pagePath", func(page, index int) string {
			return fmt.Sprintf("citation-%d.html", page*perpPage+index)
		}),
		template.WithHelper("relatedPagePath", func(index int) string {
			return fmt.Sprintf("citation-%d.html", index)
		}),
		template.WithHelper("pagination", func(page, pages int) string {
			links := []string{}

			for i := 0; i < pages+1; i++ {
				var path string
				if i == 0 {
					path = "index.html"
				} else {
					path = fmt.Sprintf("index-%d.html", i)
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

// Init Export Folder that delete container files to avoid conflict
func InitExportForlder() {
	os.RemoveAll(exportPath)
	os.Mkdir(exportPath, 0755)

	err := cp.Copy("./public", "./export")
	if err != nil {
		panic(err)
	}
}
