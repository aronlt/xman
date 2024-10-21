package component

import (
	"github.com/aronlt/toolkit/terror"
	"github.com/urfave/cli/v2"

	"github.com/aronlt/xman/component/utils"
)

type Recover struct {
}

func NewRecover() *Recover {
	return &Recover{}
}

func (r *Recover) Name() string {
	return "recover"
}

func (r *Recover) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (r *Recover) Usage() string {
	return "从git stash中恢复代码"
}

func (r *Recover) Run(_ *cli.Context) error {
	err := utils.GitCheckDirtyZone()
	if err != nil {
		return terror.Wrap(err, "call GitCheckDirtyZone fail")
	}
	err = utils.GitStashPop()
	if err != nil {
		return terror.Wrap(err, "call GitStashPop fail")
	}
	return nil
}
