package component

import (
	"fmt"

	"github.com/aronlt/xman/component/utils"
	"github.com/urfave/cli/v2"
)

type ListTags struct{}

func NewListTags() *ListTags {
	return &ListTags{}
}

func (l *ListTags) Name() string {
	return "list_tags"
}

func (l *ListTags) Usage() string {
	return "显示标签情况"
}

func (l *ListTags) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (l *ListTags) Run(_ *cli.Context) error {
	err := utils.PullRemoteRepository()
	if err != nil {
		return err
	}
	cmd := `git for-each-ref --format='%(color:cyan)%(authordate:format:%Y/%m/%d %I:%M %p)    %(align:25,left)%(color:yellow)%(authorname)%(end) %(color:reset)%(refname:strip=1)' --sort=authordate refs/tags`
	out, err := utils.RunCmdWithOutput(cmd, false)
	if err != nil {
		return err
	}
	fmt.Println()
	fmt.Printf("%s", out)
	return nil
}
