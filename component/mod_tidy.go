package component

import (
	"fmt"

	"github.com/aronlt/toolkit/ds"
	"github.com/aronlt/toolkit/terror"
	"github.com/aronlt/xman/component/utils"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type Tidy struct {
}

func NewTidy() *Tidy {
	return &Tidy{}
}

func (t *Tidy) PreCheck() error {
	err := utils.GitCheckDirtyZone()
	if err != nil {
		return terror.Wrap(err, "work zone is dirty")
	}
	return nil
}

func (t *Tidy) Name() string {
	return "mod"
}

func (t *Tidy) Usage() string {
	return "更新依赖的模块信息"
}

func (t *Tidy) Run(_ *cli.Context) error {
	if err := t.PreCheck(); err != nil {
		return terror.Wrap(err, "call PreCheck fail")
	}

	modInfo, err := NewModInfo()
	if err != nil {
		terror.Wrap(err, "call getModInfo fail")
	}
	modules, err := modInfo.ListModuleNames()
	if err != nil {
		terror.Wrap(err, "call ListModuleNames fail")
	}
	branches, err := utils.ListAllBranch()
	if err != nil {
		terror.Wrap(err, "call ListAllBranch fail")
	}
	pushBranch := utils.GetFromStdio("推送到什么分支", branches...)
	moduleName := utils.GetFromStdio("要替换的模块名", modules...)
	branchName := utils.GetFromStdio("模块替换后的分支")
	if ds.SliceIncludeUnpack("", pushBranch, moduleName, branchName) {
		return fmt.Errorf("invalid arguments:%+v, %+v, %+v", pushBranch, moduleName, branchName)
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
	err = t.replace(modInfo, moduleName, branchName)
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
	err = utils.GitPush(pushBranch)
	if err != nil {
		return err
	}
	logrus.Infof("run script success...")
	return nil
}

func (t *Tidy) replace(modInfo *ModInfo, moduleName string, branchName string) error {

	err, ok := modInfo.Replace(moduleName, branchName)
	if err != nil {
		terror.Wrap(err, "call Replace fail")
	}
	if !ok {
		logrus.Infof("not replace any module, skip")
	}
	err = modInfo.Refresh()
	if err != nil {
		return terror.Wrap(err, "call modInfo.Refresh fail")
	}
	logrus.Infof("replace success...")
	return nil
}
