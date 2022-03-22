package cli

import (
	"log"
	"testing"
)

func Test(t *testing.T) {
	app := App{
		Commands: Commands{
			&Command{
				Name: "mod",
				Commands: Commands{
					&Command{
						Name:      "init",
						Arguments: Args{"name"},
						Action: func(c *Context) error {
							log.Println("go mod init " + c.Get("name"))
							return nil
						},
					},
				},
			},
		},
	}

	args := []string{"go", "mod", "init", "kek"}
	err := app.Run(args)
	log.Println(app.Context)

	if err != nil {
		log.Fatal(err)
	}
}
