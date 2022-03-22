package cli

import "errors"

type Context map[string]string

type Args []string

type Command struct {
    Name string
    Commands Commands
    Arguments Args
    Context *Context
    Action func(c *Context) error
}

type Commands []*Command

type App Command

func (c *Context) Set(key, value string) {
    (*c)[key] = value
}

func (c *Context) Get(key string) string {
    return (*c)[key]
}

func (c *Command) ShareContext(context *Context) {
    c.Context = context
    for i := range c.Commands {
        c.Commands[i].ShareContext(context)
    }
}

func (c Command) SetArgs(args []string) {
    for i := range args {
        c.Context.Set(c.Arguments[i], args[i])
    }
}

func (c Commands) Get(name string) (*Command, error) {
    for i := range c {
        if c[i].Name == name {
            return c[i], nil
        }
    }
    return nil, errors.New("command not found: " + name)
}

func (c Commands) Search(args []string) (*Command, error) {
    command, err := c.Get(args[0])
    if err != nil {
        return nil, err
    }

    if len(command.Commands) == 0 {
        command.SetArgs(args[1:])
        return command, nil
    }

    return command.Commands.Search(args[1:])
}

func (a *App) shareContext() {
    context := &Context{}
    a.Context = context

    for i := range a.Commands {
        a.Commands[i].ShareContext(context)
    }
}

func (a *App) Run(args []string) error {
    a.shareContext()

    command, err := a.Commands.Search(args[1:])
    if err != nil {
        return err
    }

    if command.Action == nil {
        return errors.New("missing action for command " + command.Name)
    }

    command.Action(command.Context)

    return nil
}
