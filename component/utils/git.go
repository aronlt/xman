package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aronlt/toolkit/ds"
	"github.com/aronlt/toolkit/terror"
	"github.com/fatih/color"
)

func GitCheckConflict() error {
	err := RunCmd("[ $(git ls-files -u  | cut -f 2 | sort -u | wc -l) -eq 0 ] && exit 0 || exit 1")
	if err != nil {
		return err
	}
	return nil
}

func GitAddAndCommit(msg string) error {
	err := RunCmd("git add .")
	if err != nil {
		return err
	}
	if msg == "" {
		err = RunInteractiveCmd("git", []string{"commit"})
	} else {
		err = RunInteractiveCmd("git", []string{"commit", "-m", msg})
	}
	return err
}

func GitCheckDirtyZone() error {
	err := RunCmd("[ $(git status -s | wc -l) -eq 0 ] && exit 0 || exit 1")
	if err != nil {
		return terror.Wrap(err, "work zone is dirty")
	}
	return nil
}

func GitAddTag(tag string) error {
	if tag == "" {
		return fmt.Errorf("empty tag")
	}
	err := RunCmd("git tag " + tag)
	if err != nil {
		return terror.Wrap(err, "call push tag fail")
	}
	return nil
}

func GitPushTag(tag string) error {
	if tag == "" {
		return fmt.Errorf("empty tag")
	}
	err := RunCmd("git push origin " + tag)
	if err != nil {
		return terror.Wrap(err, "call push tag fail")
	}
	return nil
}

func GitPull(branch string) error {
	err := RunCmd("git pull origin " + branch)
	if err != nil {
		return terror.Wrap(err, "run cmd fail")
	}
	return nil
}

func GitPush(branch string, force ...bool) error {
	if len(force) == 0 {
		color.Red("输入y/Y开始推送到远程...")
		reader := bufio.NewReader(os.Stdin)
		y, _ := reader.ReadByte()
		if y != 'y' && y != 'Y' {
			return nil
		}
	}
	err := RunCmd("git push origin " + branch)
	if err != nil {
		return terror.Wrap(err, "run cmd fail")
	}
	return nil
}

func GitPullTags() error {
	err := RunCmd("git pull --tags")
	if err != nil {
		return terror.Wrap(err, "run cmd fail")
	}
	return nil
}

func GitTags() ([]string, error) {
	content, err := RunCmdWithOutput("git tag")
	if err != nil {
		return nil, terror.Wrap(err, "call RunCmdWithOutput fail")
	}
	lines := strings.Split(content, "\n")
	return lines, nil
}

func GitCheckout(branch string) error {
	err := RunCmd("git checkout " + branch)
	if err != nil {
		return terror.Wrap(err, "run cmd fail")
	}
	return nil
}

func GitMerge(branch string) error {
	err := RunCmd("git merge " + branch)
	if err != nil {
		return terror.Wrap(err, "run cmd fail")
	}
	return nil
}

func GitCurrentBranch() (string, error) {
	out, err := RunCmdWithOutput("git rev-parse --abbrev-ref HEAD")
	if err != nil {
		return "", terror.Wrap(err, "run cmd fail")
	}
	return out, nil
}

func GitStash(msg string) error {
	err := RunCmd(fmt.Sprintf("git stash save \"%s\"", msg))
	if err != nil {
		return terror.Wrap(err, "run cmd fail")
	}
	return nil
}

func GitStashPop() error {

	err := RunInteractiveCmd("git", []string{"stash", "list"})
	if err != nil {
		return terror.Wrap(err, "call RunInteractiveCmd fail")
	}

	content, err := RunCmdWithOutput("git stash list")
	lines := strings.Split(content, "\n")
	line := GetFromStdio("恢复的行号(下标0开始)", true)
	lineNum, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		return terror.Wrap(err, "call strconv.Atoi fail")
	}

	stashContent := lines[lineNum]
	start := strings.Index(stashContent, "branch:")
	end := strings.Index(stashContent, ";")
	branch := stashContent[start+len("branch:") : end]
	err = GitCheckout(branch)
	if err != nil {
		return terror.Wrap(err, "call GitCheckout fail")
	}

	err = RunCmd(fmt.Sprintf("git stash pop stash@{%d}", lineNum))
	if err != nil {
		return terror.Wrap(err, "call RunCmd fail")
	}
	return nil
}

func GitAddWithConfirm(msg string) error {
	err := RunInteractiveCmd("git", []string{"add", "-p", "."})
	if err != nil {
		return terror.Wrap(err, "call RunInteractiveCmd fail")
	}
	if msg != "" {
		err = RunInteractiveCmd("git", []string{"commit"})
	} else {
		err = RunInteractiveCmd("git", []string{"commit", "-m", msg})
	}
	if err != nil {
		return terror.Wrap(err, "call RunInteractiveCmd fail")
	}
	return err
}

func ListAllBranch() ([]string, error) {
	content, err := RunCmdWithOutput("git branch")
	if err != nil {
		return nil, terror.Wrap(err, "call RunInteractiveCmd fail")
	}
	branches := make([]string, 0)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.Trim(line, "*")
		line = strings.TrimSpace(line)
		branches = append(branches, line)
	}
	return branches, nil
}

func GetAddFiles() ([]string, error) {
	content, err := RunCmdWithOutput("git diff --cached --name-only")
	if err != nil {
		return nil, terror.Wrap(err, "call RunCmdWithOutput fail")
	}
	lines := strings.Split(content, "\n")
	dirs := ds.NewSet[string]()
	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = filepath.Dir(line)
		dirs.Insert(line)
	}
	return dirs.Keys(), nil
}
