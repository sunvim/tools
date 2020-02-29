package pkg

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func GenerateProtoBuffers(cmd *cobra.Command) error {
	protoImpPath := cmd.Flag("proto_path").Value.String()
	protoSrvPath := cmd.Flag("proto").Value.String()
	protoCmd := fmt.Sprintf("protoc --proto_path=%s:. --go_out=plugins=grpc:. --grpc-gateway_out=logtostderr=true:. "+
		"--govalidators_out=. --swagger_out=logtostderr=true:swagger %s/*.proto", protoImpPath, protoSrvPath)
	if err := exec.Command(protoCmd).Run(); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}

const (
	protoCmd             = "protoc"
	protoGenGo           = "protoc-gen-go"
	protoGenGateway      = "protoc-gen-gateway"
	protocGenSwagger     = "protoc-gen-swagger"
	protocGenGrpcGateway = "protoc-gen-grpc-gateway"
)

func precondition() error {
	protoCmdPath := fullPath(protoCmd, strings.Split(os.Getenv("PATH"), ":"))
	if protoCmdPath == "" {
		glog.Errorf("Failed finding plugin binary %s\n", protoCmd)
		return errors.New("can't find ptotoc command")
	}

	protoGenGoPath := fullPath(protoGenGo, strings.Split(os.Getenv("PATH"), ":"))
	if protoGenGoPath == "" {
		glog.Errorf("Failed finding plugin binary %s\n", protoGenGo)
		return errors.New("can't find protoc-gen-go command")
	}

	protoGenGatewayPath := fullPath(protoGenGateway, strings.Split(os.Getenv("PATH"), ":"))
	if protoGenGatewayPath == "" {
		glog.Errorf("Failed finding plugin binary %s\n", protoGenGateway)
		return errors.New("can't find protoc-gen-gateway command")
	}

	protocGenSwaggerPath := fullPath(protocGenSwagger, strings.Split(os.Getenv("PATH"), ":"))
	if protocGenSwaggerPath == "" {
		glog.Errorf("Failed finding plugin binary %s\n", protocGenSwagger)
		return errors.New("can't find protoc-gen-swagger command")
	}

	protocGenGrpcGatewayPath := fullPath(protocGenGrpcGateway, strings.Split(os.Getenv("PATH"), ":"))
	if protocGenGrpcGatewayPath == "" {
		glog.Errorf("Failed finding plugin binary %s\n", protocGenGrpcGateway)
		return errors.New("can't find protoc-gen-grpc-gateway command")
	}
	return nil
}

func fullPath(binary string, paths []string) string {
	if strings.Index(binary, "/") >= 0 {
		// path with path component
		return binary
	}
	for _, p := range paths {
		full := path.Join(p, binary)
		fi, err := os.Stat(full)
		if err == nil && !fi.IsDir() {
			return full
		}
	}
	return ""
}
