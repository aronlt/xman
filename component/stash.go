package component

import (
	"fmt"
	"time"

	"github.com/aronlt/xman/component/utils"

	"github.com/aronlt/toolkit/terror"
	"github.com/aronlt/toolkit/tutils"
	"github.com/urfave/cli/v2"
)

type Stash struct {
}

func NewStash() *Stash {
	return &Stash{}
}

func (s *Stash) Name() string {
	return "stash"
}

func (s *Stash) Usage() string {
	return "暂存当前代码到stash中"
}

func (s *Stash) Run(_ *cli.Context) error {
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
