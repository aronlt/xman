package component

import (
	"fmt"

	"github.com/aronlt/toolkit/terror"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	"github.com/aronlt/xman/component/utils"
)

type ListLastCommit struct{}

func NewListLastCommit() *ListLastCommit {
	return &ListLastCommit{}
}

func (l *ListLastCommit) Name() string {
	return "last_commit"
}

func (l *ListLastCommit) Usage() string {
	return "显示本地最近提交信息"
}

func (l *ListLastCommit) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (l *ListLastCommit) Run(_ *cli.Context) error {
	cmd := `for branch in $(git branch -r | grep -v HEAD);do echo -e $(git show --format="%ci %cr" $branch | head -n 1) \\t$branch; done | sort -r`
	out, err := utils.RunCmdWithOutput(cmd, false)
	if err != nil {
		return terror.Wrapf(err, "call RunCmdWithOutput fail, cmd:%s", cmd)
	}
	fmt.Println()
	color.Green("%s", out)
	return nil
}
