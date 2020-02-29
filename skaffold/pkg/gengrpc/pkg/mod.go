package pkg

import (
	"os/exec"

	"github.com/golang/glog"
)

func GenerateMod(pkgName string) error {
	err := exec.Command("go", "mod", "init", pkgName).Run()
	if err != nil {
		glog.Error(err)
		return err
	}
	return nil
}
