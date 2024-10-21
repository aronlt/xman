package component

import (
	"github.com/aronlt/toolkit/terror"
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
	return "切换到其他分支"
}

func (cb *CheckoutBranch) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "branch_name",
			Aliases: []string{"b"},
			Usage:   "切换到的分支名",
		},
		&cli.StringFlag{
			Name:    "checkout_from",
			Aliases: []string{"cf"},
			Usage:   "从特定分支切换出新的分支",
		},
	}
}

func (cb *CheckoutBranch) Run(ctx *cli.Context) error {
	targetBranch := ctx.String("branch_name")
	if targetBranch == "" {
		branches, err := utils.ListAllBranch()
		if err != nil {
			return terror.Wrap(err, "call ListAllBranch fail")
		}
		targetBranch = utils.GetFromStdio("切换到的分支", false, branches...)
	}
	currentBranch, err := utils.GitCurrentBranch()
	if err != nil {
		return terror.Wrap(err, "call utils.GitCurrentBranch fail")
	}
	err = utils.GitAddAndCommit("临时保存")
	if err != nil {
		return terror.Wrap(err, "call GitAddAndCommit fail")
	}
	err = utils.PushBranch(currentBranch)
	if err != nil {
		return terror.Wrapf(err, "call PushBranch fail, branch:%s", currentBranch)
	}
	checkoutFrom := ctx.String("checkout_from")
	if checkoutFrom == "" {
		err = utils.GitCheckout(targetBranch)
		if err != nil {
			return terror.Wrapf(err, "call GitCheckout fail, branch:%s", targetBranch)
		}
		return nil
	}
	err = utils.GitCheckout(checkoutFrom)
	if err != nil {
		return terror.Wrap(err, "call GitCheckout fail, branch:%s", checkoutFrom)
	}
	err = utils.PullBranch(checkoutFrom)
	if err != nil {
		return terror.Wrap(err, "call PullBranch fail, branch:%s", checkoutFrom)
	}
	err = utils.GitCheckout(targetBranch, true)
	if err != nil {
		return terror.Wrap(err, "call GitCheckout fail, branch:%s", targetBranch)
	}
	return nil
}
