package main

import (
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

	err = InitExportForlder()
	if err != nil {
		panic(err)
	}

	tpl, err := initTemplateEngine(config)
	if err != nil {
		log.Error(err)
		return
	}

	gen := generator.Client{
		Config:    config,
		TPL:       tpl,
		Citations: citations,
		Paths: map[string]string{
			"index":   config.Generator.Templates.Index.Dest,
			"detail":  config.Generator.Templates.Detail.Dest,
			"listing": config.Generator.Templates.Listing.Dest,
		},
		Layouts: map[string]string{
			"detail":  config.Generator.Templates.Detail.Template,
			"listing": config.Generator.Templates.Listing.Template,
		},
		PerPage: config.Generator.CitationsPerPage,
	}

	pages := gen.GeneratePages()
	GenerateFiles(pages)
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

func InitExportForlder() error {
	var err error
	err = os.RemoveAll(exportFolderPath)
	if err != nil {
		return err
	}

	err = os.Mkdir(exportFolderPath, 0755)
	if err != nil {
		return err
	}

	return cp.Copy(publicFolderPath, exportFolderPath)
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

func GenerateFiles(pages generator.PagesTree) {
	for path, page := range pages {
		PutInFile(path, page.Content)
	}
}

func PutInFile(filePath, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}
