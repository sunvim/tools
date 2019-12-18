package cert

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	apiv1beta1 "k8s.io/api/certificates/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	certReqV1Beta1 "k8s.io/client-go/kubernetes/typed/certificates/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
)

const certTpl = `
apiVersion: certificates.k8s.io/v1beta1
kind: CertificateSigningRequest
metadata:
  name: {{ .SrvName }}
spec:
  request: {{ .SrvCsr }}
  usages:
  - digital signature
  - key encipherment
  - server auth

`

const srvCertFile = "admission/server-cert.pem"

func GenServerCert(srvName, csrContent string, kubeconfig string) ([]byte, error) {
	tpl, err := template.New("cert").Parse(certTpl)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	buffer := bytes.NewBuffer([]byte{})
	arg := struct {
		SrvName string
		SrvCsr  string
	}{
		SrvName: srvName,
		SrvCsr:  strings.Replace(base64.StdEncoding.EncodeToString([]byte(csrContent)), "\n", "", -1),
	}
	if err = tpl.Execute(buffer, arg); err != nil {
		glog.Error(err)
		return nil, err
	}
	k8sCsr := &apiv1beta1.CertificateSigningRequest{}
	if err = yaml.Unmarshal(buffer.Bytes(), &k8sCsr); err != nil {
		glog.Error(err)
		return nil, err
	}
	clientSet := getCertificateSigningRequestsClient(kubeconfig)
	if clientSet == nil {
		return nil, errors.New("get csr client failed")
	}
	csrReqClient := clientSet.CertificateSigningRequests()
	// request
	k8sCsr, err = csrReqClient.Create(k8sCsr)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	// approve action
	k8sCsr.Status.Conditions = append(k8sCsr.Status.Conditions, apiv1beta1.CertificateSigningRequestCondition{
		Type:           apiv1beta1.CertificateApproved,
		Reason:         "Mobius Approve",
		Message:        "This CSR was approved by Mobius certificate approve.",
		LastUpdateTime: metav1.Now(),
	})
	_, err = csrReqClient.UpdateApproval(k8sCsr)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	// get certificate from k8s cluster and write output directory
	for {
		k8sCsr, err = csrReqClient.Get(srvName, metav1.GetOptions{})
		if err != nil {
			glog.Error(err)
			return nil, err
		}
		if len(k8sCsr.Status.Certificate) > 10 {
			break
		}
		time.Sleep(time.Second)
	}
	// write to the pem file
	if err = ioutil.WriteFile(srvCertFile, k8sCsr.Status.Certificate, 0666); err != nil {
		glog.Error(err)
		return nil, err
	}
	return k8sCsr.Status.Certificate, nil
}

func getCertificateSigningRequestsClient(kubeconfig string) *certReqV1Beta1.CertificatesV1beta1Client {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		glog.Error(err)
		os.Exit(-1)
	}
	//clientSet, err := kubernetes.NewForConfig(config)
	clientSet, err := certReqV1Beta1.NewForConfig(config)
	if err != nil {
		glog.Error(err)
		os.Exit(-1)
	}
	return clientSet
}
