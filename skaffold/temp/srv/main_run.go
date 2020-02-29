
package srv

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func MainRun(cmd *cobra.Command, args []string) {
	defer func() {
		glog.Flush()
	}()

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan struct{})

	// start grpc service
	go Service(ctx, sig, cmd)
	// start http service
	go HttpService(ctx, sig, cmd)

	select {
	case <-stopChan:
		cancel()
		time.Sleep(1 * time.Second)
		glog.Info("all services exited totally.")
	}
