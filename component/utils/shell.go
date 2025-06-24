package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aronlt/toolkit/ds"
	"github.com/aronlt/toolkit/terror"
	"github.com/aronlt/toolkit/tio"
	"github.com/aronlt/toolkit/tutils"
	"github.com/fatih/color"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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

func RunCmdWithOutput(cmd string, printLog bool) (string, error) {
	env := map[string]string{}
	result := tutils.RunCmd(cmd, env)
	if result.Error() != nil {
		if printLog {
			logrus.Errorf("run cmd:%s fail, error:%+v", cmd, result.Error())
		}
		return "", result.Error()
	}
	if printLog {
		logrus.Infof("run cmd:%s success, output:%s", cmd, result.String())
	}
	return result.String(), nil
}

func RunCmd(cmd string) error {
	_, err := RunCmdWithOutput(cmd, true)
	return err
}

func RunShellScript(scriptPath string) error {
	cmd := exec.Command("bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		return terror.Wrapf(err, "failed to run script %s", scriptPath)
	}
	return nil
}

func SimpleGetFromStdio(hint string) string {
	color.Red("%s:", hint)
	reader := bufio.NewReader(os.Stdin)
	content, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(content)
}

func GetFromStdio(hint string, simple bool, words ...string) string {
	if simple {
		return SimpleGetFromStdio(hint)
	}
	app := tview.NewApplication()
	inputField := tview.NewInputField().
		SetLabel(hint + ":").
		SetFieldWidth(30).
		SetDoneFunc(func(key tcell.Key) {
			app.Stop()
		})
	inputField.SetAutocompleteFunc(func(currentText string) (entries []string) {
		if len(currentText) == 0 {
			return
		}
		for _, word := range words {
			if strings.Contains(strings.ToLower(word), strings.ToLower(currentText)) {
				entries = append(entries, word)
			}
		}
		if len(entries) <= 1 {
			entries = nil
		}
		return
	})
	inputField.SetAutocompletedFunc(func(text string, index, source int) bool {
		if source != tview.AutocompletedNavigate {
			inputField.SetText(text)
		}
		return source == tview.AutocompletedEnter || source == tview.AutocompletedClick
	})
	if err := app.EnableMouse(true).SetRoot(inputField, true).Run(); err != nil {
		panic(err)
	}

	return inputField.GetText()
}
