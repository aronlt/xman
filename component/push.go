package component

import (
	"github.com/aronlt/xman/component/utils"

	"github.com/aronlt/toolkit/terror"
	"github.com/sirupsen/logrus"
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
			Name:    "skip_review",
			Aliases: []string{"k"},
			Usage:   "跳过add确认",
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
	skipReview := ctx.Bool("skip_review")
	force := "n"
	if !skipReview {
		force = utils.GetFromStdio("是否review git add(输入y确认否则跳过确认)", true)
	}
	if force == "y" {
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
	err = utils.GitPull(currentBranch)
	if err != nil {
		logrus.Warnf("call utils.GitPull fail")
	}
	err = utils.GitCheckConflict()
	if err != nil {
		return terror.Wrap(err, "call GitCheckConflict fail")
	}

	err = utils.GitPush(currentBranch, true)
	if err != nil {
		return terror.Wrap(err, "call utils.GitPush fail")
	}
	return nil
}
