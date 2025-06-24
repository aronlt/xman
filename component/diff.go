package component

import (
	"bytes"
	_ "embed"
	"os"
	"os/exec"

	"github.com/aronlt/toolkit/terror"
	"github.com/urfave/cli/v2"
)

type Diff struct{}

func NewDiff() *Diff {
	return &Diff{}
}

func (d *Diff) Name() string {
	return "diff"
}

func (d *Diff) Usage() string {
	return "查看文件变更情况"
}

func (d *Diff) Flags() []cli.Flag {
	return []cli.Flag{}
}

//go:embed shell/diff.sh
var diffScript string

func (d *Diff) Run(_ *cli.Context) error {
	cmd := exec.Command("bash")
	cmd.Stdin = bytes.NewBufferString(diffScript)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return terror.Wrap(err, "run diff script fail")
	}
	return nil
}
