package template

import (
	"io/ioutil"
	"os"

	"github.com/aymerick/raymond"
)

type Client struct {
	Partials map[string]string
}

type Option func(*Client) error

func WithPartial(name, path string) Option {
	return func(c *Client) error {
		partialRaw, err := LoadFile(path)
		if err != nil {
			return err
		}
		c.Partials[name] = string(partialRaw)
		raymond.RegisterPartial(name, string(partialRaw))
		return nil
	}
}

func WithHelper(name string, helper interface{}) Option {
	return func(c *Client) error {
		raymond.RegisterHelper(name, helper)
		return nil
	}
}

func New(opts ...Option) (*Client, error) {
	c := &Client{
		Partials: make(map[string]string),
	}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

type TemplateExec func(map[string]interface{}) (string, error)

func (c *Client) UseLayout(path string) (TemplateExec, error) {
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
