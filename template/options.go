package template

import "github.com/aymerick/raymond"

type Option func(*client) error

// Ici nous devriosn directement avoir le contenu du fichier et non son adresse, cette partie n'a pas a connaitre le FS
func WithPartial(name, path string) Option {
	return func(c *client) error {
		return RegisterPartialFromFile(name, path)
	}
}

func WithPartials(partials map[string]string) Option {
	return func(c *client) error {
		for name, content := range partials {
			raymond.RegisterPartial(name, content)
		}

		return nil
	}
}

func RegisterPartialFromFile(name, path string) error {
	partialRaw, err := LoadFile(path)
	if err != nil {
		return err
	}
	raymond.RegisterPartial(name, string(partialRaw))

	return nil
}

func WithHelper(name string, helper interface{}) Option {
	return func(c *client) error {
		raymond.RegisterHelper(name, helper)
		return nil
	}
}

func WithHelpers(helpers map[string]interface{}) Option {
	return func(c *client) error {
		for name, helper := range helpers {
			raymond.RegisterHelper(name, helper)
		}
		return nil
	}
}
