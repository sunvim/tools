
package srv

import (
	"context"
	"net"
	"os"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func Service(ctx context.Context, sig chan struct{}, cmd *cobra.Command) {

	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
		glog.Flush()
	}()

	// ready to start service
	port := cmd.Flag("port").Value.String()
	l, err := net.Listen("tcp", port)
	if err != nil {
		glog.Error(err)
		return
	}

	gs := grpc.NewServer()

	go func() {
		<-ctx.Done()
		glog.Info("grpc http2 service exiting ...")
		gs.GracefulStop()
	}()

	//register all kinds of services
	registerService(gs)

	// send started signal
	close(sig)

	glog.Info("grpc service starting ...., port: ", port)
	if err := gs.Serve(l); err != nil {
		glog.Error(err)
		os.Exit(-1)
	}

}

