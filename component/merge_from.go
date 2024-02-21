package component

import (
	"fmt"

	"github.com/aronlt/xman/component/utils"

	"github.com/aronlt/toolkit/terror"
	"github.com/urfave/cli/v2"
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
	return "其他分支合入当前分支"
}

func (m *MergeFrom) Flags() []cli.Flag {
	return []cli.Flag{&cli.StringFlag{
		Name:    "commit_msg",
		Aliases: []string{"m"},
		Usage:   "commit提交信息",
	}, &cli.StringFlag{
		Name:    "merge_from",
		Aliases: []string{"from"},
		Usage:   "从那个分支合入",
	}}
}

func (m *MergeFrom) Run(ctx *cli.Context) error {
	currentBranch, err := utils.GitCurrentBranch()
	if err != nil {
		return terror.Wrap(err, "call gitCurrentBranch fail")
	}
	branches, err := utils.ListAllBranch()
	if err != nil {
		return terror.Wrap(err, "call ListAllBranch fail")
	}

	err = utils.GitPull(currentBranch)
	if err != nil {
		return terror.Wrapf(err, "call GitPull fail, branch:%s", currentBranch)
	}

	err = utils.GitCheckConflict()
	if err != nil {
		return terror.Wrap(err, "call GitCheckConflict fail")
	}
	commitMsg := ctx.String("commit_msg")
	err = utils.GitAddAndCommit(commitMsg)
	if err != nil {
		return terror.Wrapf(err, "call GitAddAndCommit fail, commit msg:%s", commitMsg)
	}

	fromBranch := ctx.String("merge_from")
	if fromBranch == "" {
		fromBranch = utils.GetFromStdio("从哪个分支合入", false, branches...)
	}
	if fromBranch == currentBranch {
		return fmt.Errorf("target branch should not equal to current branch")
	}

	err = utils.GitMerge(fromBranch)
	if err != nil {
		return terror.Wrapf(err, "call GitMerge fail, from branch:%s", fromBranch)
	}
	err = utils.GitCheckConflict()
	if err != nil {
		return terror.Wrapf(err, "call GitCheckConflict fail")
	}
	err = utils.GitPush(currentBranch)
	if err != nil {
		return terror.Wrapf(err, "call GitPush fail, branch:%s", currentBranch)
	}
	return nil
}
