package frame

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/sunvim/tools/skaffold/pkg/constant"
	"github.com/sunvim/tools/skaffold/pkg/gencode/deploy"
)

func GenCode(cmd *cobra.Command, args []string) {

	// main code frame
	pkgName, err := cmd.Flags().GetString(constant.PkgName)
	if err != nil {
		glog.Error(err)
		return
	}

	// create project module
	pkgCmd := exec.Command("go", "mod", "init", pkgName)
	if err := pkgCmd.Run(); err != nil {
		glog.Error(err)
		return
	}

	srvName, err := cmd.Flags().GetString(constant.ServiceName)
	if err != nil {
		glog.Error(err)
		return
	}
	arg := struct {
		SrvName string
		PkgName string
	}{
		SrvName: srvName,
		PkgName: pkgName,
	}
	buf := bytes.NewBuffer([]byte{})
	//1: generate main frame code
	mainFrame, err := template.New("mainFrame").Parse(constant.MainFrameGoFile)
	if err != nil {
		glog.Error(err)
		return
	}
	if err = mainFrame.Execute(buf, &arg); err != nil {
		glog.Error(err)
		return
	}
	if err = ioutil.WriteFile(constant.MainGoFile, buf.Bytes(), 0666); err != nil {
		glog.Error(err)
		return
	}
	//2: generate web hook code
	if err = os.MkdirAll(constant.WebHookDir, 0775); err != nil {
		glog.Error(err)
		return
	}
	whBytes, err := base64.StdEncoding.DecodeString(constant.WebHookGoFile)
	if err != nil {
		glog.Error(err)
		return
	}
	buf.Reset()
	webHookTpl, err := template.New("webHookCode").Parse(string(whBytes))
	if err != nil {
		glog.Error(err)
		return
	}
	if err = webHookTpl.Execute(buf, &arg); err != nil {
		glog.Error(err)
		return
	}
	if err = ioutil.WriteFile(constant.WebHookDir+"/"+constant.WebHookFile, buf.Bytes(), 0666); err != nil {
		glog.Error(err)
		return
	}
	//3: generate Dockerfile file
	if err = ioutil.WriteFile(constant.DockerFileName, []byte(constant.DockerFile), 0666); err != nil {
		glog.Error(err)
		return
	}
	//4: generate Makefile file
	companyName, err := cmd.Flags().GetString(constant.CompanyName)
	if err != nil {
		glog.Error(err)
		return
	}
	makeFileTpl, err := template.New("makefile").Parse(constant.Makefile)
	if err != nil {
		glog.Error(err)
		return
	}
	mkArg := struct {
		CompanyName string
		SrvName     string
	}{
		CompanyName: companyName,
		SrvName:     srvName,
	}
	buf.Reset()
	if err = makeFileTpl.Execute(buf, &mkArg); err != nil {
		glog.Error(err)
		return
	}
	if err = ioutil.WriteFile(constant.MakeFileName, buf.Bytes(), 0666); err != nil {
		glog.Error(err)
		return
	}
	//5: generate deployment file
	if err = deploy.Generate(companyName, srvName); err != nil {
		glog.Error(err)
		return
	}
}
