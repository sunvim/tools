package deploy

import (
	"bytes"
	"html/template"
	"io/ioutil"

	"github.com/golang/glog"
)

const chartFile = `
apiVersion: v1
appVersion: "1.0"
description: A Helm chart for Kubernetes
name: {{.}}
version: 0.1.0
`
const chartFileName = "deployment/Chart.yaml"

func genChartFile(srvName string) error {
	chartTpl, err := template.New("chartFile").Parse(chartFile)
	if err != nil {
		glog.Error(err)
		return err
	}
	buf := bytes.NewBuffer([]byte{})
	if err = chartTpl.Execute(buf, &srvName); err != nil {
		glog.Error(err)
		return err
	}
	if err = ioutil.WriteFile(chartFileName, buf.Bytes(), 0666); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}
