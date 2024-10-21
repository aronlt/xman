package component

import (
	"fmt"

	"github.com/aronlt/toolkit/ds"
	"github.com/sirupsen/logrus"

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
	if err = utils.GitCheckDirtyZone(); err != nil {
		return terror.Wrap(err, "can not merge when work dir is dirty")
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
		return fmt.Errorf("target branch should not equal to current branch:%s", targetBranch)
	}
	if !ds.SliceInclude(branches, targetBranch) {
		return fmt.Errorf("targe branch:%s not in branch list", targetBranch)
	}

	logrus.Infof("1. check target branch:%s success", targetBranch)

	err = utils.GitPull(currentBranch)
	if err != nil {
		logrus.Errorf("call GitPull fail, branch:%s, error:%v", currentBranch, err)
	} else {
		err = utils.GitCheckConflict()
		if err != nil {
			return terror.Wrap(err, "call GitCheckConflict fail")
		}
	}

	logrus.Infof("2. pull current branch:%s success", currentBranch)
	err = utils.GitCheckout(targetBranch)
	if err != nil {
		return terror.Wrapf(err, "call GitCheckout fail, branch:%s", targetBranch)
	}

	err = utils.GitPull(targetBranch)
	if err != nil {
		logrus.Errorf("call GitPull fail, branch:%s, error:%v", targetBranch, err)
	} else {
		err = utils.GitCheckConflict()
		if err != nil {
			return terror.Wrapf(err, "call GitCheckConflict fail")
		}
	}
	logrus.Infof("3. pull target branch:%s success", targetBranch)

	err = utils.GitMerge(currentBranch)
	if err != nil {
		return terror.Wrapf(err, "call GitMerge fail, branch:%s", currentBranch)
	}
	err = utils.GitCheckConflict()
	if err != nil {
		return terror.Wrap(err, "call GitCheckConflict fail")
	}
	logrus.Infof("4. git merge from:%s to %s success", currentBranch, targetBranch)
	return nil
}
