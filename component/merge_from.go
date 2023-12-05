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
	return []cli.Flag{}
}

func (m *MergeFrom) Run(_ *cli.Context) error {
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
		return err
	}

	err = utils.GitCheckConflict()
	if err != nil {
		return err
	}
	err = utils.GitAddAndCommit()
	if err != nil {
		return err
	}
	fromBranch := utils.GetFromStdio("从哪个分支合入", false, branches...)
	if fromBranch == currentBranch {
		return fmt.Errorf("target branch should not equal to current branch")
	}

	err = utils.GitMerge(fromBranch)
	if err != nil {
		return err
	}
	err = utils.GitCheckConflict()
	if err != nil {
		return err
	}
	err = utils.GitPush(currentBranch)
	if err != nil {
		return err
	}
	return nil
}
