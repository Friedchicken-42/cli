package cli

import (
	"errors"
	"testing"
)

func CreateApp() *App {
	app := &App{
		Commands: Commands{
			&Command{
				Name: "init",
				Action: func(c *Context) error {
					return nil
				},
			},
			&Command{
				Name:      "clone",
				Arguments: Args{"url"},
				Options: Options{
					&Option{
						Short:  'n',
						Prompt: "name",
						Name:   "NAME",
					},
				},
				Action: func(c *Context) error {
					v, ok := c.Get("NAME")

					if !ok {
						return errors.New("missing option NAME")
					}
					if v != "test" {
						return errors.New("wrong name")
					}
					return nil
				},
			},
		},
		Options: Options{
			&Option{
				Short:  'v',
				IsFlag: true,
			},
		},
		Action: func(c *Context) error {
			if _, ok := c.Get("v"); !ok {
				return errors.New("missing flag version")
			}
			return nil
		},
	}

	return app
}

func TestGeneric(t *testing.T) {
	app := CreateApp()

	err := app.Run([]string{"git", "clone", "--name", "test", "kek"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBase(t *testing.T) {
	app := CreateApp()

	err := app.Run([]string{"git", "init"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestMissing(t *testing.T) {
	app := CreateApp()

	err := app.Run([]string{"git", "clone"})
	if err == nil {
		t.Fatal("this should have failed")
	}

}

func TestGlobalFlags(t *testing.T) {
	app := CreateApp()

	err := app.Run([]string{"git", "-v"})
	if err != nil {
		t.Fatal(err)
	}
}
