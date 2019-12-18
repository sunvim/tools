package gencert

import (
	"os"

	"github.com/sunvim/tools/skaffold/pkg/gencert/secret"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/sunvim/tools/skaffold/pkg/constant"
	"github.com/sunvim/tools/skaffold/pkg/gencert/cert"
	"github.com/sunvim/tools/skaffold/pkg/gencert/config"
	"github.com/sunvim/tools/skaffold/pkg/gencert/csr"
)

func PrepareBasicWorkEnvironments(cmd *cobra.Command, args []string) {
	if err := os.MkdirAll("admission", 0775); err != nil {
		glog.Error(err)
	}

	srvName, err := cmd.Flags().GetString(constant.ServiceName)
	if err != nil {
		glog.Error(err)
		return
	}
	namespace, err := cmd.Flags().GetString(constant.Namespace)
	if err != nil {
		glog.Error(err)
		return
	}
	nsLabel, err := cmd.Flags().GetString(constant.NsLabel)
	if err != nil {
		glog.Error(err)
		return
	}
	kubeconfig, err := cmd.Flags().GetString(constant.KubernetesConfig)
	input, err := csr.GenerateKeyCsr(srvName, namespace)
	if err != nil {
		glog.Error(err)
		return
	}
	certPem, err := cert.GenServerCert(srvName, input["csr"].(string), kubeconfig)
	if err != nil {
		glog.Error(err)
		return
	}
	config.ApplyMutatingWebHookConfig(srvName, namespace, nsLabel, kubeconfig)

	if err = secret.CreateCertSecret(srvName, namespace, string(certPem), input["key"].(string),
		kubeconfig); err != nil {
		glog.Error(err)
		return
	}
}
