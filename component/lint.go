package component

import (
	"fmt"
	"strings"

	"github.com/aronlt/toolkit/terror"
	"github.com/aronlt/xman/component/utils"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type Lint struct {
}

func NewLint() *Lint {
	return &Lint{}
}

func (l *Lint) Run(ctx *cli.Context) error {
	dirs, err := utils.GetAddFiles()
	if err != nil {
		return terror.Wrap(err, "call GetAddFiles fail")
	}
	if len(dirs) == 0 {
		logrus.Infof("暂存区无提交文件,跳过lint")
		return nil
	}
	dirContent := strings.Builder{}
	for _, dir := range dirs {
		dirContent.WriteString(" " + dir)
	}
	cmd := fmt.Sprintf("golangci-lint run %s --tests=false --timeout=5m -c .golangci.yml | tee lint_output.txt", dirContent.String())
	logrus.Infof("执行lint命令:%s", cmd)
	err = utils.RunCmd(cmd)
	if err != nil {
		return terror.Wrap(err, "call RunCmd:%s fail", cmd)
	}

	return nil
}

func (l *Lint) Name() string {
	return "lint"
}

func (l *Lint) Usage() string {
	return "lint代码"
}

func (l *Lint) Flags() []cli.Flag {
	return []cli.Flag{}
}
