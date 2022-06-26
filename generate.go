package main

import (
	"fmt"
	"os"

	"github.com/4nth0/citations-generator/internal/citations"
	"github.com/4nth0/citations-generator/internal/config"
	"github.com/4nth0/citations-generator/internal/generator"
	"github.com/4nth0/citations-generator/internal/template"
	cp "github.com/otiai10/copy"
)

const (
	publicFolderPath = "./public"
	exportFolderPath = "./export"
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
	gen.GenerateIndexPage(config.Generator.CitationsPerPage)
}

func initTemplateEngine(config *config.Config) (template.Client, error) {
	return template.New(
		template.WithPartials(config.Generator.Templates.Partials),
		template.WithHelper("pagePath", PagePathHelper(config)),
		template.WithHelper("relatedPagePath", RelatedPagePathHelper(config)),
		template.WithHelper("pagination", PaginationHelper(config)),
	)
}

func InitExportForlder() {
	os.RemoveAll(exportFolderPath)
	os.Mkdir(exportFolderPath, 0755)

	err := cp.Copy(publicFolderPath, exportFolderPath)
	if err != nil {
		panic(err)
	}
}
