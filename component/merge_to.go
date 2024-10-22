package component

import (
	"github.com/fatih/color"

	"github.com/aronlt/xman/component/utils"

	"github.com/aronlt/toolkit/terror"
	"github.com/urfave/cli/v2"
)

type MergeTo struct {
}

func NewMergeTo() *MergeTo {
	return &MergeTo{}
}

func (m *MergeTo) Name() string {
	return "merge_to"
}

func (m *MergeTo) Usage() string {
	return "并到其他分支"
}

func (m *MergeTo) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "to",
			Aliases: []string{"to"},
			Usage:   "合入的目标分支",
		}}
}

func (m *MergeTo) Run(ctx *cli.Context) error {
	to, err := utils.GitSelectBranch(ctx, "merge_to", "选择要合入的目标分支")
	if err != nil {
		return terror.Wrap(err, "call GitSelectBranch fail")
	}

	color.Green("1.Select to branch:%s success", to)

	current, err := utils.GitTryPullAndCheck()
	if err != nil {
		return terror.Wrapf(err, "call GitTryPullAndCheck fail")
	}

	color.Green("2.Pull current branch:%s success", current)
	err = utils.GitCheckout(to)
	if err != nil {
		return terror.Wrapf(err, "call GitCheckout fail, branch:%s", to)
	}
	_, err = utils.GitTryPullAndCheck()
	if err != nil {
		return terror.Wrapf(err, "call GitTryPullAndCheck fail")
	}

	color.Green("3.Pull to branch:%s success", to)

	err = utils.GitMerge(current)
	if err != nil {
		return terror.Wrapf(err, "call GitMerge fail, branch:%s", current)
	}
	err = utils.GitCheckConflict()
	if err != nil {
		return terror.Wrap(err, "call GitCheckConflict fail")
	}

	color.Green("4.Git merge from:%s to %s success", current, to)
	return nil
}
