package component

import (
	"fmt"

	"github.com/aronlt/toolkit/terror"
	"github.com/aronlt/xman/component/utils"
	"github.com/urfave/cli/v2"
)

type ListGitLocalBranch struct{}

func NewListGitBranch() *ListGitLocalBranch {
	return &ListGitLocalBranch{}
}

func (l *ListGitLocalBranch) Name() string {
	return "local_branch"
}

func (l *ListGitLocalBranch) Usage() string {
	return "显示本地分支情况"
}

func (l *ListGitLocalBranch) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (l *ListGitLocalBranch) Run(_ *cli.Context) error {
	cmd := `git for-each-ref --format='%(color:cyan)%(authordate:format:%Y/%m/%d %I:%M %p)    %(align:25,left)%(color:yellow)%(authorname)%(end) %(color:reset)%(refname:strip=2)' --sort=authordate refs/heads`
	out, err := utils.RunCmdWithOutput(cmd, false)
	if err != nil {
		return terror.Wrapf(err, "call RunCmdWithOutput fail, cmd:%s", cmd)
	}
	fmt.Println()
	fmt.Printf("%s", out)
	return nil
}
