package component

import (
	"github.com/aronlt/toolkit/terror"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	"github.com/aronlt/xman/component/utils"
)

type CheckoutBranch struct{}

func NewCheckoutBranch() *CheckoutBranch {
	return &CheckoutBranch{}
}

func (cb *CheckoutBranch) Name() string {
	return "checkout"
}

func (cb *CheckoutBranch) Usage() string {
	return "切换到其他分支 --t 切换到的分支名 --f 从哪个分支切换"
}

func (cb *CheckoutBranch) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "to",
			Aliases: []string{"t"},
			Usage:   "切换到的分支名",
		},
		&cli.StringFlag{
			Name:    "from",
			Aliases: []string{"f"},
			Usage:   "从特定分支切换出新的分支",
		},
	}
}

func (cb *CheckoutBranch) Run(ctx *cli.Context) error {
	to := ctx.String("to")
	if to == "" {
		branches, err := utils.ListAllBranch()
		if err != nil {
			return terror.Wrap(err, "call ListAllBranch fail")
		}
		to = utils.GetFromStdio("切换到的分支", false, branches...)
	}
	color.Green("1.Select to branch:%s success", to)

	err := utils.GitAddAll()
	if err != nil {
		return terror.Wrap(err, "call GitAddAndCommit fail")
	}

	color.Green("2.Git add current branch success")
	from := ctx.String("from")
	if from == "" {
		err = utils.GitCheckout(to)
		if err != nil {
			return terror.Wrapf(err, "call GitCheckout fail, branch:%s", to)
		}
		color.Green("3.Checkout to branch:%s success", to)
		return nil
	}
	err = utils.GitCheckout(from)
	if err != nil {
		return terror.Wrap(err, "call GitCheckout fail, branch:%s", from)
	}
	_, err = utils.GitTryPullAndCheck()
	if err != nil {
		return terror.Wrap(err, "call GitTryPullAndCheck fail")
	}
	err = utils.GitCheckout(to, true)
	if err != nil {
		return terror.Wrap(err, "call GitCheckout fail, branch:%s", to)
	}
	color.Green("3.Checkout from branch %s to branch:%s success", from, to)
	return nil
}
