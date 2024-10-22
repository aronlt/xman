package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aronlt/toolkit/ds"
	"github.com/aronlt/toolkit/terror"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var EmptyBranchErr = errors.New("empty branch error")
var EmptyTagErr = errors.New("empty tag error")

func GitCheckConflict() error {
	err := RunCmd("[ $(git ls-files -u  | cut -f 2 | sort -u | wc -l) -eq 0 ] && exit 0 || exit 1")
	if err != nil {
		return err
	}
	return nil
}

func GitAddAll() error {
	err := RunCmd("git add .")
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
		return EmptyTagErr
	}
	err := RunCmd("git tag " + tag)
	if err != nil {
		return terror.Wrap(err, "call push tag fail")
	}
	return nil
}

func GitPushTag(tag string) error {
	if tag == "" {
		return EmptyTagErr
	}
	err := RunCmd("git push origin " + tag)
	if err != nil {
		return terror.Wrap(err, "call push tag fail")
	}
	return nil
}

func GitTryPullAndCheck() (string, error) {
	current, err := GitCurrentBranch()
	if err != nil {
		return "", terror.Wrap(err, "call GitCurrentBranch fail")
	}
	if err = GitCheckDirtyZone(); err != nil {
		return "", terror.Wrap(err, "can not pull when work dir is dirty")
	}
	err = GitPull(current)
	if err != nil {
		color.Red("call GitPull fail, branch:%s, error:%v", current, err)
	} else {
		err = GitCheckConflict()
		if err != nil {
			return "", terror.Wrap(err, "call GitCheckConflict fail")
		}
	}
	return current, nil
}

func GitSelectBranch(ctx *cli.Context, param string, hint string) (string, error) {
	current, err := GitCurrentBranch()
	if err != nil {
		return "", terror.Wrap(err, "call GitCurrentBranch fail")
	}
	branches, err := ListAllBranch()
	if err != nil {
		return "", terror.Wrap(err, "call ListAllBranch fail")
	}
	branch := ctx.String(param)
	if branch == "" {
		branch = GetFromStdio(hint, false, branches...)
	}
	if branch == current {
		return "", fmt.Errorf("target branch should not equal to current branch:%s", branch)
	}
	if !ds.SliceInclude(branches, branch) {
		return "", fmt.Errorf("select branch:%s not in all branch", branch)
	}
	return branch, nil
}

func GitPull(branch string) error {
	if branch == "" {
		return EmptyBranchErr
	}
	err := RunCmd("git pull origin " + branch)
	if err != nil {
		return terror.Wrap(err, "run cmd fail")
	}
	return nil
}

func GitPush(branch string, force ...bool) error {
	if branch == "" {
		return EmptyBranchErr
	}
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
	content, err := RunCmdWithOutput("git tag", true)
	if err != nil {
		return nil, terror.Wrap(err, "call RunCmdWithOutput fail")
	}
	lines := strings.Split(content, "\n")
	return lines, nil
}

func GitCheckout(branch string, newBranch ...bool) error {
	if branch == "" {
		return EmptyBranchErr
	}
	cmd := "git checkout "
	if len(newBranch) != 0 && newBranch[0] {
		cmd += " -b "
	}
	err := RunCmd(cmd + branch)
	if err != nil {
		return terror.Wrap(err, "run cmd fail")
	}
	return nil
}

func GitMerge(branch string) error {
	if branch == "" {
		return EmptyBranchErr
	}
	err := RunCmd("git merge " + branch)
	if err != nil {
		return terror.Wrap(err, "run cmd fail")
	}
	return nil
}

func GitCurrentBranch() (string, error) {
	out, err := RunCmdWithOutput("git rev-parse --abbrev-ref HEAD", true)
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

	content, err := RunCmdWithOutput("git stash list", true)
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
	content, err := RunCmdWithOutput("git branch", true)
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

func PullRemoteRepository() error {
	err := GitCheckDirtyZone()
	if err != nil {
		logrus.WithError(err).Errorf("dirty working zone")
		return err
	}
	err = RunCmd("git pull")
	if err != nil {
		logrus.WithError(err).Errorf("pull remote repository fail")
		return err
	}
	return nil
}

func GetAddFiles(ignore ds.BuiltinSet[string]) ([]string, error) {
	content, err := RunCmdWithOutput("git diff --cached --name-only", true)
	if err != nil {
		return nil, terror.Wrap(err, "call RunCmdWithOutput fail")
	}
	lines := strings.Split(content, "\n")
	files := ds.NewSet[string]()
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if ignore.Has(line) {
			continue
		}
		files.Insert(line)
	}
	return files.Keys(), nil
}

func PushBranch(branch string) error {
	if branch == "" {
		return EmptyBranchErr
	}
	err := GitPull(branch)
	if err != nil {
		// 远程分支可能不存在，只报错不返回失败
		logrus.WithError(err).Errorf("call GitPull fail")
	}
	err = GitCheckConflict()
	if err != nil {
		return terror.Wrap(err, "call GitCheckConflict fail")
	}

	err = GitPush(branch, true)
	if err != nil {
		return terror.Wrap(err, "call utils.GitPush fail")
	}
	return nil
}

func PullBranch(branch string) error {
	if branch == "" {
		return EmptyBranchErr
	}
	err := GitPull(branch)
	if err != nil {
		return terror.Wrapf(err, "call GitPushfail, branch:%s", branch)
	}
	err = GitCheckConflict()
	if err != nil {
		return terror.Wrap(err, "call GitCheckConflict fail")
	}
	return nil
}
