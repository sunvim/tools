package deploy

import (
	"io/ioutil"

	"github.com/golang/glog"
)

const valuesFile = `
app:
  version: v0.1.0`

const valuesFilePath = "deployment/values.yaml"

func genValuesFile() error {
	if err := ioutil.WriteFile(valuesFilePath, []byte(valuesFile), 0666); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}
