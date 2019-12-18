package secret

import (
	"bytes"
	"encoding/base64"
	"encoding/pem"
	"html/template"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/kubernetes/client-go/tools/clientcmd"
	corev1 "k8s.io/api/core/v1"
	secv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const certSecretFile = `
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: {{ .SrvName }}-certs
  namespace: {{ .Namespace }}
data:
  cert.pem: $CertPem$
  key.pem:  $KeyPem$
`

func CreateCertSecret(srvName, namespace, certPem, keyPem, kubeconfig string) error {
	tpl, err := template.New("secrets").Parse(certSecretFile)
	if err != nil {
		glog.Error(err)
		return err
	}
	cps := getPemBody(certPem)
	kps := getPemBody(keyPem)
	arg := struct {
		SrvName   string
		Namespace string
	}{
		SrvName:   srvName,
		Namespace: namespace,
	}
	buf := bytes.NewBuffer([]byte{})
	if err = tpl.Execute(buf, &arg); err != nil {
		glog.Error(err)
		return err
	}
	bs := strings.Replace(buf.String(), "$CertPem$", cps, -1)
	bs = strings.Replace(bs, "$KeyPem$", kps, -1)
	sec := corev1.Secret{}
	if err = yaml.Unmarshal([]byte(bs), &sec); err != nil {
		glog.Error(err)
		return err
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		glog.Error(err)
		os.Exit(-1)
	}
	coreV1Client, err := secv1.NewForConfig(config)
	if err != nil {
		glog.Error(err)
		return err
	}
	if _, err = coreV1Client.Secrets(namespace).Create(&sec); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}

func getPemBody(body string) string {
	block, _ := pem.Decode([]byte(body))
	return base64.StdEncoding.EncodeToString(block.Bytes)
}
