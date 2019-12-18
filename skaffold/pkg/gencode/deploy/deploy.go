package deploy

import (
	"os"

	"github.com/golang/glog"
)

const dirPath = "deployment/templates"

func Generate(companyName, srvName string) error {
	if err := os.MkdirAll(dirPath, 0775); err != nil {
		glog.Error(err)
		return err
	}
	//1: generate chart file
	if err := genChartFile(srvName); err != nil {
		glog.Error(err)
		return err
	}
	//2: generate value file
	if err := genValuesFile(); err != nil {
		glog.Error(err)
		return err
	}
	//3: generate template files
	if err := genTemplates(companyName, srvName); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}
