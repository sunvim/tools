package pkg

import (
	"bytes"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/golang/glog"
	"github.com/sunvim/tools/skaffold/pkg/gengrpc/tmpl"
)

func CmdServerCode(pkgName string) error {
	tpl, err := template.New("cmd_server_code").Parse(tmpl.CmdServerTpl)
	if err != nil {
		glog.Error(err)
		return err
	}
	buf := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(buf, &pkgName); err != nil {
		glog.Error(err)
		return err
	}

	if err := os.MkdirAll("srv", 0775); err != nil {
		glog.Error(err)
		return err
	}

	if err := ioutil.WriteFile("cmd/server.go", buf.Bytes(), 0666); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}
