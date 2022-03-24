
# Cli

A package for handling command line arguments in Go

Still work in progress

## Related

Based on [urfave/cli](https://github.com/urfave/cli)

## Installation

```bash
  go get github.com/Friedchicken-42/cli
```
    
## Usage
Basic functionality
```Go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/Friedchicken-42/cli"
)

func main() {
    app := &cli.App{
        Action: func(c *cli.Context) error {
            log.Println("basic action")
            return nil
        },
    }

    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}
```

Complete example
```Go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/Friedchicken-42/cli"
)

func main() {
    app := &cli.App{
    	Commands:  cli.Commands{
            &cli.Command{
            	Name:      "sub1",
            	Arguments: cli.Args{"name"},
            	Action: func(c *cli.Context) error {
                    v, _ := c.Get("name")
                    fmt.Printf("command sub1 with arg: %s\n", v)
                    return nil
            	},
            },
            &cli.Command{
                Name: "sub2",
                Options: cli.Options{
                    &cli.Option{
                    	Name:   "add",
                    },
                    &cli.Option{
                        Prompt: "sub",
                        Short: 's',
                    },
                    &cli.Option{
                        Name: "verbose",
                        Short: 'v',
                        IsFlag: true,
                    },
                },
                Action: func(c *cli.Context) error {
                    fmt.Println(c.Get("add"))
                    fmt.Println(c.Get("sub"))
                    fmt.Println(c.Get("verbose"))
                    return nil
                },
            },
        },
    }

    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}

```
This will handle 2 subcommands
 - sub1 with a single required argument
 - sub2 with 3 optional flags
## TODO
- [ ]  generate help menu
