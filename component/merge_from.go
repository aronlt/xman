package component

import (
	"fmt"

	"github.com/aronlt/toolkit/ds"
	"github.com/sirupsen/logrus"

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
	if err = utils.GitCheckDirtyZone(); err != nil {
		return terror.Wrap(err, "can not merge when work dir is dirty")
	}

	err = utils.GitPull(currentBranch)
	if err != nil {
		logrus.Errorf("call GitPull fail, branch:%s, error:%v", currentBranch, err)
	} else {
		err = utils.GitCheckConflict()
		if err != nil {
			return terror.Wrap(err, "call GitCheckConflict fail")
		}
	}
	logrus.Infof("1. pull current branch:%s success", currentBranch)

	branches, err := utils.ListAllBranch()
	if err != nil {
		return terror.Wrap(err, "call ListAllBranch fail")
	}
	fromBranch := ctx.String("merge_from")
	if fromBranch == "" {
		fromBranch = utils.GetFromStdio("从哪个分支合入", false, branches...)
	}
	if fromBranch == currentBranch {
		return fmt.Errorf("target branch should not equal to current branch:%s", fromBranch)
	}
	if !ds.SliceInclude(branches, fromBranch) {
		return fmt.Errorf("from branch:%s not in all branch", fromBranch)
	}

	logrus.Infof("2. check merge from branch:%s success", fromBranch)
	err = utils.GitCheckout(fromBranch)
	if err != nil {
		return fmt.Errorf("call git checkout:%s fail", fromBranch)
	}
	err = utils.GitPull(fromBranch)
	if err != nil {
		logrus.Errorf("call GitPull fail, branch:%s, error:%v", currentBranch, err)
	} else {
		err = utils.GitCheckConflict()
		if err != nil {
			return terror.Wrap(err, "call GitCheckConflict fail")
		}
	}
	err = utils.GitCheckout(currentBranch)
	if err != nil {
		return fmt.Errorf("call git checkout:%s fail", fromBranch)
	}
	logrus.Infof("3. call git update from branch:%s success", fromBranch)

	err = utils.GitMerge(fromBranch)
	if err != nil {
		return terror.Wrapf(err, "call GitMerge fail, from branch:%s", fromBranch)
	}
	err = utils.GitCheckConflict()
	if err != nil {
		return terror.Wrapf(err, "call GitCheckConflict fail")
	}

	logrus.Infof("3. merge from %s to %s success", fromBranch, currentBranch)
	return nil
}
