package component

import (
	"fmt"

	"github.com/aronlt/toolkit/terror"
	"github.com/aronlt/xman/component/utils"
	"github.com/urfave/cli/v2"
)

type ListGitRemoteBranch struct{}

func NewListGitRemoteBranch() *ListGitRemoteBranch {
	return &ListGitRemoteBranch{}
}

func (l *ListGitRemoteBranch) Name() string {
	return "remote_branch"
}

func (l *ListGitRemoteBranch) Usage() string {
	return "显示远程分支情况"
}

func (l *ListGitRemoteBranch) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (l *ListGitRemoteBranch) Run(_ *cli.Context) error {
	err := utils.PullRemoteRepository()
	if err != nil {
		return terror.Wrap(err, "call PullRemoteRepository")
	}
	cmd := `git for-each-ref --format='%(color:cyan)%(authordate:format:%Y/%m/%d %I:%M %p)    %(align:25,left)%(color:yellow)%(authorname)%(end) %(color:reset)%(refname:strip=2)' --sort=authordate refs/remotes`
	out, err := utils.RunCmdWithOutput(cmd, false)
	if err != nil {
		return terror.Wrapf(err, "call RunCmdWithOutput fail, cmd:%s", cmd)
	}
	fmt.Println()
	fmt.Printf("%s", out)
	return nil
}
