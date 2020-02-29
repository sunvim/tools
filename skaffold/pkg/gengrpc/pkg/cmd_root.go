package pkg

import (
	"bytes"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/golang/glog"
	"github.com/sunvim/tools/skaffold/pkg/gengrpc/tmpl"
)

func CmdRootCode() error {
	tpl, err := template.New("cmd_root_code").Parse(tmpl.CmdRootTpl)
	if err != nil {
		glog.Error(err)
		return err
	}
	buf := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(buf, nil); err != nil {
		glog.Error(err)
		return err
	}
	// generate cmd directory
	if err := os.MkdirAll("cmd", 0775); err != nil {
		glog.Error(err)
		return err
	}
	// generate root file
	if err := ioutil.WriteFile("cmd/root.go", buf.Bytes(), 0666); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}
