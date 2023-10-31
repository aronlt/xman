package component

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aronlt/xman/component/utils"

	"github.com/aronlt/toolkit/ds"
	"github.com/aronlt/toolkit/terror"
	"github.com/aronlt/toolkit/tio"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Tidy(_ *cli.Context) error {
	pushBranch := utils.GetFromStdio("推送到什么分支")
	moduleName := utils.GetFromStdio("要替换的模块名")
	branchName := utils.GetFromStdio("模块替换后的分支")
	if ds.SliceIncludeUnpack("", pushBranch, moduleName, branchName) {
		return fmt.Errorf("invalid arguments:%+v, %+v, %+v", pushBranch, moduleName, branchName)
	}
	err := utils.GitCheckDirtyZone()
	if err != nil {
		return terror.Wrap(err, "work zone is dirty")
	}
	err = utils.GitCheckout(pushBranch)
	if err != nil {
		return err
	}
	err = utils.GitPull(pushBranch)
	if err != nil {
		logrus.WithError(err).Errorf("call GitPull:%s fail", pushBranch)
	}

	err = utils.GitCheckConflict()
	if err != nil {
		return err
	}
	err = replace(moduleName, branchName)
	if err != nil {
		return err
	}

	err = utils.RunCmd("go mod tidy")
	if err != nil {
		return err
	}

	err = utils.GitAddAndCommit()
	if err != nil {
		return err
	}
	utils.WaitForKeyPress()
	err = utils.GitPush(pushBranch)
	if err != nil {
		return err
	}
	logrus.Infof("run script success...")
	return nil
}

func replace(moduleName string, branchName string) error {
	wd, err := os.Getwd()
	if err != nil {
		err = terror.Wrap(err, "call os.Getwd fail")
		return err
	}
	filename := filepath.Join(wd, "go.mod")
	lines, err := tio.ReadLines(filename)
	if err != nil {
		err = terror.Wrap(err, "call tio.ReadLines fail")
		return err
	}
	inRequire := false
	replaceCount := 0
	for i, lineContent := range lines {
		line := string(lineContent)
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "require") {
			inRequire = true
			continue
		}
		if line == ")" {
			inRequire = false
			continue
		}
		if inRequire {
			parts := strings.Split(line, " ")
			if len(parts) < 2 {
				return fmt.Errorf("warning: invalid module, line:%s", line)
			}
			module := parts[0]
			if strings.HasSuffix(module, moduleName) {
				parts[1] = branchName
				newLine := strings.Join(parts, " ")
				lines[i] = []byte(newLine)
				logrus.Infof("replace module from:%s -> %s", line, newLine)
				replaceCount += 1
			}
		}
	}
	if replaceCount == 0 {
		logrus.Infof("not replace any module, skip")
		return nil
	}
	_, err = tio.WriteFile(filename, bytes.Join(lines, []byte("\n")), false)
	if err != nil {
		return terror.Wrap(err, "tio.WriteFile fail")
	}
	err = os.Chmod(filename, 0755)
	if err != nil {
		return terror.Wrap(err, "os.Chmod fail")
	}
	logrus.Infof("replace success...")
	return nil
}
