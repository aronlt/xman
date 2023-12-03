package component

import (
	"sort"
	"strconv"
	"strings"

	"github.com/aronlt/toolkit/ds"
	"github.com/aronlt/toolkit/terror"
	"github.com/aronlt/xman/component/utils"
	"github.com/hashicorp/go-version"
	"github.com/urfave/cli/v2"
)

type Tag struct {
}

func NewTag() *Tag {
	return &Tag{}
}

func (t *Tag) Run(ctx *cli.Context) error {
	prefix := ctx.String("prefix")
	suffix := ctx.String("suffix")
	err := utils.GitCheckDirtyZone()
	if err != nil {
		err = utils.GitAddAndCommit()
		if err != nil {
			return terror.Wrap(err, "call GitAddAndCommit fail")
		}
	}
	err = utils.GitPullTags()
	if err != nil {
		return terror.Wrap(err, "call GitPullTags fail")
	}
	lines, err := utils.GitTags()
	if err != nil {
		return terror.Wrap(err, "call GitTags fail")
	}
	lines = ds.SliceIterFilter(lines, func(a []string, i int) bool {
		ok := false
		if prefix != "" {
			ok = strings.HasPrefix(a[i], prefix)
		} else {
			ok = !strings.Contains(a[i], "/")
		}
		return ok
	})

	lines = ds.SliceIterFilter(lines, func(a []string, i int) bool {
		ok := false
		if suffix != "" {
			ok = strings.HasSuffix(a[i], suffix)
		} else {
			ok = !strings.Contains(a[i], "-")
		}
		return ok
	})

	if prefix != "" || suffix != "" {
		var p string
		var s string
		if prefix != "" {
			p = prefix + "/"
		}
		if suffix != "" {
			s = "-" + suffix
		}
		ds.SliceIter(lines, func(a []string, i int) {
			if p != "" {
				a[i] = strings.TrimPrefix(a[i], p)
			}
			if s != "" {
				a[i] = strings.TrimSuffix(a[i], s)
			}
		})
	}
	var max string
	if len(lines) != 0 {
		versions := make([]*version.Version, len(lines))
		for i, raw := range lines {
			v, err := version.NewVersion(raw)
			if err != nil {
				return terror.Wrapf(err, "call NewVersion fail, raw:%s", raw)
			}
			versions[i] = v
		}
		// After this, the versions are properly sorted
		sort.Sort(version.Collection(versions))
		max = ds.SliceGetTail(versions).String()
	} else {
		max = "1.0.0"
	}

	parts := strings.Split(max, ".")
	lastV, err := strconv.Atoi(ds.SliceGetTail(parts))
	if err != nil {
		return terror.Wrap(err, "call atoi fail")
	}
	lastV += 1
	lastVS := strconv.Itoa(lastV)
	parts[len(parts)-1] = lastVS
	maxV := strings.Join(parts, ".")
	if !strings.HasPrefix(maxV, "v") {
		maxV = "v" + maxV
	}
	if prefix != "" {
		maxV = prefix + "/" + maxV
	}
	if suffix != "" {
		maxV = maxV + "-" + suffix
	}
	err = utils.GitAddTag(maxV)
	if err != nil {
		return terror.Wrap(err, "call AddTag fail")
	}
	err = utils.GitPushTag(maxV)
	if err != nil {
		return terror.Wrap(err, "call PushTag fail")
	}
	return nil
}

func (t *Tag) Name() string {
	return "tag"

}

func (t *Tag) Usage() string {
	return "自动生成tag"
}

func (t *Tag) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "prefix",
			Aliases: []string{"p"},
			Usage:   "tag 前缀信息",
		},
		&cli.StringFlag{
			Name:    "suffix",
			Aliases: []string{"s"},
			Usage:   "tag 后缀信息",
		}}
}
