package main

import (
	"os"

	"github.com/aronlt/xman/component"

	"github.com/urfave/cli/v2"
)

type Action interface {
	Run(ctx *cli.Context) error
	Name() string
	Usage() string
	Flags() []cli.Flag
}

func Register() []*cli.Command {
	actions := []Action{
		component.NewMergeTo(),
		component.NewMergeFrom(),
		component.NewStash(),
		component.NewRecover(),
		component.NewPush(),
		component.NewTag(),
		component.NewLint(),
		component.NewListGitBranch(),
		component.NewListGitRemoteBranch(),
		component.NewListTags(),
		component.NewListLastCommit(),
		component.NewCheckoutBranch(),
	}
	commands := make([]*cli.Command, 0, len(actions))
	for i := range actions {
		action := actions[i]
		commands = append(commands, &cli.Command{
			Name:  action.Name(),
			Usage: action.Usage(),
			Flags: action.Flags(),
			Action: func(ctx *cli.Context) error {
				return action.Run(ctx)
			},
		})
	}
	return commands
}

func main() {
	app := &cli.App{
		Commands: Register(),
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
