package component

import (
	"fmt"

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
			Name:    "commit_msg",
			Aliases: []string{"m"},
			Usage:   "commit提交信息",
		},
		&cli.StringFlag{
			Name:    "merge_to",
			Aliases: []string{"to"},
			Usage:   "合入的目标分支",
		}}
}

func (m *MergeTo) Run(ctx *cli.Context) error {
	currentBranch, err := utils.GitCurrentBranch()
	if err != nil {
		return terror.Wrap(err, "call gitCurrentBranch fail")
	}
	branches, err := utils.ListAllBranch()
	if err != nil {
		return terror.Wrap(err, "call ListAllBranch fail")
	}

	targetBranch := ctx.String("merge_to")
	if targetBranch == "" {
		targetBranch = utils.GetFromStdio("要合入的目标分支", false, branches...)
	}
	if targetBranch == currentBranch {
		return fmt.Errorf("target branch should not equal to current branch")
	}

	commitMsg := ctx.String("commit_msg")
	err = utils.GitAddAndCommit(commitMsg)
	if err != nil {
		return terror.Wrapf(err, "call GitAddAndCommit fail, commit msg:%s", commitMsg)
	}
	err = utils.GitPull(currentBranch)
	if err != nil {
		return terror.Wrapf(err, "call GitPull fail, branch:%s", currentBranch)
	}

	err = utils.GitCheckConflict()
	if err != nil {
		return terror.Wrapf(err, "call GitCheckConflict fail")
	}

	err = utils.GitPush(currentBranch)
	if err != nil {
		return terror.Wrapf(err, "call GitPush fail, branch:%s", currentBranch)
	}

	err = utils.GitCheckout(targetBranch)
	if err != nil {
		return terror.Wrapf(err, "call GitCheckout fail, branch:%s", targetBranch)
	}

	err = utils.GitPull(targetBranch)
	if err != nil {
		return terror.Wrapf(err, "call GitPull fail, branch:%s", targetBranch)
	}
	err = utils.GitCheckConflict()
	if err != nil {
		return terror.Wrapf(err, "call GitCheckConflict fail")
	}

	err = utils.GitMerge(currentBranch)
	if err != nil {
		return terror.Wrapf(err, "call GitMerge fail, branch:%s", currentBranch)
	}
	err = utils.GitPush(targetBranch)
	if err != nil {
		return terror.Wrapf(err, "call GitPush fail, branch:%s", targetBranch)
	}
	return nil
}
