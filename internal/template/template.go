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
	Partials map[string]string
}

type Option func(*client) error

func WithPartial(name, path string) Option {
	return func(c *client) error {
		partialRaw, err := LoadFile(path)
		if err != nil {
			return err
		}
		c.Partials[name] = string(partialRaw)
		raymond.RegisterPartial(name, string(partialRaw))
		return nil
	}
}

func WithPartials(partials map[string]string) Option {
	return func(c *client) error {
		for name, path := range partials {
			partialRaw, err := LoadFile(path)
			if err != nil {
				return err
			}
			c.Partials[name] = string(partialRaw)
			raymond.RegisterPartial(name, string(partialRaw))
		}

		return nil
	}
}

func WithHelper(name string, helper interface{}) Option {
	return func(c *client) error {
		raymond.RegisterHelper(name, helper)
		return nil
	}
}

func New(opts ...Option) (Client, error) {
	c := &client{
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
