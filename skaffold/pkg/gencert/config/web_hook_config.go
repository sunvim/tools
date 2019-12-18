package config

import (
	"bytes"
	"encoding/base64"
	"os"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	metav1beta1 "k8s.io/api/admissionregistration/v1beta1"
	admv1beta1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1beta1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const whConfig = `
apiVersion: admissionregistration.k8s.io/v1beta1
kind: MutatingWebhookConfiguration
metadata:
  name: {{ .SrvName }}-webhook-cfg
  labels:
    app: kms-injector
webhooks:
  - name: {{ .SrvName }}.{{ .Namespace }}.mobius
    clientConfig:
      service:
        name: {{ .SrvName }}
        namespace: {{ .Namespace }}
        path: "/mutate"
      caBundle: {{ .CaBundle }}
    rules:
      - operations: [ "CREATE" ]
        apiGroups: [""]
        apiVersions: ["v1"]
        resources: ["pods"]
    namespaceSelector:
      matchLabels:
        {{ .NsLabel }}: enabled
`

func ApplyMutatingWebHookConfig(srvName, namespace, nsLabel string, kubeconfig string) {
	tpl, err := template.New("mutatingWebHookConfig").Parse(whConfig)
	if err != nil {
		glog.Error(err)
		return
	}

	config, caData := getCaBundleFromCluster(kubeconfig)
	buffer := bytes.NewBuffer([]byte{})

	arg := struct {
		SrvName   string
		Namespace string
		CaBundle  string
		NsLabel   string
	}{
		SrvName:   srvName,
		Namespace: namespace,
		CaBundle:  caData,
		NsLabel:   nsLabel,
	}

	if err = tpl.Execute(buffer, arg); err != nil {
		glog.Error(err)
		return
	}

	if err = applyConfig(config, buffer.Bytes()); err != nil {
		glog.Error(err)
		return
	}
}

func getCaBundleFromCluster(kubeconfig string) (*rest.Config, string) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		glog.Error(err)
		os.Exit(-1)
	}

	return config, base64.StdEncoding.EncodeToString(config.CAData)
}

func applyConfig(config *rest.Config, mwhConfig []byte) error {
	admClient, err := admv1beta1.NewForConfig(config)
	if err != nil {
		return err
	}
	mwh := &metav1beta1.MutatingWebhookConfiguration{}
	if err = yaml.Unmarshal(mwhConfig, &mwh); err != nil {
		return err
	}
	_, err = admClient.MutatingWebhookConfigurations().Create(mwh)
	if err != nil {
		return err
	}
	return nil
}
