package pkg

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"github.com/golang/glog"
	"github.com/sunvim/tools/skaffold/pkg/gengrpc/tmpl"
)

func SrvGrpcServiceCode() error {
	tpl, err := template.New("srv_grpc_service_code").Parse(tmpl.SrvGrpcServiceTpl)
	if err != nil {
		glog.Error(err)
		return err
	}
	buf := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(buf, nil); err != nil {
		glog.Error(err)
		return err
	}

	if err := ioutil.WriteFile("srv/grpc_service.go", buf.Bytes(), 0666); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}
