package pkg

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"github.com/golang/glog"
	"github.com/sunvim/tools/skaffold/pkg/gengrpc/tmpl"
)

const (
	mainFileName = "main.go"
)

func MainCode(pkgName string) error {
	tpl, err := template.New("main_code").Parse(tmpl.MainCodeTpl)
	if err != nil {
		glog.Error(err)
		return err
	}
	buf := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(buf, &pkgName); err != nil {
		glog.Error(err)
		return err
	}

	if err := ioutil.WriteFile(mainFileName, buf.Bytes(), 0666); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}
