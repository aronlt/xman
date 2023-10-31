package main

import (
	"os"

	"github.com/aronlt/xman/component"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "mod",
				Usage: "更新依赖的模块信息",
				Flags: []cli.Flag{},
				Action: func(ctx *cli.Context) error {
					err := component.Tidy(ctx)
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "recover",
				Usage: "恢复当前模块",
				Action: func(ctx *cli.Context) error {
					err := component.Recover(ctx)
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "stash",
				Usage: "暂存当前模块",
				Action: func(ctx *cli.Context) error {
					err := component.Stash(ctx)
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "merge",
				Usage: "分支合并",
				Flags: []cli.Flag{},
				Action: func(ctx *cli.Context) error {
					err := component.Merge(ctx)
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "push",
				Usage: "分支推送",
				Flags: []cli.Flag{},
				Action: func(ctx *cli.Context) error {
					err := component.Push(ctx)
					if err != nil {
						return err
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
