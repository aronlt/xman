package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aronlt/toolkit/ds"
	"github.com/aronlt/toolkit/terror"
	"github.com/aronlt/toolkit/tio"
	"github.com/aronlt/toolkit/tutils"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

func SearchFile(cmd string) (string, error) {
	dirs := []string{"/usr/bin/", "/usr/local/bin/"}

	var ok bool
	var err error
	var file string
	for _, dir := range dirs {
		file = filepath.Join(dir, cmd)
		ok, err = tio.ExistFile(file)
		if err != nil {
			return "", terror.Wrap(err, "call tio.ExistFile fail")
		}
		if !ok {
			logrus.Warnf("call tio.ExistFile fail, file:%s not exist", file)
			continue
		}
		return file, nil
	}

	return "", fmt.Errorf("can't find cmd:%s in all dirs", cmd)
}

func RunInteractiveCmd(cmd string, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return terror.Wrap(err, "call os.Getwd fail")
	}
	pa := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Dir:   cwd,
	}

	logrus.Infof("starting a new interactive shell")
	ds.SliceOpInsert(&args, 0, cmd)
	file, err := SearchFile(cmd)
	if err != nil {
		return terror.Wrap(err, "call SearchFile fail")
	}

	proc, err := os.StartProcess(file, args, &pa)
	if err != nil {
		return terror.Wrap(err, "call StartProcess fail")
	}

	state, err := proc.Wait()
	if err != nil {
		return terror.Wrap(err, "call proc.Wait fail")
	}

	logrus.Infof("exited interactive shell: %s", state.String())
	return nil
}

func RunCmdWithOutput(cmd string) (string, error) {
	env := map[string]string{}
	result := tutils.RunCmd(cmd, env)
	if result.Error() != nil {
		logrus.Infof("run cmd:%s fail, error:%+v", cmd, result.Error())
		return "", result.Error()
	}
	logrus.Infof("run cmd:%s success, output:%s", cmd, result.String())
	return result.String(), nil
}

func RunCmd(cmd string) error {
	_, err := RunCmdWithOutput(cmd)
	return err
}

func WaitForKeyPress() {
	color.Red("输入任意键继续, 输入Ctl+C退出...")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadByte()
}

func GetFromStdio(hint string) string {
	color.Red("请输入%s:", hint)
	reader := bufio.NewReader(os.Stdin)
	content, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(content)
}
