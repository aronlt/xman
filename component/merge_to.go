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
	return []cli.Flag{}
}

func (m *MergeTo) Run(_ *cli.Context) error {
	currentBranch, err := utils.GitCurrentBranch()
	if err != nil {
		return terror.Wrap(err, "call gitCurrentBranch fail")
	}
	branches, err := utils.ListAllBranch()
	if err != nil {
		return terror.Wrap(err, "call ListAllBranch fail")
	}

	targetBranch := utils.GetFromStdio("要合入的目标分支", false, branches...)
	if targetBranch == currentBranch {
		return fmt.Errorf("target branch should not equal to current branch")
	}
	err = utils.GitAddAndCommit()
	if err != nil {
		return err
	}
	err = utils.GitPull(currentBranch)
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

	err = utils.GitCheckout(targetBranch)
	if err != nil {
		return err
	}

	err = utils.GitPull(targetBranch)
	if err != nil {
		return err
	}
	err = utils.GitCheckConflict()
	if err != nil {
		return err
	}

	err = utils.GitMerge(currentBranch)
	if err != nil {
		return err
	}
	err = utils.GitPush(targetBranch)
	if err != nil {
		return err
	}
	return nil
}
