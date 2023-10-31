package component

import (
	"github.com/aronlt/xman/component/utils"

	"github.com/aronlt/toolkit/terror"
	"github.com/urfave/cli/v2"
)

func Merge(_ *cli.Context) error {
	currentBranch, err := utils.GitCurrentBranch()
	if err != nil {
		return terror.Wrap(err, "call gitCurrentBranch fail")
	}
	targetBranch := utils.GetFromStdio("要合入的目标分支")
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
