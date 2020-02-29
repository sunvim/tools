package gengrpc

import (
	"flag"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/sunvim/tools/skaffold/pkg/gengrpc/pkg"
)

func GrpcRun(cmd *cobra.Command, args []string) {
	flag.Parse()
	pkgName := cmd.Flag("package").Value.String()
	if pkgName == "" {
		glog.Error("")
	}

	// generate main code file
	glog.Info("generate main code file")
	if err := pkg.MainCode(pkgName); err != nil {
		glog.Error(err)
		return
	}

	// generate mod file
	glog.Info("generate mod file")
	if err := pkg.GenerateMod(pkgName); err != nil {
		glog.Error(err)
		return
	}

	// generate root command file
	glog.Info("generate root command file")
	if err := pkg.CmdRootCode(); err != nil {
		glog.Error(err)
		return
	}

	// generate server file
	glog.Info("generate server file")
	if err := pkg.CmdServerCode(pkgName); err != nil {
		glog.Error(err)
		return
	}

	// generate srv/main_run.go file
	glog.Info("generate srv/main_run.go file")
	if err := pkg.SrvMainRunCode(); err != nil {
		glog.Error(err)
		return
	}

	// generate srv/grpc_service.go file
	glog.Info("generate srv/grpc_service.go file")
	if err := pkg.SrvGrpcServiceCode(); err != nil {
		glog.Error(err)
		return
	}

	// generate gateway files
	glog.Info("generate gateway files")
	if err := pkg.CreateGateWayFiles(); err != nil {
		glog.Error(err)
		return
	}

	// generate proto buffer files
	glog.Info("generate proto buffer files")
	if err := pkg.GenerateProtoBuffers(); err != nil {
		glog.Error(err)
		return
	}

	// generate swagger files
	glog.Info("generate swagger files")
	if err := pkg.GenerateSwaggerFiles(); err != nil {
		glog.Error(err)
		return
	}
}
