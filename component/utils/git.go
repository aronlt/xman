package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aronlt/toolkit/terror"
)

func GitCheckConflict() error {
	err := RunCmd("[ $(git ls-files -u  | cut -f 2 | sort -u | wc -l) -eq 0 ] && exit 0 || exit 1")
	if err != nil {
		return err
	}
	return nil
}

func GitAddAndCommit() error {
	err := RunCmd("git add .")
	if err != nil {
		return err
	}
	err = RunInteractiveCmd("git", []string{"commit"})
	return err
}

func GitCheckDirtyZone() error {
	err := RunCmd("[ $(git status -s | wc -l) -eq 0 ] && exit 0 || exit 1")
	if err != nil {
		return terror.Wrap(err, "work zone is dirty")
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

func GitPush(branch string) error {
	err := RunCmd("git push origin " + branch)
	if err != nil {
		return terror.Wrap(err, "run cmd fail")
	}
	return nil
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
	line := GetFromStdio("恢复的行号(下标0开始)")
	lineNum, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		return terror.Wrap(err, "call strconv.Atoi fail")
	}
	err = RunCmd(fmt.Sprintf("git stash pop stash@{%d}", lineNum))
	if err != nil {
		return terror.Wrap(err, "call RunCmd fail")
	}
	return nil
}

func GitAddWithConfirm() error {
	err := RunInteractiveCmd("git", []string{"add", "-p", "."})
	if err != nil {
		return terror.Wrap(err, "call RunInteractiveCmd fail")
	}
	err = RunInteractiveCmd("git", []string{"commit"})
	if err != nil {
		return terror.Wrap(err, "call RunInteractiveCmd fail")
	}
	return err
}
