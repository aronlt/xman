package component

import (
	"github.com/fatih/color"

	"github.com/aronlt/xman/component/utils"

	"github.com/aronlt/toolkit/terror"
	"github.com/urfave/cli/v2"
)

type Push struct {
}

func NewPush() *Push {
	return &Push{}
}

func (p *Push) Name() string {
	return "push"
}

func (p *Push) Usage() string {
	return "编译后再分支推送"
}

func (p *Push) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "commit_msg",
			Aliases: []string{"m"},
			Usage:   "commit提交信息",
		},
		&cli.StringFlag{
			Name:    "review",
			Aliases: []string{"r"},
			Usage:   "add确认",
		}}
}

func (p *Push) Run(ctx *cli.Context) error {
	currentBranch, err := utils.GitCurrentBranch()
	if err != nil {
		return terror.Wrap(err, "call utils.GitCurrentBranch fail")
	}
	content, err := utils.RunCmdWithOutput("go mod tidy", true)
	if err != nil {
		color.Red("call go mod tidy fail, content:%s", content)
		return terror.Wrap(err, "call go mod tidy fail")
	}
	color.Green("1.Go mod tidy success")
	content, err = utils.RunCmdWithOutput("go build .", true)
	if err != nil {
		color.Red("call go build fail, content:%s", content)
		return terror.Wrap(err, "call go build fail")
	}
	color.Green("2.Go build success")
	commitMsg := ctx.String("commit_msg")
	review := ctx.Bool("review")
	if review {
		err = utils.GitAddWithConfirm(commitMsg)
		if err != nil {
			return terror.Wrap(err, "call utils.GitAddWithConfirm fail")
		}
	} else {
		err = utils.GitAddAndCommit(commitMsg)
		if err != nil {
			return terror.Wrap(err, "call utils.GitAddAndCommit fail")
		}
	}
	color.Green("3.Add and commit success, commit msg:%s", commitMsg)
	err = utils.PushBranch(currentBranch)
	if err != nil {
		return terror.Wrap(err, "call utils.PushBranch fail")
	}
	color.Green("4.Push current branch:%s success", currentBranch)
	return nil
}
