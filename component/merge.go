package component

import (
	"github.com/aronlt/xman/component/utils"

	"github.com/aronlt/toolkit/terror"
	"github.com/urfave/cli/v2"
)

type Merge struct {
}

func NewMerge() *Merge {
	return &Merge{}
}

func (m *Merge) Name() string {
	return "merge"
}

func (m *Merge) Usage() string {
	return "分支合并"
}

func (m *Merge) Run(_ *cli.Context) error {
	currentBranch, err := utils.GitCurrentBranch()
	if err != nil {
		return terror.Wrap(err, "call gitCurrentBranch fail")
	}
	branches, err := utils.ListAllBranch()
	if err != nil {
		terror.Wrap(err, "call ListAllBranch fail")
	}

	targetBranch := utils.GetFromStdio("要合入的目标分支", branches...)
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
