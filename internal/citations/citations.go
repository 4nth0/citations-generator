package citations

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Client struct {
	Path      string
	Citations []string
}

func New(path string) *Client {
	return &Client{
		Path: path,
	}
}

func (c *Client) LoadCitations() ([]string, error) {
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
