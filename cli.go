package cli

import (
	"errors"
	"fmt"
)

type Context struct {
	Args  map[string]string
	Flags map[string]bool
}

type Args []string

type Option struct {
	Name   string
	Prompt string
	Short  rune
	IsFlag bool
}

type Options []*Option

type Command struct {
	Name      string
	Commands  Commands
	Arguments Args
	Options   Options
	context   *Context
	Action    func(c *Context) error
}

type Commands []*Command

type App Command

func (c *Context) Get(name string) (string, bool) {
	if v, ok := c.Args[name]; ok {
		return v, ok
	}

	if ok := c.Flags[name]; ok {
		return "", ok
	}

	return "", false
}

func (c *Context) size() int {
    return len(c.Args) + len(c.Flags)
}

func (c *Command) shareContext(context *Context) {
	c.context = context
	for i := range c.Commands {
		c.Commands[i].shareContext(context)
	}
}

func (c *Command) setArg(key, value string) {
	c.context.Args[key] = value
}

func (c *Command) findOption(name string) *Option {
    if name == "help" || name == "h" {
        return &Option{
        	Name:   "help",
        	Prompt: "help",
        	Short:  'h',
        	IsFlag: true,
        }
    }

	for i := range c.Options {
		if c.Options[i].Prompt == name {
			return c.Options[i]
		}
		if c.Options[i].Short == rune(name[0]) {
			return c.Options[i]
		}
		if c.Options[i].Name == name {
			return c.Options[i]
		}
	}
	return nil
}

func (c *Command) getOption(name string) *Option {
	if name[1] == '-' {
		return c.findOption(name[2:])
	}
	return c.findOption(name[1:])
}

func (c *Command) setOption(name, value string) {
	c.context.Args[name] = value
}

func (c *Command) setFlag(name string) {
	c.context.Flags[name] = true
}

func (c *Command) setArgs(args []string) error {
	a := 0

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg[0] == '-' {
			option := c.getOption(arg)

			if option == nil {
				return errors.New("missing option: " + arg)
			}

			var name string
			if option.Short != 0 {
				name = string(option.Short)
			}
			if option.Prompt != "" {
				name = option.Prompt
			}
			if option.Name != "" {
				name = option.Name
			}

			if option.IsFlag == false {
				i++
				value := args[i]

				c.setOption(name, value)
			} else {
				c.setFlag(name)
			}

		} else {
			if a == len(c.Arguments) {
				break
			}
			c.setArg(c.Arguments[a], arg)
			a++
		}
	}

	if len(c.Arguments) != a {
        if _, h := c.context.Get("help"); !h {
            return errors.New("mismatcing arguments")
        }
	}

	return nil
}

func (c Command) get(name string) (*Command, error) {
	for i := range c.Commands {
		if c.Commands[i].Name == name {
			return c.Commands[i], nil
		}
	}
	return nil, errors.New("command not found: " + name)
}

func (c *Command) search(args []string) (*Command, []string, error) {
	if len(args) == 0 {
		return c, nil, nil
	}

	command, err := c.get(args[0])
	if err == nil {
		return command.search(args[1:])
	}

	return c, args, nil
}

func (c *Command) printHelp() {
    fmt.Printf("Usage: %s", c.Name)

    if len(c.Options) > 0 {
        fmt.Printf(" [options]")
    }

    for _, arg := range c.Arguments {
        fmt.Printf(" <%s>", arg)
    }

    if len(c.Options) > 0 {
        fmt.Printf("\n\n")
        fmt.Printf("Options:\n")
    }
    for _, option := range c.Options {
        fmt.Printf("\t")
        if option.Short != 0 {
            fmt.Printf("-%c", option.Short)
        }
        if option.Short != 0 && option.Prompt != "" {
            fmt.Printf(" | ")
        }
        if option.Prompt != "" {
            fmt.Printf("--%s", option.Prompt)
        }

        if !option.IsFlag {
            fmt.Printf(" = <%s>", option.Name)
        }
        fmt.Printf("\n")
    }

    if len(c.Commands) > 0 {
        fmt.Printf("\n")
        fmt.Printf("Subcommands:\n")
    }
    for _, command := range c.Commands {
        fmt.Printf("\t%s\n", command.Name)
    }
    fmt.Printf("\n")
}

func (a *App) search(args []string) (*Command, []string, error) {
	c := Command(*a)
	command := &c

	return command.search(args)
}

func (a *App) printHelp() {
    c := Command(*a)
    command := &c

    command.printHelp()
}

func (a *App) shareContext() {
	context := &Context{}
	context.Args = make(map[string]string)
	context.Flags = make(map[string]bool)

	a.context = context

	for i := range a.Commands {
		a.Commands[i].shareContext(context)
	}
}

func (a *App) Run(args []string) error {
    if a.Name == "" {
        a.Name = "app"
    }

	a.shareContext()

	command, args, err := a.search(args[1:])
	if err != nil {
		return err
	}

	err = command.setArgs(args)
	if err != nil {
		return err
	}

    if _, h := a.context.Get("help"); h && a.context.size() == 1 {
        a.printHelp()
        return nil
    }

	if command.Action == nil {
		return errors.New("missing action for command " + command.Name)
	}

	return command.Action(command.context)
}
