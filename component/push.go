package component

import (
	"github.com/aronlt/xman/component/utils"

	"github.com/aronlt/toolkit/terror"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Push(_ *cli.Context) error {
	currentBranch, err := utils.GitCurrentBranch()
	if err != nil {
		return terror.Wrap(err, "call utils.GitCurrentBranch fail")
	}
	force := utils.GetFromStdio("是否review git add(输入y确认否则跳过确认)")
	if force == "y" {
		err = utils.GitAddWithConfirm()
		if err != nil {
			return terror.Wrap(err, "call utils.GitAddWithConfirm fail")
		}
	} else {
		err = utils.GitAddAndCommit()
		if err != nil {
			return terror.Wrap(err, "call utils.GitAddAndCommit fail")
		}
	}
	err = utils.GitPull(currentBranch)
	if err != nil {
		logrus.Warnf("call utils.GitPull fail")
	}
	err = utils.GitCheckConflict()
	if err != nil {
		return terror.Wrap(err, "call GitCheckConflict fail")
	}
	err = utils.GitPush(currentBranch)
	if err != nil {
		return terror.Wrap(err, "call utils.GitPush fail")
	}
	return nil
}