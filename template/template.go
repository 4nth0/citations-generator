package template

import (
	"io/ioutil"
	"os"

	"github.com/aymerick/raymond"
)

type Client interface {
	UseLayout(path string) (func(map[string]interface{}) (string, error), error)
}

type client struct {
}

func New(opts ...Option) (Client, error) {
	c := &client{}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *client) UseLayout(path string) (func(map[string]interface{}) (string, error), error) {
	raw, err := LoadFile(path)
	if err != nil {
		return nil, err
	}

	tpl, err := raymond.Parse(string(raw))
	if err != nil {
		return nil, err
	}

	return func(templateContext map[string]interface{}) (string, error) {
		return tpl.Exec(templateContext)
	}, nil
}

func LoadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}
