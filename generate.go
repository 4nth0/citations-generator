package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/4nth0/citations-generator/citations"
	"github.com/4nth0/citations-generator/config"
	"github.com/4nth0/citations-generator/generator"
	"github.com/4nth0/citations-generator/template"
	cp "github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
)

const (
	publicFolderPath = "./public"
	exportFolderPath = "./export"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Error(err)
		return
	}

	citationsClient := citations.New(config.Generator.Source)
	citations, err := citationsClient.LoadCitations()
	if err != nil {
		log.Error(err)
		return
	}

	InitExportForlder()

	tpl, err := initTemplateEngine(config)
	if err != nil {
		log.Error(err)
		return
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

	pages := gen.GeneratePages()
	fmt.Println(pages)

	gen.GenerateDetailPages()
	gen.GenerateIndexPage(config.Generator.CitationsPerPage)
}

func initTemplateEngine(config *config.Config) (template.Client, error) {
	partials, err := LoadPartials(config.Generator.Templates.Partials)
	if err != nil {
		return nil, err
	}

	return template.New(
		template.WithPartials(partials),
		template.WithHelpers(map[string]interface{}{
			"pagePath":        PagePathHelper(config),
			"relatedPagePath": RelatedPagePathHelper(config),
			"pagination":      PaginationHelper(config),
		}),
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

func LoadPartials(partials map[string]string) (map[string]string, error) {
	partialsMap := make(map[string]string)
	for name, path := range partials {
		partial, err := LoadFile(path)
		if err != nil {
			return nil, err
		}
		partialsMap[name] = string(partial)
	}
	return partialsMap, nil
}

func LoadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}
