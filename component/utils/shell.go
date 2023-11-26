package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aronlt/toolkit/ds"
	"github.com/aronlt/toolkit/terror"
	"github.com/aronlt/toolkit/tio"
	"github.com/aronlt/toolkit/tutils"
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

func GetFromStdio(hint string, words ...string) string {
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
