package pkg

import (
	"bytes"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/golang/glog"
	"github.com/sunvim/tools/skaffold/pkg/gengrpc/tmpl"
)

func CreateGateWayFiles() error {
	createDir()

	glog.Info("create gateway main file")
	if err := createGatewayMainFile(); err != nil {
		glog.Error(err)
		return err
	}

	glog.Info("create gateway handler file")
	if err := createGatewayHandlers(); err != nil {
		glog.Error(err)
		return err
	}

	glog.Info("create gateway file")
	if err := createGatewayFile(); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}

func createDir() {
	if err := os.MkdirAll("gateway", 0666); err != nil {
		glog.Error(err)
	}
}

func createGatewayMainFile() error {
	tpl, err := template.New("gateway_main").Parse(tmpl.GatewayMainTpl)
	if err != nil {
		glog.Error(err)
		return err
	}
	buf := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(buf, nil); err != nil {
		glog.Error(err)
		return err
	}

	if err := ioutil.WriteFile("gateway/main.go", buf.Bytes(), 0666); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}

func createGatewayHandlers() error {
	tpl, err := template.New("gateway_handler").Parse(tmpl.GatewayHandlerTpl)
	if err != nil {
		glog.Error(err)
		return err
	}
	buf := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(buf, nil); err != nil {
		glog.Error(err)
		return err
	}

	if err := ioutil.WriteFile("gateway/handler.go", buf.Bytes(), 0666); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}

func createGatewayFile() error {
	tpl, err := template.New("gateway_gateway").Parse(tmpl.GateWayTpl)
	if err != nil {
		glog.Error(err)
		return err
	}
	buf := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(buf, nil); err != nil {
		glog.Error(err)
		return err
	}

	if err := ioutil.WriteFile("gateway/gateway.go", buf.Bytes(), 0666); err != nil {
		glog.Error(err)
		return err
	}
	return nil
	return nil
}
