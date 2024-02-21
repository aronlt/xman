package component

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aronlt/toolkit/ds"
	"github.com/aronlt/toolkit/terror"
	"github.com/aronlt/toolkit/tio"
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
	err := utils.GitAddAll()
	if err != nil {
		return terror.Wrap(err, "call GitAddAll fail")
	}
	files, err := utils.GetAddFiles(ds.SetFromSlice([]string{"lint_output.txt", ".golangci.yml"}))
	if err != nil {
		return terror.Wrap(err, "call GetAddFiles fail")
	}
	if len(files) == 0 {
		logrus.Infof("暂存区无提交文件,跳过lint")
		return nil
	}
	dirs := make(map[string][]string)
	for _, file := range files {
		dir := filepath.Dir(file)
		ds.MapOpAppendValue(dirs, dir, file)
	}
	lintDir := ctx.Bool("dir")
	count := 0
	for dir, dirFiles := range dirs {
		scanContent := strings.Builder{}
		if lintDir {
			scanContent.WriteString(dir)
		} else {
			for _, file := range dirFiles {
				scanContent.WriteString(" " + file)
			}
		}
		appendOut := ""
		if count > 0 {
			appendOut = "-a"
		}
		cmd := fmt.Sprintf("golangci-lint run %s --tests=false --timeout=5m | tee %s lint_output.txt", scanContent.String(), appendOut)
		if ok, err := tio.ExistFile("~/.golangci.yml"); ok && err == nil {
			cmd = fmt.Sprintf("golangci-lint run %s --tests=false --timeout=5m -c ~/.golangci.yml | tee %s lint_output.txt", scanContent.String(), appendOut)
		}
		if ok, err := tio.ExistFile(".golangci.yml"); ok && err == nil {
			cmd = fmt.Sprintf("golangci-lint run %s --tests=false --timeout=5m -c .golangci.yml | tee %s lint_output.txt", scanContent.String(), appendOut)
		}
		count++
		logrus.Infof("执行lint命令:%s", cmd)
		err = utils.RunCmd(cmd)
		if err != nil {
			return terror.Wrap(err, "call RunCmd:%s fail", cmd)
		}
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
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    "lint_dir",
			Aliases: []string{"dir"},
			Usage:   "是否以目录维度lint",
		},
	}
}
