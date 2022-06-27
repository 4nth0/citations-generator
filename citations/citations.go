package citations

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Citation struct {
	Citation      string `json:"citation"`
	Author        string `json:"author"`
	AuthorPicture string `json:"author_picture"`
}

type Client interface {
	LoadCitations() ([]Citation, error)
}

type client struct {
	Path      string
	Citations []Citation
}

func New(path string) Client {
	return &client{
		Path: path,
	}
}

func (c *client) LoadCitations() ([]Citation, error) {
	raw, err := LoadFile(c.Path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(raw, &c.Citations)
	if err != nil {
		return nil, err
	}

	return c.Citations, nil
}

func LoadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}
