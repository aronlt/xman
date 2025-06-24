package component

import (
	"github.com/aronlt/toolkit/terror"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	"github.com/aronlt/xman/component/utils"
)

type MergeFrom struct {
}

func NewMergeFrom() *MergeFrom {
	return &MergeFrom{}
}

func (m *MergeFrom) Name() string {
	return "merge_from"
}

func (m *MergeFrom) Usage() string {
	return "其他分支合入当前分支 --f 从哪个分支合入"
}

func (m *MergeFrom) Flags() []cli.Flag {
	return []cli.Flag{&cli.StringFlag{
		Name:    "from",
		Aliases: []string{"f"},
		Usage:   "从那个分支合入",
	}}
}

func (m *MergeFrom) Run(ctx *cli.Context) error {
	current, err := utils.GitTryPullAndCheck()
	if err != nil {
		return terror.Wrapf(err, "call GitTryPullAndCheck fail")
	}
	color.Green("1.Pull current branch:%s success", current)
	from, err := utils.GitSelectBranch(ctx, "from", "选择从哪个分支合入")
	if err != nil {
		return terror.Wrapf(err, "call GitSelectBranch fail")
	}
	color.Green("2.Select merge from branch:%s success", from)
	err = utils.GitCheckout(from)
	if err != nil {
		return terror.Wrapf(err, "call git checkout:%s fail", from)
	}
	_, err = utils.GitTryPullAndCheck()
	if err != nil {
		return terror.Wrapf(err, "call GitTryPullAndCheck fail")
	}
	err = utils.GitCheckout(current)
	if err != nil {
		return terror.Wrapf(err, "call git checkout:%s fail", current)
	}
	color.Green("3.Update git from branch:%s success", from)

	err = utils.GitMerge(from)
	if err != nil {
		return terror.Wrapf(err, "call GitMerge fail, from branch:%s", from)
	}
	err = utils.GitCheckConflict()
	if err != nil {
		return terror.Wrap(err, "call GitCheckConflict fail")
	}

	color.Green("3.Merge git from branch %s to branch %s success", from, current)
	return nil
}
