package generator

import (
	"fmt"
	"strings"

	"github.com/4nth0/citations-generator/citations"
	"github.com/4nth0/citations-generator/config"
)

type TemplateManager interface {
	UseLayout(path string) (func(map[string]interface{}) (string, error), error)
}

type Client struct {
	Config    *config.Config
	TPL       TemplateManager
	Citations []citations.Citation
	Paths     map[string]string
	Layouts   map[string]string
	PerPage   int
}

type RelatedCitation struct {
	Citation citations.Citation
	Index    int
}

type Page struct {
	Path    string
	Content string
	Context map[string]interface{}
}

type PagesTree map[string]Page

func (c Client) GeneratePages() PagesTree {
	export := PagesTree{}

	c.generateDetailsPages(export)
	c.generateIndexPages(export)
	c.generateSiteMap(export)

	return export
}

func (c Client) generateSiteMap(pages PagesTree) {
	pagesPaths := []string{}
	for _, page := range pages {
		pagesPaths = append(pagesPaths, page.Path)
	}

	pages["./export/sitemap.txt"] = Page{
		Content: strings.Join(pagesPaths, "\n"),
	}
}

func (c Client) generateDetailsPages(pages PagesTree) {
	generate, err := c.TPL.UseLayout(c.Layouts["detail"])
	if err != nil {
		panic(err)
	}

	for idx, citation := range c.Citations {
		filePath := fmt.Sprintf(c.Paths["detail"], idx)
		absolutePath := fmt.Sprintf(c.Config.Base+c.Config.Generator.Paths.Detail, idx)
		page := Page{
			Path:    absolutePath,
			Context: map[string]interface{}{},
		}

		page.Context["citation"] = citation

		prevIndex := idx - 1
		if idx == 0 {
			prevIndex = len(c.Citations) - 1
		}
		page.Context["citation_prev"] = RelatedCitation{
			Citation: c.Citations[prevIndex],
			Index:    prevIndex,
		}

		nextIndex := idx + 1
		if nextIndex >= len(c.Citations) {
			nextIndex = 0
		}

		page.Context["citation_next"] = RelatedCitation{
			Citation: c.Citations[nextIndex],
			Index:    nextIndex,
		}

		page.Context["og"] = map[string]string{
			"title":       c.Config.Author.Name,
			"url":         absolutePath,
			"description": citation.Citation,
		}

		rendered, err := generate(page.Context)
		if err != nil {
			return
		}

		page.Content = rendered

		pages[filePath] = page
	}
}

func (c Client) generateIndexPages(pages PagesTree) {
	pagesCount := len(c.Citations) / c.PerPage
	generate, err := c.TPL.UseLayout(c.Layouts["listing"])
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(c.Citations); i += c.PerPage {

		page := i / c.PerPage
		offset := i + c.PerPage
		if offset > len(c.Citations) {
			offset = len(c.Citations)
		}
		citations := c.Citations[i:offset]

		var absolutePath string
		var exportPath string
		if i > 0 {
			exportPath = fmt.Sprintf(c.Paths["listing"], page)
			absolutePath = fmt.Sprintf(c.Config.Base+c.Config.Generator.Paths.Listing, page)
		} else {
			exportPath = c.Paths["index"]
			absolutePath = c.Config.Base + c.Config.Generator.Paths.Listing
		}

		export := Page{
			Context: map[string]interface{}{},
			Path:    absolutePath,
		}

		export.Context = map[string]interface{}{
			"citations": citations,
			"page":      page,
			"pages":     pagesCount,
		}

		generated, err := generate(export.Context)
		if err != nil {
			panic(err)
		}

		export.Content = generated

		pages[exportPath] = export
	}
}
