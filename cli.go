package cli

import (
	"errors"
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
	Context   *Context
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

func (c *Command) ShareContext(context *Context) {
	c.Context = context
	for i := range c.Commands {
		c.Commands[i].ShareContext(context)
	}
}

func (c *Command) SetArg(key, value string) {
	c.Context.Args[key] = value
}

func (c *Command) FindOptionPrompt(name string) *Option {
	for i := range c.Options {
		if c.Options[i].Prompt == name {
			return c.Options[i]
		}
	}
	return nil
}

func (c *Command) FindOptionShort(name rune) *Option {
	for i := range c.Options {
		if c.Options[i].Short == name {
			return c.Options[i]
		}
	}
	return nil
}

func (c *Command) GetOption(name string) *Option {
	if name[1] == '-' {
		return c.FindOptionPrompt(name[2:])
	}
	return c.FindOptionShort(rune(name[1]))
}

func (c *Command) SetOption(name, value string) {
	c.Context.Args[name] = value
}

func (c *Command) SetFlag(name string) {
	c.Context.Flags[name] = true
}

func (c *Command) SetArgs(args []string) error {
	a := 0

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg[0] == '-' {
			option := c.GetOption(arg)

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

				c.SetOption(name, value)
			} else {
				c.SetFlag(name)
			}

		} else {
			c.SetArg(c.Arguments[a], arg)
			a++
		}
	}

	if len(c.Arguments) != a {
		return errors.New("mismatcing arguments")
	}

	return nil
}

func (c Command) Get(name string) (*Command, error) {
	for i := range c.Commands {
		if c.Commands[i].Name == name {
			return c.Commands[i], nil
		}
	}
	return nil, errors.New("command not found: " + name)
}

func (c *Command) Search(args []string) (*Command, []string, error) {
    if len(args) == 0 {
        return c, nil, nil
    }

	command, err := c.Get(args[0])
	if err == nil {
		return command.Search(args[1:])
	}

	return c, args, nil
}

func (a *App) shareContext() {
	context := &Context{}
	context.Args = make(map[string]string)
	context.Flags = make(map[string]bool)

	a.Context = context

	for i := range a.Commands {
		a.Commands[i].ShareContext(context)
	}
}

func (a *App) Search(args []string) (*Command, []string, error) {
	c := Command(*a)
	command := &c

	return command.Search(args)
}

func (a *App) Run(args []string) error {
	a.shareContext()

	command, args, err := a.Search(args[1:])
	if err != nil {
		return err
	}

	err = command.SetArgs(args)
	if err != nil {
		return err
	}

	if command.Action == nil {
		return errors.New("missing action for command " + command.Name)
	}

	return command.Action(command.Context)
}
