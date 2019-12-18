package constant

const MainFrameGoFile = `
package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"{{ .PkgName }}/webhook"
	"github.com/golang/glog"
)

func main() {
	var parameters webhook.WhSvrParameters

	// get command line parameters
	flag.IntVar(&parameters.Port, "port", 443, "web hook server port.")
	flag.StringVar(&parameters.CertFile, "tlsCertFile", "/etc/hook/certs/cert.pem",
		"file containing the x509 certificate for HTTPS.")
	flag.StringVar(&parameters.KeyFile, "tlsKeyFile", "/etc/hook/certs/key.pem",
		"file containing the x509 private key to --tlsCertFile.")
	flag.StringVar(&parameters.InjectorFile, "injectorFile", "/etc/hook/config/injector.yaml",
		"file containing the mutation configuration.")
	if err := flag.Set("alsologtostderr", "true"); err != nil {
		glog.Errorf("set command line flag failed: %v \n", err)
		return
	}
	if err := flag.Set("stderrthreshold", "INFO"); err != nil {
		glog.Errorf("set command line flag failed: %v \n", err)
		return
	}
	flag.Parse()

	injectConfig, err := webhook.LoadConfig(parameters.InjectorFile)
	if err != nil {
		glog.Errorf("Filed to load configuration: %v", err)
		return
	}

	pair, err := tls.LoadX509KeyPair(parameters.CertFile, parameters.KeyFile)
	if err != nil {
		glog.Errorf("Filed to load key pair: %v", err)
		return
	}

	whSvr := &webhook.Server{
		InjectorConfig: injectConfig,
		Server: &http.Server{
			Addr:      fmt.Sprintf(":%v", parameters.Port),
			TLSConfig: &tls.Config{Certificates: []tls.Certificate{pair}},
		},
	}

	glog.Infof("server inject config:\n %+v \n", whSvr.InjectorConfig)

	// define http server and server handler
	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", whSvr.Serve)
	whSvr.Server.Handler = mux

	// start web hook server in new goroutine
	go func() {
		glog.Info("kms messenger service starting ...")
		if err := whSvr.Server.ListenAndServeTLS("", ""); err != nil {
			glog.Errorf("Filed to listen and serve web hook server: %v", err)
			return
		}
	}()

	// listening OS shutdown signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	glog.Infof("Got OS shutdown signal, shutting down web hook server gracefully...")
	if err = whSvr.Server.Shutdown(context.Background()); err != nil {
		glog.Error(err)
	}
}
`
