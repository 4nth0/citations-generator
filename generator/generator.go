package generator

import (
	"fmt"
	"os"

	"github.com/4nth0/citations-generator/citations"
	"github.com/4nth0/citations-generator/config"
)

type TemplateManager interface {
	UseLayout(path string) (func(map[string]interface{}) (string, error), error)
}

type Client struct {
	config    *config.Config
	tpl       TemplateManager
	Citations []citations.Citation
	Paths     map[string]string
	Layouts   map[string]string
}

func New(
	config *config.Config,
	tpl TemplateManager,
	citations []citations.Citation,
	layouts map[string]string,
	paths map[string]string,
) *Client {
	return &Client{
		config:    config,
		tpl:       tpl,
		Citations: citations,
		Paths:     paths,
		Layouts:   layouts,
	}
}

type RelatedCitation struct {
	Citation citations.Citation
	Index    int
}

type Page struct {
	Path    string
	Content string
}

type PagesTree map[string]Page

func (c Client) GeneratePages() PagesTree {
	export := PagesTree{}

	c.generateDetailsPages(export)

	return export
}

func (c Client) generateDetailsPages(pages PagesTree) {
	generate, err := c.tpl.UseLayout(c.Layouts["detail"])
	if err != nil {
		panic(err)
	}

	for idx, citation := range c.Citations {
		filePath := fmt.Sprintf(c.Paths["detail"], idx)

		tplCtx := map[string]interface{}{
			"citation": citation,
		}

		prevIndex := idx - 1
		if idx == 0 {
			prevIndex = len(c.Citations) - 1
		}
		tplCtx["citation_prev"] = RelatedCitation{
			Citation: c.Citations[prevIndex],
			Index:    prevIndex,
		}

		nextIndex := idx + 1
		if nextIndex >= len(c.Citations) {
			nextIndex = 0
		}

		tplCtx["citation_next"] = RelatedCitation{
			Citation: c.Citations[nextIndex],
			Index:    nextIndex,
		}

		tplCtx["og"] = map[string]string{
			"title":       c.config.Author.Name,
			"url":         fmt.Sprintf(c.config.Base+c.config.Generator.Paths.Detail, idx),
			"description": citation.Citation,
		}

		rendered, err := generate(tplCtx)
		if err != nil {
			return
		}

		pages[filePath] = Page{
			Content: rendered,
		}
	}
}

func (c Client) GenerateDetailPages() error {
	generate, err := c.tpl.UseLayout(c.Layouts["detail"])
	if err != nil {
		panic(err)
	}

	for idx, citation := range c.Citations {
		filePath := fmt.Sprintf(c.Paths["detail"], idx)

		tplCtx := map[string]interface{}{
			"citation": citation,
		}

		prevIndex := idx - 1
		if idx == 0 {
			prevIndex = len(c.Citations) - 1
		}
		tplCtx["citation_prev"] = RelatedCitation{
			Citation: c.Citations[prevIndex],
			Index:    prevIndex,
		}

		nextIndex := idx + 1
		if nextIndex >= len(c.Citations) {
			nextIndex = 0
		}

		tplCtx["citation_next"] = RelatedCitation{
			Citation: c.Citations[nextIndex],
			Index:    nextIndex,
		}

		tplCtx["og"] = map[string]string{
			"title":       c.config.Author.Name,
			"url":         fmt.Sprintf(c.config.Base+c.config.Generator.Paths.Detail, idx),
			"description": citation.Citation,
		}

		rendered, err := generate(tplCtx)
		if err != nil {
			return err
		}

		err = PutInFile(filePath, rendered)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c Client) GenerateIndexPage(perPage int) error {
	pages := len(c.Citations) / perPage
	generate, err := c.tpl.UseLayout(c.Layouts["listing"])
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(c.Citations); i += perPage {
		offset := i + perPage
		if offset > len(c.Citations) {
			offset = len(c.Citations)
		}
		citationsPerPage := c.Citations[i:offset]
		page := i / perPage

		generated, err := generate(map[string]interface{}{
			"citations": citationsPerPage,
			"page":      page,
			"pages":     pages,
		})
		if err != nil {
			panic(err)
		}

		var pathToSave string
		if i > 0 {
			pathToSave = fmt.Sprintf(c.Paths["listing"], page)
		} else {
			pathToSave = c.Paths["index"]
		}

		PutInFile(pathToSave, generated)
	}
	return nil
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
