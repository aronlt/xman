package component

import (
	"fmt"
	"time"

	"github.com/aronlt/xman/component/utils"

	"github.com/aronlt/toolkit/terror"
	"github.com/aronlt/toolkit/tutils"
	"github.com/urfave/cli/v2"
)

func Stash(_ *cli.Context) error {
	currentBranch, err := utils.GitCurrentBranch()
	if err != nil {
		return terror.Wrap(err, "call gitCurrentBranch fail")
	}
	t := tutils.TimeToString(time.Now())
	msg := utils.GetFromStdio("描述信息")
	err = utils.GitStash(fmt.Sprintf("branch:%s;time:%s;msg:%s", currentBranch, t, msg))
	if err != nil {
		return terror.Wrap(err, "call GitStash fail")
	}
	return nil
}

func Recover(_ *cli.Context) error {
	err := utils.GitCheckDirtyZone()
	if err != nil {
		return terror.Wrap(err, "call GitCheckDirtyZone fail")
	}
	branch := utils.GetFromStdio("切换的分支")
	err = utils.GitCheckout(branch)
	if err != nil {
		return terror.Wrap(err, "call GitCheckDirtyZone fail")
	}
	err = utils.GitStashPop()
	if err != nil {
		return terror.Wrap(err, "call GitStashPop fail")
	}
	return nil
}
