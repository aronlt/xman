package component

import (
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
	return "分支推送"
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
	err = utils.RunCmd("go mod tidy")
	if err != nil {
		return terror.Wrap(err, "call go mod tidy fail")
	}
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
	err = utils.PushBranch(currentBranch)
	if err != nil {
		return terror.Wrap(err, "call utils.PushBranch fail")
	}
	return nil
}
