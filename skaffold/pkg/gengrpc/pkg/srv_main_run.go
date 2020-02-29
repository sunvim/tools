package pkg

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"github.com/golang/glog"
	"github.com/sunvim/tools/skaffold/pkg/gengrpc/tmpl"
)

func SrvMainRunCode() error {
	tpl, err := template.New("srv_main_run_code").Parse(tmpl.SrvMainRunTpl)
	if err != nil {
		glog.Error(err)
		return err
	}
	buf := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(buf, nil); err != nil {
		glog.Error(err)
		return err
	}

	if err := ioutil.WriteFile("srv/main_run.go", buf.Bytes(), 0666); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}
