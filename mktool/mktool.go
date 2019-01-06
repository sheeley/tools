package mktool

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"unicode"

	"github.com/richardwilkes/toolbox/errs"
)

type Input struct {
	Verbose  bool
	ToolName string
}

type Output struct {
	ToolDir, CmdFile string
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(err)
	}

	return true
}

func clean(s string) string {
	o := strings.Builder{}
	first := true
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			first = true
			continue
		}
		if first {
			first = false
			r = unicode.ToUpper(r)
		}

		o.WriteRune(r)
	}
	return o.String()
}

func Mktool(in *Input) (*Output, error) {
	in.ToolName = strings.TrimSpace(in.ToolName)
	if in.ToolName == "" {
		return nil, errs.New("must include a tool name")
	}
	gp, ok := os.LookupEnv("GOPATH") // TODO: maybe? build.Default.GOPATH
	if !ok {
		return nil, errors.New("$GOPATH not set")
	}

	err := setTemplates()
	if err != nil {
		return nil, errs.Wrap(err)
	}

	toolDir := path.Join(gp, "src/github.com/sheeley/tools/")

	cleaned := clean(in.ToolName)
	pkg := strings.ToLower(cleaned)

	libDir := path.Join(toolDir, pkg)
	libFile := path.Join(libDir, strings.ToLower(in.ToolName)+".go")

	cmdDir := path.Join(libDir, "cmd", in.ToolName)
	cmdFile := path.Join(cmdDir, "main.go")

	t := &tool{
		Package:  pkg,
		FuncName: cleaned,
	}

	// fmt.Println(libFile, cmdFile, t)

	if exists(libDir) {
		return nil, errors.New(fmt.Sprintf("package %s exists", pkg))
	}

	err = os.MkdirAll(cmdDir, os.ModePerm)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	// create the lib file
	f, err := os.Create(libFile)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	err = templateLib.Execute(f, t)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	f.Close()

	// create the cmd file
	f, err = os.Create(cmdFile)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	err = templateMain.Execute(f, t)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	f.Close()

	return &Output{
		ToolDir: toolDir,
		CmdFile: cmdFile,
	}, nil
}
