package component

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aronlt/toolkit/ds"
	"github.com/aronlt/toolkit/terror"
	"github.com/aronlt/toolkit/tio"
	"github.com/sirupsen/logrus"
)

type ModInfo struct {
	Filename string
	Lines    [][]byte
}

func NewModInfo() (*ModInfo, error) {
	wd, err := os.Getwd()
	if err != nil {
		err = terror.Wrap(err, "call os.Getwd fail")
		return nil, err
	}
	filename := filepath.Join(wd, "go.mod")
	lines, err := tio.ReadLines(filename)
	if err != nil {
		err = terror.Wrap(err, "call tio.ReadLines fail")
		return nil, err
	}
	return &ModInfo{
		Filename: filename,
		Lines:    lines,
	}, nil
}

func (m *ModInfo) Refresh() error {
	_, err := tio.WriteFile(m.Filename, bytes.Join(m.Lines, []byte("\n")), false)
	if err != nil {
		return terror.Wrap(err, "tio.WriteFile fail")
	}
	err = os.Chmod(m.Filename, 0755)
	if err != nil {
		return terror.Wrap(err, "os.Chmod fail")
	}
	return nil
}

func (m *ModInfo) ListModuleNames() ([]string, error) {
	inRequire := false
	modules := make([]string, 0, len(m.Lines))
	for _, lineContent := range m.Lines {
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
				return nil, fmt.Errorf("warning: invalid module, line:%s", line)
			}
			module := parts[0]
			segments := strings.Split(module, "/")
			last := ds.SliceGetNthTail(segments, 1)
			if strings.HasPrefix(last, "v") {
				sub := last[1:]
				if _, err := strconv.Atoi(sub); err == nil {
					modules = append(modules, ds.SliceGetNthTail(segments, 2))
				} else {
					modules = append(modules, last)
				}
			} else {
				modules = append(modules, last)
			}
		}
	}
	return modules, nil
}

func (m *ModInfo) Replace(moduleName string, branchName string) (error, bool) {
	inRequire := false
	replaceCount := 0
	for i, lineContent := range m.Lines {
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
				return fmt.Errorf("warning: invalid module, line:%s", line), false
			}
			module := parts[0]
			if strings.HasSuffix(module, moduleName) {
				parts[1] = branchName
				newLine := strings.Join(parts, " ")
				m.Lines[i] = []byte(newLine)
				logrus.Infof("replace module from:%s -> %s", line, newLine)
				replaceCount += 1
			}
		}
	}
	return nil, replaceCount > 0
}
