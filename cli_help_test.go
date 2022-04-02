package cli

import "testing"

func CreateAppTesting() *App {
	app := &App{
		Options: Options{
			&Option{
				Name:   "verbose",
				Prompt: "verbose",
				Short:  'v',
				IsFlag: true,
			},
			&Option{
				Name:   "name",
				Prompt: "name",
				Short:  'n',
			},
		},
		Commands: Commands{
			&Command{
				Name:      "test",
				Arguments: Args{"c", "d"},
				Action: func(c *Context) error {
					return nil
				},
			},
		},
		Arguments: Args{"a", "b"},
		Action: func(c *Context) error {
			return nil
		},
	}

	return app
}

func TestHelp(t *testing.T) {
	app := CreateAppTesting()

	err := app.Run([]string{"app", "--help"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestHelpSubcommand(t *testing.T) {
    app := CreateAppTesting()

    err := app.Run([]string{"app", "test", "--help"})
    if err != nil {
        t.Fatal(err)
    }
}
